// The body of GET /api/status, shared by the server and the page.

export type TargetResult = {
  name: string
  group: string
  optional: boolean
  // The target URL with its password, if any, replaced by ***.
  url: string
  up: boolean
  // False when nothing answered (refused, closed, timed out); an optional target in that
  // state just isn't running.
  reachable: boolean
  latencyMs: number
  error?: string
}

export type StatusResponse = {
  checkedAt: string
  results: TargetResult[]
}
