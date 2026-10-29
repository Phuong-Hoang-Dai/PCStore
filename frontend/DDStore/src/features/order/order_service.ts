import type { Order, OrderItemRequest } from "./order_model";

const url = import.meta.env.VITE_API_URL

const readError = async (res: Response) => {
  try {
    const data = await res.json()
    return data["error"] ?? `Request failed with status ${res.status}`
  } catch {
    return `Request failed with status ${res.status}`
  }
}

export const createOrder = async (items: OrderItemRequest[]): Promise<number> => {
  const res = await fetch(url + "/order", {
    method: "POST",
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ data: items }),
    credentials: "include",
  })
  if (!res.ok){
    throw new Error(await readError(res));
  }

  const data = await res.json()
  return data["data"]["id"]
}

export const createCheckoutSession = async (orderId: number): Promise<string> => {
  const res = await fetch(url + "/payment/checkout/" + orderId, {
    method: "POST",
    credentials: "include",
  })
  if (!res.ok){
    throw new Error(await readError(res));
  }

  const data = await res.json()
  return data["data"]["url"]
}

export const fetchOrderHistory = async (limit = 20, offset = 0): Promise<Order[]> => {
  const res = await fetch(url + `/order/history?limit=${limit}&offset=${offset}`, {
    method: "GET",
    credentials: "include",
  })
  if (!res.ok){
    throw new Error(await readError(res));
  }

  const data = await res.json()
  return data["data"] ?? []
}
