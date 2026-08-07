// QuitoCoin proof of work.
// The backend validates: SHA-256(fmt.Sprintf("%d%s", nonce, data))
// That is: the decimal nonce comes FIRST, followed by the `data` string.

export const DIFFICULTY = Number(process.env.NEXT_PUBLIC_MINING_DIFFICULTY || 7)
export const TARGET_PREFIX = "0".repeat(DIFFICULTY)

const encoder = new TextEncoder()

export async function hashAttempt(nonce: number, data: string): Promise<string> {
  const buf = await crypto.subtle.digest("SHA-256", encoder.encode(`${nonce}${data}`))
  return Array.from(new Uint8Array(buf))
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("")
}

export type MineProgress = {
  nonce: number
  hash: string
  hashRate: number
  attempts: number
}

export type MineResult = {
  nonce: number
  hash: string
  attempts: number
  elapsedMs: number
}

/**
 * Searches for a nonce whose hash starts with `TARGET_PREFIX`.
 * Runs in slices to avoid freezing the UI and reports progress periodically.
 */
export async function mineBlock(
  data: string,
  opts: {
    signal?: { aborted: boolean }
    onProgress?: (p: MineProgress) => void
    batchSize?: number
  } = {},
): Promise<MineResult | null> {
  const { signal, onProgress, batchSize = 500 } = opts
  const start = performance.now()
  let nonce = 0
  let lastReport = start

  while (true) {
    if (signal?.aborted) return null

    for (let i = 0; i < batchSize; i++) {
      const hash = await hashAttempt(nonce, data)

      if (hash.startsWith(TARGET_PREFIX)) {
        onProgress?.({ nonce, hash, attempts: nonce + 1, hashRate: 0 })
        return { nonce, hash, attempts: nonce + 1, elapsedMs: performance.now() - start }
      }

      const now = performance.now()
      if (now - lastReport > 80) {
        lastReport = now
        onProgress?.({
          nonce,
          hash,
          attempts: nonce + 1,
          hashRate: Math.round((nonce + 1) / ((now - start) / 1000)),
        })
      }
      nonce++
    }

    // Yield control back to the browser between slices.
    await new Promise((r) => setTimeout(r, 0))
  }
}
