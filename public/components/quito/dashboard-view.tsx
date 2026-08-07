"use client"

import { useState } from "react"
import useSWR from "swr"
import { ArrowDownLeft, ArrowUpRight, Check, Copy, Clock, Pickaxe, Wallet } from "lucide-react"
import { api, type HistoryRecord } from "@/lib/api"
import { formatQtc, relativeTime, truncate } from "@/lib/quito-data"
import { useAuth } from "./auth-provider"
import { EmptyState, ErrorState, LoadingRows } from "./states"

type RecentItem = {
  key: string
  kind: "miner" | "sent" | "received"
  amount: number
  positive: boolean
  when?: string
  blockIndex: number
  otherParty?: string
}

function toRecent(records: HistoryRecord[] | null | undefined): RecentItem[] {
  if (!records) return []
  return records.map((r, i) => {
    const base = { key: `${r.block_index}-${r.type}-${i}`, blockIndex: r.block_index }
    if (r.type === "miner") {
      return { ...base, kind: "miner", amount: r.value ?? 0, positive: true, when: r.date }
    }
    if (r.type === "receiver") {
      return {
        ...base,
        kind: "received",
        amount: r.value ?? 0,
        positive: true,
        when: r.date,
        otherParty: r.other_party,
      }
    }
    return {
      ...base,
      kind: "sent",
      amount: r.value ?? 0,
      positive: false,
      when: r.date,
      otherParty: r.other_party,
    }
  })
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

  const recent = toRecent(history.data?.blocks).slice(0, 8)
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
            <EmptyState message="Nenhuma atividade ainda. Minere um bloco ou receba uma transferência." />
          ) : (
            <ul className="flex flex-col divide-y divide-zinc-800">
              {recent.map((item) => {
                const isReceived = item.kind === "received"
                const isMined = item.kind === "miner"
                return (
                  <li key={item.key} className="flex items-center gap-4 py-3">
                    <span
                      className={`flex size-10 shrink-0 items-center justify-center rounded-full ${
                        isReceived ? "bg-emerald-500/10" : isMined ? "bg-yellow-400/10" : "bg-zinc-800"
                      }`}
                    >
                      {isReceived ? (
                        <ArrowDownLeft className="size-5 text-emerald-400" aria-hidden="true" />
                      ) : isMined ? (
                        <Pickaxe className="size-5 text-yellow-400" aria-hidden="true" />
                      ) : (
                        <ArrowUpRight className="size-5 text-zinc-400" aria-hidden="true" />
                      )}
                    </span>
                    <div className="min-w-0 flex-1">
                      <p className="text-sm font-medium text-zinc-100">
                        {isMined ? "Bloco minerado" : isReceived ? "Recebido" : "Enviado"}
                      </p>
                      <p className="truncate font-mono text-xs text-zinc-500">
                        {isMined ? `Bloco #${item.blockIndex}` : truncate(item.otherParty)}
                      </p>
                    </div>
                    <div className="text-right">
                      <p
                        className={`text-sm font-semibold ${
                          isReceived ? "text-emerald-400" : isMined ? "text-yellow-400" : "text-zinc-300"
                        }`}
                      >
                        {item.positive ? "+" : "−"}
                        {formatQtc(item.amount)}
                      </p>
                      <p className="text-xs text-zinc-600">{relativeTime(item.when)}</p>
                    </div>
                  </li>
                )
              })}
            </ul>
          )}
        </div>
      </div>
    </div>
  )
}
