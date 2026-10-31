import type { IncomingMessage, ServerResponse } from 'node:http'
import type { Plugin } from 'vite'
import { checkAll } from './probe'
import { targets } from './targets'
import type { StatusResponse } from './types'

// statusApi serves GET /api/status from the dev and preview servers, so the checks run on the
// host and not in the browser: a browser can't speak the database protocols, and the services
// don't send CORS headers for this page.
export function statusApi(): Plugin {
  async function handle(req: IncomingMessage, res: ServerResponse) {
    if (req.method !== 'GET') {
      res.statusCode = 405
      res.end()
      return
    }

    const body: StatusResponse = {
      checkedAt: new Date().toISOString(),
      results: await checkAll(targets),
    }
    res.setHeader('Content-Type', 'application/json')
    res.setHeader('Cache-Control', 'no-store')
    res.end(JSON.stringify(body))
  }

  return {
    name: 'status-api',
    configureServer(server) {
      server.middlewares.use('/api/status', handle)
    },
    configurePreviewServer(server) {
      server.middlewares.use('/api/status', handle)
    },
  }
}
