// Everything the status board checks. It runs on the host, so it reaches the containers
// through the ports backend/DDStore/docker-compose.yaml publishes on 127.0.0.1.
//
// The URL scheme picks the check (see probe.ts):
//   http://, https://  GET, up below HTTP 400. user:password in the URL is sent as basic auth.
//   postgres://, mongodb://, redis://  the server answers a handshake in its own protocol.
//
// Optional targets only run with a compose profile (search, monitoring, db, tunnel), so they
// don't count toward the overall status, and show as "Off" when nothing answers.
export type Target = {
  name: string
  group: string
  url: string
  optional?: boolean
}

export const targets: Target[] = [
  { name: 'frontend', group: 'Services', url: 'http://127.0.0.1:5173/' },
  { name: 'api_gateway', group: 'Services', url: 'http://127.0.0.1:9999/healthcheck' },
  { name: 'user_service', group: 'Services', url: 'http://127.0.0.1:6666/user/healthcheck' },
  { name: 'product_service', group: 'Services', url: 'http://127.0.0.1:8888/product/healthcheck' },
  { name: 'order_service', group: 'Services', url: 'http://127.0.0.1:7777/order/healthcheck' },
  { name: 'search_service', group: 'Services', url: 'http://127.0.0.1:8889/health', optional: true },

  // 5499 and 27018, not the default ports: Windows already runs its own Postgres and MongoDB
  // on 5432 and 27017.
  { name: 'postgres', group: 'Infrastructure', url: 'postgres://127.0.0.1:5499' },
  { name: 'mongo', group: 'Infrastructure', url: 'mongodb://127.0.0.1:27018' },
  { name: 'redis', group: 'Infrastructure', url: 'redis://127.0.0.1:6379' },
  {
    name: 'rabbitmq',
    group: 'Infrastructure',
    url: 'http://ddstore:ddstore@127.0.0.1:15672/api/health/checks/alarms',
  },
  {
    name: 'elasticsearch',
    group: 'Infrastructure',
    // Answers 408 instead of 200 while the cluster is red.
    url: 'http://127.0.0.1:9200/_cluster/health?wait_for_status=yellow&timeout=1s',
    optional: true,
  },

  { name: 'kibana', group: 'Tools', url: 'http://127.0.0.1:5601/api/status', optional: true },
  { name: 'prometheus', group: 'Tools', url: 'http://127.0.0.1:9090/-/healthy', optional: true },
  { name: 'grafana', group: 'Tools', url: 'http://127.0.0.1:3000/api/health', optional: true },
  { name: 'pgadmin', group: 'Tools', url: 'http://127.0.0.1:5050/misc/ping', optional: true },
  { name: 'ngrok', group: 'Tools', url: 'http://127.0.0.1:4040/api/tunnels', optional: true },
]
