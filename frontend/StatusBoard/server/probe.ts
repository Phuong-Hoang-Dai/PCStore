import http from 'node:http'
import https from 'node:https'
import net from 'node:net'
import type { Target } from './targets'
import type { TargetResult } from './types'

const TIMEOUT_MS = 3000

// reachable: something answered, just not the way a healthy server does.
class ProbeError extends Error {
  reachable: boolean

  constructor(message: string, reachable: boolean) {
    super(message)
    this.reachable = reachable
  }
}

// checkAll checks every target in parallel, so one slow target doesn't hold up the rest.
export function checkAll(targets: Target[]): Promise<TargetResult[]> {
  return Promise.all(targets.map(check))
}

async function check(target: Target): Promise<TargetResult> {
  const url = new URL(target.url)
  const start = performance.now()
  let error: ProbeError | undefined
  try {
    await probe(url)
  } catch (err) {
    error = toProbeError(err)
  }

  return {
    name: target.name,
    group: target.group,
    optional: target.optional ?? false,
    url: redact(url),
    up: !error,
    reachable: !error || error.reachable,
    latencyMs: Math.round(performance.now() - start),
    error: error?.message,
  }
}

function probe(url: URL): Promise<void> {
  switch (url.protocol) {
    case 'http:':
    case 'https:':
      return probeHttp(url)
    case 'postgres:':
      return probePostgres(url)
    case 'mongodb:':
      return probeMongo(url)
    case 'redis:':
      return probeRedis(url)
  }
  return Promise.reject(new ProbeError(`unsupported scheme ${url.protocol}`, false))
}

// Not fetch: it refuses the ports browsers block, and user_service listens on one (6666).
// node:http sends user:password from the URL as basic auth. agent: false opens a new
// connection per check, so a kept-alive socket to a since-restarted container can't fail it.
function probeHttp(url: URL): Promise<void> {
  return new Promise((resolve, reject) => {
    const req = (url.protocol === 'https:' ? https : http).get(url, { agent: false }, (res) => {
      res.resume()
      const status = res.statusCode ?? 0
      if (status >= 400) {
        reject(new ProbeError(`HTTP ${status} ${res.statusMessage ?? ''}`.trim(), true))
      } else {
        resolve()
      }
    })
    req.setTimeout(TIMEOUT_MS, () =>
      req.destroy(new ProbeError(`timed out after ${TIMEOUT_MS / 1000}s`, false)),
    )
    req.once('error', reject)
  })
}

async function probeRedis(url: URL) {
  const reply = await exchange(url, Buffer.from('PING\r\n'), (r) => r.includes('\n'))
  const line = reply.toString().trim()
  if (line !== '+PONG') {
    throw new ProbeError(`unexpected reply ${JSON.stringify(line)}`, true)
  }
}

// SSLRequest is the one message Postgres answers (with a single byte, S or N) before login.
const sslRequest = Buffer.from([0, 0, 0, 8, 0x04, 0xd2, 0x16, 0x2f])

async function probePostgres(url: URL) {
  const reply = await exchange(url, sslRequest, (r) => r.length >= 1)
  const answer = String.fromCharCode(reply[0])
  if (answer !== 'S' && answer !== 'N') {
    throw new ProbeError('unexpected reply to SSLRequest', true)
  }
}

const OP_MSG = 2013

function int32(n: number) {
  const b = Buffer.alloc(4)
  b.writeInt32LE(n)
  return b
}

// helloRequest is an OP_MSG running {hello: 1} on admin, which mongod answers without login.
function helloRequest() {
  const elements = Buffer.concat([
    Buffer.from([0x10]), Buffer.from('hello\0'), int32(1),
    Buffer.from([0x02]), Buffer.from('$db\0'), int32('admin'.length + 1), Buffer.from('admin\0'),
  ])
  const doc = Buffer.concat([int32(4 + elements.length + 1), elements, Buffer.from([0])])
  const body = Buffer.concat([int32(0), Buffer.from([0]), doc])
  return Buffer.concat([int32(16 + body.length), int32(1), int32(0), int32(OP_MSG), body])
}

async function probeMongo(url: URL) {
  const reply = await exchange(url, helloRequest(), (r) => r.length >= 16 && r.length >= r.readInt32LE(0))
  if (reply.readInt32LE(12) !== OP_MSG) {
    throw new ProbeError('unexpected reply to hello', true)
  }
  // The reply document holds ok: 1.0 (a BSON double) on success.
  const ok = reply.indexOf(Buffer.from([0x01, 0x6f, 0x6b, 0x00]))
  if (ok === -1 || reply.readDoubleLE(ok + 4) !== 1) {
    throw new ProbeError('hello failed', true)
  }
}

// exchange sends request on a new connection and resolves with the reply once done says it
// is complete.
function exchange(url: URL, request: Buffer, done: (reply: Buffer) => boolean): Promise<Buffer> {
  return new Promise((resolve, reject) => {
    const socket = net.connect({ host: url.hostname, port: Number(url.port) })
    let reply = Buffer.alloc(0)

    const timer = setTimeout(() => {
      socket.destroy()
      reject(new ProbeError(`timed out after ${TIMEOUT_MS / 1000}s`, false))
    }, TIMEOUT_MS)

    socket.once('connect', () => socket.write(request))
    socket.on('data', (chunk) => {
      reply = Buffer.concat([reply, chunk])
      if (done(reply)) {
        resolve(reply)
        socket.destroy()
      }
    })
    socket.once('error', reject)
    // Docker Desktop accepts connections on a published port even when nothing listens inside
    // the container, then closes them, so a plain TCP connect would count as up.
    socket.once('close', () => {
      clearTimeout(timer)
      reject(new ProbeError('connection closed without a reply', false))
    })
  })
}

function toProbeError(err: unknown): ProbeError {
  if (err instanceof ProbeError) {
    return err
  }
  switch ((err as NodeJS.ErrnoException)?.code) {
    case 'ECONNREFUSED':
      return new ProbeError('connection refused', false)
    case 'ECONNRESET':
      return new ProbeError('connection closed without a reply', false)
    case 'ENOTFOUND':
      return new ProbeError('host not found', false)
  }
  return new ProbeError(err instanceof Error ? err.message : String(err), false)
}

function redact(url: URL) {
  if (!url.password) {
    return url.href
  }
  const copy = new URL(url)
  copy.password = '***'
  return copy.href
}
