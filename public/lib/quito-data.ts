export type View = "auth" | "dashboard" | "transfer" | "explorer" | "mining" | "account"

export function formatQtc(value: number | undefined | null, digits = 2): string {
  const n = typeof value === "number" && Number.isFinite(value) ? value : 0
  return n.toLocaleString("pt-BR", {
    minimumFractionDigits: digits,
    maximumFractionDigits: digits,
  })
}

export function truncate(code: string | undefined, head = 6, tail = 4): string {
  if (!code) return "—"
  if (code.length <= head + tail) return code
  return `${code.slice(0, head)}…${code.slice(-tail)}`
}

export function formatDateTime(iso: string | undefined): string {
  if (!iso) return "—"
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return "—"
  return d.toLocaleString("pt-BR", { dateStyle: "short", timeStyle: "short" })
}

export function relativeTime(iso: string | undefined): string {
  if (!iso) return "—"
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return "—"
  const mins = Math.floor((Date.now() - d.getTime()) / 60000)
  if (mins < 1) return "agora"
  if (mins < 60) return `${mins} min atrás`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours} h atrás`
  return `${Math.floor(hours / 24)} d atrás`
}
