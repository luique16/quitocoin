export type View = "auth" | "dashboard" | "transfer" | "explorer" | "mining"

export type TxStatus = "pending" | "completed"
export type TxDirection = "received" | "sent"

export type Transaction = {
  id: string
  from: string
  to: string
  amount: number
  timestamp: string
  status: TxStatus
  direction: TxDirection
}

export type Block = {
  number: number
  hash: string
  miner: string
  timestamp: string
  txCount: number
}

export type RichEntry = {
  rank: number
  code: string
  balance: number
}

export const USER = {
  name: "João",
  publicCode: "QTC-8F3A-9C21-B7E0-4D6F",
  balance: 1250.0,
}

export const RECENT_TX: Transaction[] = [
  {
    id: "tx-1",
    from: "QTC-1A2B-3C4D-5E6F-7A8B",
    to: USER.publicCode,
    amount: 320.5,
    timestamp: "Hoje, 14:22",
    status: "completed",
    direction: "received",
  },
  {
    id: "tx-2",
    from: USER.publicCode,
    to: "QTC-9Z8Y-7X6W-5V4U-3T2S",
    amount: 75.0,
    timestamp: "Hoje, 11:08",
    status: "completed",
    direction: "sent",
  },
  {
    id: "tx-3",
    from: "QTC-4D5E-6F7A-8B9C-0D1E",
    to: USER.publicCode,
    amount: 500.0,
    timestamp: "Ontem, 19:47",
    status: "completed",
    direction: "received",
  },
]

export const TRANSFERS: Transaction[] = [
  {
    id: "tr-1",
    from: USER.publicCode,
    to: "QTC-A1B2-C3D4-E5F6-A7B8",
    amount: 120.0,
    timestamp: "Agora mesmo",
    status: "pending",
    direction: "sent",
  },
  {
    id: "tr-2",
    from: USER.publicCode,
    to: "QTC-B2C3-D4E5-F6A7-B8C9",
    amount: 45.75,
    timestamp: "Há 2 min",
    status: "pending",
    direction: "sent",
  },
  {
    id: "tr-3",
    from: USER.publicCode,
    to: "QTC-C3D4-E5F6-A7B8-C9D0",
    amount: 210.0,
    timestamp: "Há 18 min",
    status: "completed",
    direction: "sent",
  },
  {
    id: "tr-4",
    from: USER.publicCode,
    to: "QTC-D4E5-F6A7-B8C9-D0E1",
    amount: 88.2,
    timestamp: "Há 1 h",
    status: "completed",
    direction: "sent",
  },
]

export const BLOCKS: Block[] = Array.from({ length: 8 }).map((_, i) => {
  const number = 42 - i
  return {
    number,
    hash: `0x${randomHex(40)}`,
    miner: `QTC-${randomHex(4).toUpperCase()}-${randomHex(4).toUpperCase()}`,
    timestamp: `14:${String(59 - i * 3).padStart(2, "0")}:0${i % 10}`,
    txCount: 2 + (i % 4),
  }
})

export const RICH_LIST: RichEntry[] = [
  { rank: 1, code: "QTC-000A-CAFE-BABE-0001", balance: 98432.11 },
  { rank: 2, code: "QTC-1F2E-3D4C-5B6A-7089", balance: 74210.0 },
  { rank: 3, code: "QTC-DEAD-BEEF-1234-5678", balance: 51988.42 },
  { rank: 4, code: "QTC-8F3A-9C21-B7E0-4D6F", balance: 1250.0 },
  { rank: 5, code: "QTC-AA11-BB22-CC33-DD44", balance: 980.9 },
  { rank: 6, code: "QTC-5566-7788-99AA-BBCC", balance: 742.33 },
  { rank: 7, code: "QTC-0102-0304-0506-0708", balance: 611.0 },
  { rank: 8, code: "QTC-FEED-FACE-C0DE-BEEF", balance: 500.5 },
  { rank: 9, code: "QTC-1357-2468-9BDF-ACEB", balance: 321.75 },
  { rank: 10, code: "QTC-9182-7364-5546-3728", balance: 210.0 },
]

export const MEMPOOL: Transaction[] = [
  {
    id: "mp-1",
    from: "QTC-1A2B-3C4D-5E6F-7A8B",
    to: "QTC-9Z8Y-7X6W-5V4U-3T2S",
    amount: 42.0,
    timestamp: "Há 12 s",
    status: "pending",
    direction: "sent",
  },
  {
    id: "mp-2",
    from: "QTC-4D5E-6F7A-8B9C-0D1E",
    to: "QTC-A1B2-C3D4-E5F6-A7B8",
    amount: 128.5,
    timestamp: "Há 34 s",
    status: "pending",
    direction: "sent",
  },
  {
    id: "mp-3",
    from: "QTC-B2C3-D4E5-F6A7-B8C9",
    to: "QTC-C3D4-E5F6-A7B8-C9D0",
    amount: 7.25,
    timestamp: "Há 51 s",
    status: "pending",
    direction: "sent",
  },
]

export function randomHex(length: number): string {
  const chars = "0123456789abcdef"
  let out = ""
  for (let i = 0; i < length; i++) {
    out += chars[Math.floor(Math.random() * chars.length)]
  }
  return out
}

export function truncate(code: string, head = 6, tail = 4): string {
  if (code.length <= head + tail) return code
  return `${code.slice(0, head)}…${code.slice(-tail)}`
}

export function formatQtc(value: number): string {
  return value.toLocaleString("pt-BR", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
}
