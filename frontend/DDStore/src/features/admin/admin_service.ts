import type { Order, OrderDetail } from "../order/order_model";
import type { AdminProduct, Category, ProductInput, Stock } from "./admin_model";

const url = import.meta.env.VITE_API_URL

// Calls the gateway with the session cookie and throws the gateway's error message.
const request = async <T>(path: string, init: RequestInit = {}): Promise<T> => {
  const res = await fetch(url + path, {
    ...init,
    headers: init.body ? { 'Content-Type': 'application/json' } : undefined,
    credentials: "include",
  })

  let data: unknown = null
  try {
    data = await res.json()
  } catch {
    // Empty or non-JSON body.
  }

  if (!res.ok) {
    // Some gateway errors serialize as {} instead of a string.
    const error = (data as { error?: unknown } | null)?.error
    throw new Error(typeof error === "string" && error ? error : `Request failed with status ${res.status}`)
  }
  return data as T
}

type Envelope<T> = { data: T }

const toProductBody = (input: ProductInput) => JSON.stringify({
  name: input.name,
  description: input.description,
  brand: input.brand,
  price: input.price,
  type: "standard",
  cate: { id: input.categoryId },
  images: input.imageUrl ? [{ url: input.imageUrl, role: "thumbnail" }] : [],
})

// ---------- products ----------

export const fetchProducts = async (limit: number, offset: number, categoryId?: number) => {
  const path = categoryId
    ? `/product/category/${categoryId}?limit=${limit}&offset=${offset}`
    : `/product?limit=${limit}&offset=${offset}`
  const res = await request<Envelope<AdminProduct[] | null>>(path)
  return res.data ?? []
}

export const fetchProduct = async (id: string) =>
  (await request<Envelope<AdminProduct>>(`/product/${id}`)).data

// Creation is asynchronous: the gateway queues the product and the catalog creates it.
export const createProduct = (input: ProductInput) =>
  request<unknown>("/product", { method: "POST", body: toProductBody(input) })

export const updateProduct = (id: string, input: ProductInput) =>
  request<unknown>(`/product/${id}`, { method: "PUT", body: toProductBody(input) })

export const deleteProduct = (id: string) =>
  request<unknown>(`/product/${id}`, { method: "DELETE" })

// ---------- stock ----------

export const fetchStocks = async (productIds: string[]): Promise<Record<string, number>> => {
  if (productIds.length === 0) {
    return {}
  }
  const res = await request<Envelope<Stock[] | null>>(`/stock?ids=${productIds.join(",")}`)
  return Object.fromEntries((res.data ?? []).map((s) => [s.productId, s.quantity]))
}

export const addStock = async (productId: string, quantity: number) =>
  (await request<Envelope<Stock>>(`/stock/${productId}/add`, {
    method: "POST",
    body: JSON.stringify({ quantity }),
  })).data.quantity

// ---------- categories ----------

export const fetchCategories = async () =>
  (await request<Envelope<Category[] | null>>("/category")).data ?? []

export const createCategory = (name: string) =>
  request<unknown>("/category", { method: "POST", body: JSON.stringify({ name }) })

export const renameCategory = (id: number, name: string) =>
  request<unknown>(`/category/${id}`, { method: "PUT", body: JSON.stringify({ name }) })

export const deleteCategory = (id: number) =>
  request<unknown>(`/category/${id}`, { method: "DELETE" })

// ---------- orders ----------

export const fetchOrders = async (limit: number, offset: number) =>
  (await request<Envelope<Order[] | null>>(`/order?limit=${limit}&offset=${offset}`)).data ?? []

export const fetchOrderDetail = async (id: number) =>
  (await request<Envelope<OrderDetail>>(`/order/${id}`)).data

export const completeOrder = (id: number) =>
  request<unknown>(`/order/completeOrder/${id}`, { method: "POST" })

export const cancelOrder = (id: number) =>
  request<unknown>(`/order/${id}`, { method: "DELETE" })
