const baseUrl = import.meta.env.VITE_API_URL

// newTransactionId makes the id of one user action. Pass it to each apiFetch call of that
// action, so the backend logs them all under it (Kibana: transaction_id : "<id>").
export const newTransactionId = (): string => crypto.randomUUID()

// apiFetch calls the gateway like fetch, sending transactionId (a new one for this request if
// none is given) in the X-Transaction-Id header. A failed request is logged to the console with
// its id, so it can be looked up in Kibana.
export const apiFetch = async (
  path: string,
  init: RequestInit = {},
  transactionId: string = newTransactionId(),
): Promise<Response> => {
  const headers = new Headers(init.headers)
  headers.set("X-Transaction-Id", transactionId)
  const request = `${init.method ?? "GET"} ${path}`

  let res: Response
  try {
    res = await fetch(baseUrl + path, { ...init, headers })
  } catch (error) {
    console.warn(`${request} failed, transaction id ${transactionId}:`, error)
    throw error
  }
  if (!res.ok) {
    console.warn(`${request} returned ${res.status}, transaction id ${transactionId}`)
  }
  return res
}
