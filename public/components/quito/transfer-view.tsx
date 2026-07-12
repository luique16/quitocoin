"use client"

import { useState } from "react"
import { Clock, CheckCircle2, Send } from "lucide-react"
import { TRANSFERS, USER, formatQtc, truncate } from "@/lib/quito-data"

export function TransferView() {
  const [recipient, setRecipient] = useState("")
  const [amount, setAmount] = useState("")

  const pending = TRANSFERS.filter((t) => t.status === "pending")
  const completed = TRANSFERS.filter((t) => t.status === "completed")

  return (
    <div className="mx-auto max-w-5xl">
      <header className="mb-8">
        <p className="text-sm text-zinc-500">Carteira</p>
        <h1 className="text-2xl font-semibold tracking-tight text-zinc-50">Transferir</h1>
      </header>

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        {/* Action form */}
        <div className="rounded-2xl border border-zinc-800 bg-zinc-900 p-6">
          <h2 className="text-sm font-semibold text-zinc-100">Enviar QuitoCoins</h2>
          <p className="mt-1 text-xs text-zinc-500">
            Saldo disponível: <span className="text-yellow-400">{formatQtc(USER.balance)} QTC</span>
          </p>

          <form
            className="mt-6 flex flex-col gap-5"
            onSubmit={(e) => {
              e.preventDefault()
              setRecipient("")
              setAmount("")
            }}
          >
            <label className="block">
              <span className="mb-1.5 block text-xs font-medium text-zinc-400">
                Código público do destinatário
              </span>
              <input
                value={recipient}
                onChange={(e) => setRecipient(e.target.value)}
                placeholder="QTC-XXXX-XXXX-XXXX-XXXX"
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
                  onClick={() => setAmount(String(USER.balance))}
                  className="absolute right-2 top-1/2 -translate-y-1/2 rounded-lg bg-zinc-800 px-2.5 py-1 text-xs font-semibold text-yellow-400 transition-colors hover:bg-zinc-700"
                >
                  MÁX
                </button>
              </div>
            </label>

            <button
              type="submit"
              className="mt-2 flex items-center justify-center gap-2 rounded-xl bg-yellow-400 py-3 text-sm font-semibold text-zinc-950 shadow-[0_0_24px_-6px] shadow-yellow-400/60 transition-all hover:bg-yellow-300 hover:shadow-yellow-400/80"
            >
              <Send className="size-4" />
              Enviar QuitoCoins
            </button>
          </form>
        </div>

        {/* Receipts */}
        <div className="flex flex-col gap-4">
          <ReceiptGroup
            title="Pendentes"
            icon={<Clock className="size-4 text-yellow-400" />}
            items={pending}
            accent="pending"
          />
          <ReceiptGroup
            title="Concluídas"
            icon={<CheckCircle2 className="size-4 text-emerald-400" />}
            items={completed}
            accent="completed"
          />
        </div>
      </div>
    </div>
  )
}

function ReceiptGroup({
  title,
  icon,
  items,
  accent,
}: {
  title: string
  icon: React.ReactNode
  items: typeof TRANSFERS
  accent: "pending" | "completed"
}) {
  return (
    <div className="rounded-2xl border border-zinc-800 bg-zinc-900 p-5">
      <div className="mb-3 flex items-center gap-2">
        {icon}
        <h3 className="text-sm font-semibold text-zinc-100">{title}</h3>
        <span className="ml-auto rounded-full bg-zinc-800 px-2 py-0.5 text-xs text-zinc-400">
          {items.length}
        </span>
      </div>
      <ul className="flex flex-col gap-2">
        {items.map((tx) => (
          <li
            key={tx.id}
            className="relative flex items-center justify-between gap-3 rounded-xl border border-dashed border-zinc-800 bg-zinc-950/40 p-3"
          >
            <div className="min-w-0">
              <p className="truncate font-mono text-xs text-zinc-300">{truncate(tx.to)}</p>
              <p className="text-[11px] text-zinc-600">{tx.timestamp}</p>
            </div>
            <div className="text-right">
              <p className="text-sm font-semibold text-zinc-100">−{formatQtc(tx.amount)}</p>
              <span
                className={`text-[10px] uppercase tracking-wide ${
                  accent === "pending" ? "text-yellow-400" : "text-emerald-400"
                }`}
              >
                {accent === "pending" ? "Confirmando" : "Confirmado"}
              </span>
            </div>
          </li>
        ))}
      </ul>
    </div>
  )
}
