"use client"

import { useState } from "react"
import { ArrowDownLeft, ArrowUpRight, Check, Copy, TrendingUp, Wallet } from "lucide-react"
import { RECENT_TX, USER, formatQtc, truncate } from "@/lib/quito-data"

export function DashboardView() {
  const [copied, setCopied] = useState(false)

  const copy = () => {
    navigator.clipboard?.writeText(USER.publicCode)
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
        {/* Hero balance */}
        <div className="relative overflow-hidden rounded-2xl border border-zinc-800 bg-zinc-900 p-6 md:col-span-2">
          <div className="absolute -right-16 -top-16 size-48 rounded-full bg-yellow-400/10 blur-3xl" />
          <div className="relative">
            <div className="flex items-center gap-2 text-sm text-zinc-400">
              <Wallet className="size-4 text-yellow-400" />
              Boa tarde, {USER.name}!
            </div>
            <p className="mt-6 text-xs uppercase tracking-widest text-zinc-500">
              Saldo disponível
            </p>
            <div className="mt-1 flex items-baseline gap-2">
              <span className="text-5xl font-bold tracking-tight text-yellow-400">
                {formatQtc(USER.balance)}
              </span>
              <span className="text-lg font-semibold text-zinc-500">QTC</span>
            </div>
            <div className="mt-4 inline-flex items-center gap-1.5 rounded-full bg-zinc-800 px-3 py-1 text-xs text-zinc-300">
              <TrendingUp className="size-3.5 text-emerald-400" />
              +12,4% esta semana
            </div>
          </div>
        </div>

        {/* Identity */}
        <div className="rounded-2xl border border-zinc-800 bg-zinc-900 p-6">
          <p className="text-xs uppercase tracking-widest text-zinc-500">Código público</p>
          <p className="mt-3 text-sm text-zinc-400">Sua identidade na rede</p>
          <div className="mt-4 flex items-center gap-2 rounded-xl border border-zinc-800 bg-zinc-950/60 p-2 pl-3">
            <span className="flex-1 truncate font-mono text-xs text-zinc-300">
              {USER.publicCode}
            </span>
            <button
              type="button"
              onClick={copy}
              aria-label="Copiar código público"
              className="flex size-8 shrink-0 items-center justify-center rounded-lg bg-zinc-800 text-zinc-300 transition-colors hover:bg-yellow-400 hover:text-zinc-950"
            >
              {copied ? <Check className="size-4" /> : <Copy className="size-4" />}
            </button>
          </div>
          <p className="mt-3 font-mono text-[11px] text-zinc-600">
            {copied ? "Copiado para a área de transferência" : "Toque para copiar"}
          </p>
        </div>

        {/* Recent activity */}
        <div className="rounded-2xl border border-zinc-800 bg-zinc-900 p-6 md:col-span-3">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-sm font-semibold text-zinc-100">Atividade recente</h2>
            <span className="text-xs text-zinc-500">Últimas 3 transações</span>
          </div>
          <ul className="flex flex-col divide-y divide-zinc-800">
            {RECENT_TX.map((tx) => {
              const received = tx.direction === "received"
              return (
                <li key={tx.id} className="flex items-center gap-4 py-3">
                  <span
                    className={`flex size-10 items-center justify-center rounded-full ${
                      received ? "bg-emerald-500/10" : "bg-zinc-800"
                    }`}
                  >
                    {received ? (
                      <ArrowDownLeft className="size-5 text-emerald-400" />
                    ) : (
                      <ArrowUpRight className="size-5 text-zinc-400" />
                    )}
                  </span>
                  <div className="min-w-0 flex-1">
                    <p className="text-sm font-medium text-zinc-100">
                      {received ? "Recebido" : "Enviado"}
                    </p>
                    <p className="truncate font-mono text-xs text-zinc-500">
                      {received ? truncate(tx.from) : truncate(tx.to)}
                    </p>
                  </div>
                  <div className="text-right">
                    <p
                      className={`text-sm font-semibold ${
                        received ? "text-emerald-400" : "text-zinc-300"
                      }`}
                    >
                      {received ? "+" : "−"}
                      {formatQtc(tx.amount)}
                    </p>
                    <p className="text-xs text-zinc-600">{tx.timestamp}</p>
                  </div>
                </li>
              )
            })}
          </ul>
        </div>
      </div>
    </div>
  )
}
