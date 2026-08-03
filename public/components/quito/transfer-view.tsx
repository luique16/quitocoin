"use client"

import { useState } from "react"
import useSWR from "swr"
import { CheckCircle2, Clock, Send } from "lucide-react"
import { api, type Block, type Transaction } from "@/lib/api"
import { formatQtc, relativeTime, truncate } from "@/lib/quito-data"
import { useAuth } from "./auth-provider"
import { EmptyState, ErrorState, LoadingRows, Spinner } from "./states"

export function TransferView() {
  const { user, refreshUser } = useAuth()
  const [recipient, setRecipient] = useState("")
  const [amount, setAmount] = useState("")
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const publicId = user?.public_id ?? ""
  const balance = user?.balance ?? 0

  const pending = useSWR(publicId ? "pending-me" : null, () => api.pendingMe(), {
    revalidateOnFocus: false,
  })
  const sentHistory = useSWR(publicId ? ["history", "sender"] : null, () => api.history("sender", 20), {
    revalidateOnFocus: false,
  })

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    setSuccess(null)

    const normalized = amount.replace(",", ".")
    const value = Number(normalized)

    if (!recipient.trim()) {
      setError("Informe o código público do destinatário.")
      return
    }
    if (!Number.isFinite(value) || value <= 0) {
      setError("Informe uma quantidade válida maior que zero.")
      return
    }
    if (value > balance) {
      setError("Saldo insuficiente para essa transferência.")
      return
    }
    if (recipient.trim() === publicId) {
      setError("Você não pode transferir para o seu próprio código.")
      return
    }

    setSubmitting(true)
    try {
      const res = await api.transfer(recipient.trim(), value)
      setSuccess(`Transferência de ${formatQtc(res.amount)} QTC enviada para o mempool.`)
      setRecipient("")
      setAmount("")
      await Promise.all([pending.mutate(), refreshUser()])
    } catch (err) {
      setError(err instanceof Error ? err.message : "Falha ao criar a transferência.")
    } finally {
      setSubmitting(false)
    }
  }

  const pendingItems = pending.data?.transactions ?? []
  const completedItems = flattenSent(sentHistory.data?.blocks, publicId)

  return (
    <div className="mx-auto max-w-5xl">
      <header className="mb-8">
        <p className="text-sm text-zinc-500">Carteira</p>
        <h1 className="text-2xl font-semibold tracking-tight text-zinc-50">Transferir</h1>
      </header>

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        {/* Formulário */}
        <div className="rounded-2xl border border-zinc-800 bg-zinc-900 p-6">
          <h2 className="text-sm font-semibold text-zinc-100">Enviar QuitoCoins</h2>
          <p className="mt-1 text-xs text-zinc-500">
            Saldo disponível: <span className="text-yellow-400">{formatQtc(balance)} QTC</span>
          </p>

          <form className="mt-6 flex flex-col gap-5" onSubmit={handleSubmit}>
            <label className="block">
              <span className="mb-1.5 block text-xs font-medium text-zinc-400">
                Código público do destinatário
              </span>
              <input
                value={recipient}
                onChange={(e) => setRecipient(e.target.value)}
                placeholder="Cole o código público"
                className="w-full rounded-xl border border-zinc-700 bg-zinc-950/60 px-3 py-2.5 font-mono text-sm text-zinc-100 placeholder:text-zinc-600 outline-none transition-shadow focus:border-yellow-400 focus:ring-2 focus:ring-yellow-400/30"
              />
            </label>

            <label className="block">
              <span className="mb-1.5 block text-xs font-medium text-zinc-400">Quantidade</span>
              <div className="relative">
                <input
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                  inputMode="decimal"
                  placeholder="0,00"
                  className="w-full rounded-xl border border-zinc-700 bg-zinc-950/60 py-2.5 pl-3 pr-20 text-sm text-zinc-100 placeholder:text-zinc-600 outline-none transition-shadow focus:border-yellow-400 focus:ring-2 focus:ring-yellow-400/30"
                />
                <button
                  type="button"
                  onClick={() => setAmount(String(balance))}
                  className="absolute right-2 top-1/2 -translate-y-1/2 rounded-lg bg-zinc-800 px-2.5 py-1 text-xs font-semibold text-yellow-400 transition-colors hover:bg-zinc-700"
                >
                  MÁX
                </button>
              </div>
            </label>

            {error && (
              <p
                role="alert"
                className="rounded-lg border border-red-900/50 bg-red-950/30 px-3 py-2 text-xs leading-relaxed text-red-200"
              >
                {error}
              </p>
            )}
            {success && (
              <p
                role="status"
                className="rounded-lg border border-emerald-900/50 bg-emerald-950/30 px-3 py-2 text-xs leading-relaxed text-emerald-200"
              >
                {success}
              </p>
            )}

            <button
              type="submit"
              disabled={submitting}
              className="mt-2 flex items-center justify-center gap-2 rounded-xl bg-yellow-400 py-3 text-sm font-semibold text-zinc-950 shadow-[0_0_24px_-6px] shadow-yellow-400/60 transition-all hover:bg-yellow-300 hover:shadow-yellow-400/80 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {submitting ? (
                <>
                  <Spinner />
                  Enviando…
                </>
              ) : (
                <>
                  <Send className="size-4" aria-hidden="true" />
                  Enviar QuitoCoins
                </>
              )}
            </button>
          </form>
        </div>

        {/* Recibos */}
        <div className="flex flex-col gap-4">
          <div className="rounded-2xl border border-zinc-800 bg-zinc-900 p-5">
            <div className="mb-3 flex items-center gap-2">
              <Clock className="size-4 text-yellow-400" aria-hidden="true" />
              <h3 className="text-sm font-semibold text-zinc-100">Pendentes</h3>
              <span className="ml-auto rounded-full bg-zinc-800 px-2 py-0.5 text-xs text-zinc-400">
                {pendingItems.length}
              </span>
            </div>
            {pending.isLoading ? (
              <LoadingRows rows={2} />
            ) : pending.error ? (
              <ErrorState
                message={pending.error instanceof Error ? pending.error.message : "Erro ao carregar."}
                onRetry={() => pending.mutate()}
              />
            ) : pendingItems.length === 0 ? (
              <EmptyState message="Nenhuma transferência aguardando confirmação." />
            ) : (
              <ul className="flex flex-col gap-2">
                {pendingItems.map((tx, i) => (
                  <Receipt key={i} tx={tx} publicId={publicId} accent="pending" />
                ))}
              </ul>
            )}
          </div>

          <div className="rounded-2xl border border-zinc-800 bg-zinc-900 p-5">
            <div className="mb-3 flex items-center gap-2">
              <CheckCircle2 className="size-4 text-emerald-400" aria-hidden="true" />
              <h3 className="text-sm font-semibold text-zinc-100">Enviadas e confirmadas</h3>
              <span className="ml-auto rounded-full bg-zinc-800 px-2 py-0.5 text-xs text-zinc-400">
                {completedItems.length}
              </span>
            </div>
            {sentHistory.isLoading ? (
              <LoadingRows rows={2} />
            ) : sentHistory.error ? (
              <ErrorState
                message={sentHistory.error instanceof Error ? sentHistory.error.message : "Erro ao carregar."}
                onRetry={() => sentHistory.mutate()}
              />
            ) : completedItems.length === 0 ? (
              <EmptyState message="Você ainda não enviou transferências confirmadas." />
            ) : (
              <ul className="flex flex-col gap-2">
                {completedItems.map((tx, i) => (
                  <Receipt key={i} tx={tx} publicId={publicId} accent="completed" />
                ))}
              </ul>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

function flattenSent(blocks: Block[] | null | undefined, publicId: string): Transaction[] {
  if (!blocks) return []
  const out: Transaction[] = []
  for (const block of blocks) {
    const txs = Array.isArray(block.transactions) ? block.transactions : []
    for (const tx of txs) {
      if (tx.from === publicId) out.push({ ...tx, created_at: tx.created_at ?? block.created_at })
    }
  }
  return out.slice(0, 10)
}

function Receipt({
  tx,
  publicId,
  accent,
}: {
  tx: Transaction
  publicId: string
  accent: "pending" | "completed"
}) {
  const outgoing = tx.from === publicId
  const counterparty = outgoing ? tx.to : tx.from
  const amount = typeof tx.amount === "number" ? tx.amount : 0

  return (
    <li className="flex items-center justify-between gap-3 rounded-xl border border-dashed border-zinc-800 bg-zinc-950/40 p-3">
      <div className="min-w-0">
        <p className="truncate font-mono text-xs text-zinc-300">{truncate(counterparty as string)}</p>
        <p className="text-[11px] text-zinc-600">
          {outgoing ? "Para" : "De"} · {relativeTime(tx.created_at)}
        </p>
      </div>
      <div className="text-right">
        <p className={`text-sm font-semibold ${outgoing ? "text-zinc-100" : "text-emerald-400"}`}>
          {outgoing ? "−" : "+"}
          {formatQtc(amount)}
        </p>
        <span
          className={`text-[10px] uppercase tracking-wide ${
            accent === "pending" ? "text-yellow-400" : "text-emerald-400"
          }`}
        >
          {accent === "pending" ? "Confirmando" : "Confirmado"}
        </span>
      </div>
    </li>
  )
}
