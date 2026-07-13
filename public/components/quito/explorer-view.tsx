"use client"

import { useEffect, useRef, useState } from "react"
import { Box, Search } from "lucide-react"
import { BLOCKS, RICH_LIST, USER, formatQtc, randomHex, truncate } from "@/lib/quito-data"

type LogLine = { id: number; time: string; text: string; kind: "info" | "block" | "tx" }

function nowTime() {
  return new Date().toLocaleTimeString("pt-BR", { hour12: false })
}

const LOG_TEMPLATES: { text: () => string; kind: LogLine["kind"] }[] = [
  { text: () => `Bloco #${43 + Math.floor(Math.random() * 20)} validado`, kind: "block" },
  { text: () => `Nova transação 0x${randomHex(8)} no mempool`, kind: "tx" },
  { text: () => `Peer QTC-${randomHex(4).toUpperCase()} conectado`, kind: "info" },
  { text: () => `Hash confirmado 0x${randomHex(10)}`, kind: "info" },
  { text: () => `Recompensa de +50 QTC distribuída`, kind: "block" },
  { text: () => `Sincronizando ${Math.floor(Math.random() * 100)}% da cadeia`, kind: "info" },
]

export function ExplorerView() {
  return (
    <div className="mx-auto max-w-6xl">
      <header className="mb-8">
        <p className="text-sm text-zinc-500">Rede</p>
        <h1 className="text-2xl font-semibold tracking-tight text-zinc-50">
          Explorador da Blockchain
        </h1>
      </header>

      {/* Block timeline */}
      <section className="mb-6">
        <div className="mb-3 flex items-center gap-2">
          <Box className="size-4 text-yellow-400" />
          <h2 className="text-sm font-semibold text-zinc-100">Blocos recentes</h2>
        </div>
        <div className="flex flex-nowrap gap-3 overflow-x-auto pb-3">
          {BLOCKS.map((b) => (
            <div
              key={b.number}
              className="group w-52 shrink-0 rounded-2xl border border-zinc-800 bg-zinc-900 p-4 transition-colors hover:border-yellow-400/50"
            >
              <div className="flex items-center justify-between">
                <span className="text-lg font-bold text-yellow-400">#{b.number}</span>
                <span className="rounded-full bg-zinc-800 px-2 py-0.5 text-[10px] text-zinc-400">
                  {b.txCount} tx
                </span>
              </div>
              <p className="mt-3 text-[10px] uppercase tracking-widest text-zinc-600">Hash</p>
              <p className="truncate font-mono text-xs text-zinc-300">{b.hash}</p>
              <p className="mt-2 text-[10px] uppercase tracking-widest text-zinc-600">Minerador</p>
              <p className="truncate font-mono text-xs text-zinc-400">{b.miner}</p>
              <p className="mt-3 font-mono text-[11px] text-zinc-600">{b.timestamp}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Split: rich list + logs */}
      <div className="grid grid-cols-1 items-stretch gap-4 lg:grid-cols-2">
        <RichList />
        <NetworkLogs />
      </div>
    </div>
  )
}

function RichList() {
  const [query, setQuery] = useState("")
  const filtered = RICH_LIST.filter((r) =>
    r.code.toLowerCase().includes(query.toLowerCase()),
  )

  return (
    <div className="rounded-2xl border border-zinc-800 bg-zinc-900 p-5">
      <h2 className="mb-3 text-sm font-semibold text-zinc-100">Top 10 — Rich List</h2>
      <div className="relative mb-4">
        <Search className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-zinc-500" />
        <input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Buscar saldo por código"
          className="w-full rounded-xl border border-zinc-700 bg-zinc-950/60 py-2.5 pl-10 pr-3 font-mono text-xs text-zinc-100 placeholder:text-zinc-600 outline-none transition-shadow focus:border-yellow-400 focus:ring-2 focus:ring-yellow-400/30"
        />
      </div>
      <ul className="flex flex-col gap-1">
        {filtered.map((r) => {
          const isUser = r.code === USER.publicCode
          return (
            <li
              key={r.code}
              className={`flex items-center gap-3 rounded-lg px-2 py-2 ${
                isUser ? "bg-yellow-400/10 ring-1 ring-yellow-400/30" : ""
              }`}
            >
              <span
                className={`flex size-6 items-center justify-center rounded-md text-xs font-bold ${
                  r.rank <= 3 ? "bg-yellow-400 text-zinc-950" : "bg-zinc-800 text-zinc-400"
                }`}
              >
                {r.rank}
              </span>
              <span className="flex-1 truncate font-mono text-xs text-zinc-300">
                {truncate(r.code, 10, 4)}
                {isUser && <span className="ml-2 text-[10px] text-yellow-400">você</span>}
              </span>
              <span className="font-mono text-xs font-semibold text-zinc-100">
                {formatQtc(r.balance)}
              </span>
            </li>
          )
        })}
        {filtered.length === 0 && (
          <li className="py-6 text-center text-xs text-zinc-600">Nenhum código encontrado.</li>
        )}
      </ul>
    </div>
  )
}

function NetworkLogs() {
  const [logs, setLogs] = useState<LogLine[]>([])
  const idRef = useRef(0)
  const scrollRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const push = () => {
      const t = LOG_TEMPLATES[Math.floor(Math.random() * LOG_TEMPLATES.length)]
      setLogs((prev) => {
        const next = [...prev, { id: idRef.current++, time: nowTime(), text: t.text(), kind: t.kind }]
        return next.slice(-60)
      })
    }
    push()
    const interval = setInterval(push, 1800)
    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    // Scroll only the log container to the bottom — never the page.
    const el = scrollRef.current
    if (el) el.scrollTop = el.scrollHeight
  }, [logs])

  return (
    <div className="flex h-full min-h-[420px] flex-col overflow-hidden rounded-2xl border border-zinc-800 bg-black/50">
      <div className="flex shrink-0 items-center gap-2 border-b border-zinc-800 px-4 py-3">
        <span className="relative flex size-2.5">
          <span className="absolute inline-flex size-full animate-ping rounded-full bg-emerald-400/70" />
          <span className="relative inline-flex size-2.5 rounded-full bg-emerald-400" />
        </span>
        <span className="text-xs font-semibold text-zinc-300">Logs da rede</span>
        <span className="ml-auto font-mono text-[10px] uppercase tracking-widest text-zinc-600">
          websocket · live
        </span>
      </div>
      <div ref={scrollRef} className="min-h-0 flex-1 overflow-y-auto p-4 font-mono text-xs leading-relaxed">
        {logs.map((l) => (
          <div key={l.id} className="flex gap-2">
            <span className="shrink-0 text-zinc-600">[{l.time}]</span>
            <span
              className={
                l.kind === "block"
                  ? "text-yellow-400"
                  : l.kind === "tx"
                    ? "text-emerald-400"
                    : "text-zinc-500"
              }
            >
              {l.text}
            </span>
          </div>
        ))}
      </div>
    </div>
  )
}
