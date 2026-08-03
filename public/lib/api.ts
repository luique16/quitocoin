// Cliente HTTP da API QuitoCoin.
// Base URL configurável — o swagger aponta para localhost:4000.
export const API_URL = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080").replace(/\/$/, "")

const TOKEN_KEY = "quito.token"

export function getToken(): string | null {
  if (typeof window === "undefined") return null
  return window.localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string) {
  window.localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken() {
  window.localStorage.removeItem(TOKEN_KEY)
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
    if (!token) throw new ApiError("Sessão expirada. Faça login novamente.", 401)
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
      `Não foi possível conectar à API em ${API_URL}. Verifique se o servidor está rodando e se o CORS está habilitado.`,
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
    request<{ blocks: Block[] | null }>(`/blockchain/history?role=${role}&limit=${limit}`),

  nextBlock: () => request<NextBlock>("/blockchain/next-block"),

  mine: (nonce: number) => request<Block & { reward: number }>("/blockchain/mine", { method: "POST", body: { nonce } }),

  transfer: (to: string, amount: number) =>
    request<{ from: string; to: string; amount: number }>("/transfer", { method: "POST", body: { to, amount } }),

  pending: () => request<{ transactions: Transaction[] | null }>("/transfer/pending"),

  pendingMe: () => request<{ transactions: Transaction[] | null }>("/transfer/pending/me"),
}
