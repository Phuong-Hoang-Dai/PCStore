// Places ORDERS orders at the same moment (one per virtual user), then checks that stock stayed
// consistent: every order that went through holds exactly one unit, so nothing was oversold or
// lost under concurrency.
//
// Run from backend/DDStore with the stack up (k6 runs in Docker, on the compose network):
//   bash:       docker run --rm -i --network ddstore_default -e BASE_URL=http://api_gateway:10000/api/v1 grafana/k6:2.2.0 run - < loadtest/place-orders.js
//   PowerShell: Get-Content loadtest/place-orders.js -Raw | docker run --rm -i --network ddstore_default -e BASE_URL=http://api_gateway:10000/api/v1 grafana/k6:2.2.0 run -
// Options, as -e NAME=value: ORDERS (default 1000), MODE (user | guest), PRODUCT_ID (default: the
// product with the most stock), ADMIN_NAME / ADMIN_PASSWORD (default: the seeded admin1).
// Watch it live in Grafana (http://localhost:3000) when started with --profile monitoring.
//
// MODE=user (default): one temporary account places every order through POST /order. setup()
//   adds ORDERS units of stock so each order can succeed; teardown() cancels the orders, which
//   gives that stock back, and deletes the account. The added stock stays on the product.
// MODE=guest: POST /order/guest with a different email per virtual user (k6-guest-<n>@example.com,
//   reused across runs). Each order also creates a Stripe customer and checkout session in test
//   mode, and Stripe's test-mode rate limit fails most of 1000 at once, so this mostly measures
//   Stripe. teardown() cancels the orders that are left.

import http from "k6/http";
import { check, fail } from "k6";
import { Counter } from "k6/metrics";

const BASE_URL = __ENV.BASE_URL || "http://localhost:9999/api/v1";
const ORDERS = parseInt(__ENV.ORDERS || "1000", 10);
const MODE = __ENV.MODE || "user";
const ADMIN_NAME = __ENV.ADMIN_NAME || "admin1";
const ADMIN_PASSWORD = __ENV.ADMIN_PASSWORD || "admin01";

// Mirrors model.Cancelled in order_service.
const CANCELLED = 4;

const placed = new Counter("orders_placed");
const rejected = new Counter("orders_rejected");

export const options = {
  scenarios: {
    // Every VU is initialised first, then all of them send their one order together.
    place_orders: {
      executor: "per-vu-iterations",
      vus: ORDERS,
      iterations: 1,
      maxDuration: "2m",
    },
  },
  setupTimeout: "1m",
  teardownTimeout: "10m",
  thresholds: {
    "http_req_failed{name:place order}": ["rate<0.01"],
    "http_req_duration{name:place order}": ["p(95)<5000"],
    checks: ["rate==1"],
  },
  summaryTrendStats: ["min", "med", "avg", "p(90)", "p(95)", "p(99)", "max"],
};

const jsonHeaders = { "Content-Type": "application/json" };

// The gateway reads the JWT from the access_token cookie.
function auth(token) {
  return { headers: { ...jsonHeaders, Cookie: `access_token=${token}` } };
}

function login(name, password) {
  const res = http.post(`${BASE_URL}/user/login`, JSON.stringify({ name, password }), { headers: jsonHeaders });
  const cookie = res.cookies.access_token;
  if (res.status !== 200 || !cookie || cookie.length === 0) {
    fail(`login as ${name}: ${res.status} ${res.body}`);
  }
  return cookie[0].value;
}

function getStock(admin, productId) {
  const rows = http.get(`${BASE_URL}/stock?ids=${productId}`, auth(admin)).json("data") || [];
  return rows.length > 0 ? rows[0].quantity : 0;
}

export function setup() {
  const runId = `${Date.now()}`;
  const admin = login(ADMIN_NAME, ADMIN_PASSWORD);

  let productId = __ENV.PRODUCT_ID;
  if (!productId) {
    const stocks = http.get(`${BASE_URL}/stock`, auth(admin)).json("data") || [];
    if (stocks.length === 0) {
      fail("no product has stock yet; add some in /manage or set PRODUCT_ID");
    }
    productId = stocks.sort((a, b) => b.quantity - a.quantity)[0].productId;
  }
  const stockBefore = getStock(admin, productId);

  // Enough stock for every order, so failures come from the load, not from running out.
  const add = http.post(`${BASE_URL}/stock/${productId}/add`, JSON.stringify({ quantity: ORDERS }), auth(admin));
  if (add.status !== 200) {
    fail(`add stock: ${add.status} ${add.body}`);
  }

  let user = null;
  if (MODE === "user") {
    const name = `k6_${runId}`;
    const password = `k6-${runId}`;
    const res = http.post(
      `${BASE_URL}/user`,
      JSON.stringify({ name, email: `k6-${runId}@example.com`, password }),
      { headers: jsonHeaders }
    );
    if (res.status !== 200) {
      fail(`create test user: ${res.status} ${res.body}`);
    }
    user = { id: res.json("data.id"), token: login(name, password) };
  }

  console.log(`placing ${ORDERS} ${MODE} orders of product ${productId} (stock ${stockBefore} + ${ORDERS})`);
  // A few seconds of margin in case the k6 and service clocks differ slightly.
  return { admin, user, productId, stockBefore, startedAt: Date.now() - 5000 };
}

export default function (data) {
  const items = [{ productId: data.productId, quantity: 1 }];
  const res =
    MODE === "guest"
      ? http.post(
          `${BASE_URL}/order/guest`,
          JSON.stringify({ name: `K6 Guest ${__VU}`, email: `k6-guest-${__VU}@example.com`, data: items }),
          { headers: jsonHeaders, tags: { name: "place order" } }
        )
      : http.post(`${BASE_URL}/order`, JSON.stringify({ data: items }), {
          ...auth(data.user.token),
          tags: { name: "place order" },
        });

  const ok = check(res, {
    "order placed": (r) => r.status === 200 && r.json("data.id") > 0,
  });
  if (ok) {
    placed.add(1);
  } else {
    rejected.add(1);
    // Logs about 20 of the failures, not one line per order.
    if (Math.random() < 20 / ORDERS) {
      console.warn(`order failed: ${res.status} ${String(res.body).slice(0, 200)}`);
    }
  }
}

export function teardown(data) {
  const { admin, productId } = data;

  // This run's orders: created since setup by k6 accounts (k6-...@example.com).
  const all = http.get(`${BASE_URL}/order?limit=0`, auth(admin)).json("data") || [];
  const recent = all.filter((o) => Date.parse(o.create_at) >= data.startedAt);
  const emails = {};
  for (const id of new Set(recent.map((o) => o.userId))) {
    const res = http.get(`${BASE_URL}/user/${id}`);
    emails[id] = res.status === 200 ? res.json("data.email") : "";
  }
  const mine = recent.filter((o) => (emails[o.userId] || "").startsWith("k6-"));
  const active = mine.filter((o) => o.state !== CANCELLED);

  // Each order that went through reserved exactly one unit.
  const reserved = data.stockBefore + ORDERS - getStock(admin, productId);
  console.log(
    `orders created: ${mine.length} (${active.length} open, ${mine.length - active.length} cancelled by the service), units reserved: ${reserved}`
  );
  check(null, { "stock matches orders": () => reserved === active.length });

  for (const o of active) {
    http.del(`${BASE_URL}/order/${o.orderId}`, null, auth(admin));
  }
  check(null, {
    "cancelling releases the stock": () => getStock(admin, productId) === data.stockBefore + ORDERS,
  });

  if (data.user) {
    http.del(`${BASE_URL}/user/${data.user.id}`, null, auth(admin));
  }
}
