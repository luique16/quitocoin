// QuitoCoin API HTTP client.
// Configurable base URL — the swagger points to localhost:4000.
export const API_URL = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080").replace(/\/$/, "")

export const TOKEN_KEY = "quito.token"
const TOKEN_MAX_AGE = 60 * 60 * 24 * 7

export function getToken(): string | null {
  if (typeof window === "undefined") return null
  const prefix = `${TOKEN_KEY}=`
  const found = document.cookie.split("; ").find((c) => c.startsWith(prefix))
  return found ? decodeURIComponent(found.slice(prefix.length)) : null
}

export function setToken(token: string) {
  document.cookie = `${TOKEN_KEY}=${encodeURIComponent(token)}; Path=/; Max-Age=${TOKEN_MAX_AGE}; SameSite=Lax`
}

export function clearToken() {
  document.cookie = `${TOKEN_KEY}=; Path=/; Max-Age=0; SameSite=Lax`
}

export class ApiError extends Error {
  status: number
  constructor(message: string, status: number) {
    super(message)
    this.name = "ApiError"
    this.status = status
  }
}

type Options = {
  method?: string
  body?: unknown
  auth?: boolean
}

async function request<T>(path: string, { method = "GET", body, auth = true }: Options = {}): Promise<T> {
  const headers: Record<string, string> = { "Content-Type": "application/json" }

  if (auth) {
    const token = getToken()
    if (!token) throw new ApiError("Session expired. Please sign in again.", 401)
    headers.Authorization = `Bearer ${token}`
  }

  let res: Response
  try {
    res = await fetch(`${API_URL}${path}`, {
      method,
      headers,
      body: body === undefined ? undefined : JSON.stringify(body),
    })
  } catch {
    throw new ApiError(
      `Could not connect to the API at ${API_URL}. Check that the server is running and that CORS is enabled.`,
      0,
    )
  }

  if (res.status === 204) return undefined as T

  const text = await res.text()
  let data: unknown = null
  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      data = null
    }
  }

  if (!res.ok) {
    const message =
      (data && typeof data === "object" && "error" in data && String((data as { error: unknown }).error)) ||
      `Erro ${res.status}`
    throw new ApiError(message, res.status)
  }

  return data as T
}

/* ---------------------------------- Tipos --------------------------------- */

export type User = {
  id: string
  public_id: string
  name: string
  email: string
  balance: number
  created_at: string
}

export type Transaction = {
  from?: string
  to?: string
  amount?: number
  created_at?: string
  [key: string]: unknown
}

export type Block = {
  index: number
  hash: string
  previous_hash: string
  nonce: number
  miner: string
  reward: number
  transactions?: Transaction[] | null
  tx_count?: number
  created_at: string
}

export type RankingEntry = {
  public_id: string
  balance: number
}

export type HistoryRecord = {
  type: "miner" | "sender" | "receiver"
  value: number
  date: string
  other_party: string
  block_index: number
}

export type NextBlock = {
  data: string
  mined: boolean
}

/* -------------------------------- Endpoints ------------------------------- */

export const api = {
  login: (email: string, password: string) =>
    request<{ token: string }>("/auth/login", { method: "POST", body: { email, password }, auth: false }),

  register: (name: string, email: string, password: string) =>
    request<{ token: string }>("/auth/register", {
      method: "POST",
      body: { name, email, password },
      auth: false,
    }),

  me: () => request<User>("/me"),

  updateMe: (input: { name: string; email: string }) => request<User>("/me", { method: "PUT", body: input }),

  updatePassword: (oldPassword: string, newPassword: string) =>
    request<{ message: string }>("/me/password", { method: "PUT", body: { oldPassword, newPassword } }),

  deleteMe: () => request<void>("/me", { method: "DELETE" }),

  ranking: () => request<{ richest: RankingEntry[] | null }>("/ranking"),

  blocks: (limit = 20, offset = 0) =>
    request<{ blocks: Block[] | null; total_count: number }>(`/blockchain/blocks?limit=${limit}&offset=${offset}`),

  block: (index: number) => request<{ block: Block }>(`/blockchain/blocks/${index}`),

  history: (role = "any", limit = 100) =>
    request<{ blocks: HistoryRecord[] | null }>(`/blockchain/history?role=${role}&limit=${limit}`),

  nextBlock: () => request<NextBlock>("/blockchain/next-block"),

  mine: (nonce: number) => request<Block & { reward: number }>("/blockchain/mine", { method: "POST", body: { nonce } }),

  transfer: (to: string, amount: number) =>
    request<{ from: string; to: string; amount: number }>("/transfer", { method: "POST", body: { to, amount } }),

  pending: () => request<{ transactions: Transaction[] | null }>("/transfer/pending"),

  pendingMe: () => request<{ transactions: Transaction[] | null }>("/transfer/pending/me"),
}
