export type View = "auth" | "dashboard" | "transfer" | "explorer" | "mining" | "account"

export function formatQtc(value: number | undefined | null, digits = 2): string {
  const n = typeof value === "number" && Number.isFinite(value) ? value : 0
  return n.toLocaleString("en-US", {
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
  return d.toLocaleString("en-US", { dateStyle: "short", timeStyle: "short" })
}

export function relativeTime(iso: string | undefined): string {
  if (!iso) return "—"
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return "—"
  const mins = Math.floor((Date.now() - d.getTime()) / 60000)
  if (mins < 1) return "now"
  if (mins < 60) return `${mins} min ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours} h ago`
  return `${Math.floor(hours / 24)} d ago`
}
