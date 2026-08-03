"use client"

import { useState } from "react"
import useSWR from "swr"
import { ArrowDownLeft, ArrowUpRight, Check, Copy, Clock, Pickaxe, Wallet } from "lucide-react"
import { api, type Block, type Transaction } from "@/lib/api"
import { formatQtc, relativeTime, truncate } from "@/lib/quito-data"
import { useAuth } from "./auth-provider"
import { EmptyState, ErrorState, LoadingRows } from "./states"

type FlatTx = {
  key: string
  from?: string
  to?: string
  amount: number
  received: boolean
  when?: string
}

function flatten(blocks: Block[] | null | undefined, publicId: string): FlatTx[] {
  if (!blocks) return []
  const out: FlatTx[] = []
  for (const block of blocks) {
    const txs = Array.isArray(block.transactions) ? block.transactions : []
    txs.forEach((tx: Transaction, i) => {
      const from = typeof tx.from === "string" ? tx.from : undefined
      const to = typeof tx.to === "string" ? tx.to : undefined
      if (from !== publicId && to !== publicId) return
      out.push({
        key: `${block.index}-${i}`,
        from,
        to,
        amount: typeof tx.amount === "number" ? tx.amount : 0,
        received: to === publicId,
        when: typeof tx.created_at === "string" ? tx.created_at : block.created_at,
      })
    })
  }
  return out
}

export function DashboardView() {
  const { user } = useAuth()
  const [copied, setCopied] = useState(false)

  const publicId = user?.public_id ?? ""

  const history = useSWR(publicId ? ["history", "any"] : null, () => api.history("any", 30), {
    revalidateOnFocus: false,
  })
  const mined = useSWR(publicId ? ["history", "miner"] : null, () => api.history("miner", 100), {
    revalidateOnFocus: false,
  })
  const pending = useSWR(publicId ? "pending-me" : null, () => api.pendingMe(), {
    revalidateOnFocus: false,
  })

  const recent = flatten(history.data?.blocks, publicId).slice(0, 5)
  const minedCount = mined.data?.blocks?.length ?? 0
  const pendingCount = pending.data?.transactions?.length ?? 0

  const copy = () => {
    if (!publicId) return
    navigator.clipboard?.writeText(publicId)
    setCopied(true)
    setTimeout(() => setCopied(false), 1600)
  }

  return (
    <div className="mx-auto max-w-5xl">
      <header className="mb-8">
        <p className="text-sm text-zinc-500">Visão geral</p>
        <h1 className="text-2xl font-semibold tracking-tight text-zinc-50">Dashboard</h1>
      </header>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
        {/* Saldo */}
        <div className="relative overflow-hidden rounded-2xl border border-zinc-800 bg-zinc-900 p-6 md:col-span-2">
          <div className="absolute -right-16 -top-16 size-48 rounded-full bg-yellow-400/10 blur-3xl" />
          <div className="relative">
            <div className="flex items-center gap-2 text-sm text-zinc-400">
              <Wallet className="size-4 text-yellow-400" aria-hidden="true" />
              Olá, {user?.name ?? "—"}!
            </div>
            <p className="mt-6 text-xs uppercase tracking-widest text-zinc-500">Saldo disponível</p>
            <div className="mt-1 flex items-baseline gap-2">
              <span className="text-5xl font-bold tracking-tight text-yellow-400">
                {formatQtc(user?.balance)}
              </span>
              <span className="text-lg font-semibold text-zinc-500">QTC</span>
            </div>
            <div className="mt-4 flex flex-wrap gap-2">
              <span className="inline-flex items-center gap-1.5 rounded-full bg-zinc-800 px-3 py-1 text-xs text-zinc-300">
                <Pickaxe className="size-3.5 text-yellow-400" aria-hidden="true" />
                {minedCount} {minedCount === 1 ? "bloco minerado" : "blocos minerados"}
              </span>
              <span className="inline-flex items-center gap-1.5 rounded-full bg-zinc-800 px-3 py-1 text-xs text-zinc-300">
                <Clock className="size-3.5 text-zinc-400" aria-hidden="true" />
                {pendingCount} {pendingCount === 1 ? "pendente" : "pendentes"}
              </span>
            </div>
          </div>
        </div>

        {/* Identidade */}
        <div className="rounded-2xl border border-zinc-800 bg-zinc-900 p-6">
          <p className="text-xs uppercase tracking-widest text-zinc-500">Código público</p>
          <p className="mt-3 text-sm text-zinc-400">Sua identidade na rede</p>
          <div className="mt-4 flex items-center gap-2 rounded-xl border border-zinc-800 bg-zinc-950/60 p-2 pl-3">
            <span className="flex-1 truncate font-mono text-xs text-zinc-300">{publicId || "—"}</span>
            <button
              type="button"
              onClick={copy}
              disabled={!publicId}
              aria-label="Copiar código público"
              className="flex size-8 shrink-0 items-center justify-center rounded-lg bg-zinc-800 text-zinc-300 transition-colors hover:bg-yellow-400 hover:text-zinc-950 disabled:opacity-50"
            >
              {copied ? <Check className="size-4" /> : <Copy className="size-4" />}
            </button>
          </div>
          <p className="mt-3 font-mono text-[11px] text-zinc-600">
            {copied ? "Copiado para a área de transferência" : "Toque para copiar"}
          </p>
        </div>

        {/* Atividade recente */}
        <div className="rounded-2xl border border-zinc-800 bg-zinc-900 p-6 md:col-span-3">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-sm font-semibold text-zinc-100">Atividade recente</h2>
            <span className="text-xs text-zinc-500">Confirmada na blockchain</span>
          </div>

          {history.isLoading ? (
            <LoadingRows rows={3} />
          ) : history.error ? (
            <ErrorState
              message={history.error instanceof Error ? history.error.message : "Erro ao carregar histórico."}
              onRetry={() => history.mutate()}
            />
          ) : recent.length === 0 ? (
            <EmptyState message="Nenhuma transação confirmada ainda. Minere um bloco ou receba uma transferência." />
          ) : (
            <ul className="flex flex-col divide-y divide-zinc-800">
              {recent.map((tx) => (
                <li key={tx.key} className="flex items-center gap-4 py-3">
                  <span
                    className={`flex size-10 shrink-0 items-center justify-center rounded-full ${
                      tx.received ? "bg-emerald-500/10" : "bg-zinc-800"
                    }`}
                  >
                    {tx.received ? (
                      <ArrowDownLeft className="size-5 text-emerald-400" aria-hidden="true" />
                    ) : (
                      <ArrowUpRight className="size-5 text-zinc-400" aria-hidden="true" />
                    )}
                  </span>
                  <div className="min-w-0 flex-1">
                    <p className="text-sm font-medium text-zinc-100">{tx.received ? "Recebido" : "Enviado"}</p>
                    <p className="truncate font-mono text-xs text-zinc-500">
                      {tx.received ? truncate(tx.from) : truncate(tx.to)}
                    </p>
                  </div>
                  <div className="text-right">
                    <p
                      className={`text-sm font-semibold ${tx.received ? "text-emerald-400" : "text-zinc-300"}`}
                    >
                      {tx.received ? "+" : "−"}
                      {formatQtc(tx.amount)}
                    </p>
                    <p className="text-xs text-zinc-600">{relativeTime(tx.when)}</p>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  )
}
