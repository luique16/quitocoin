"use client"

import { useEffect, useRef, useState } from "react"
import { Cpu, Layers, Pickaxe, Sparkles, Square } from "lucide-react"
import { MEMPOOL, formatQtc, randomHex, truncate } from "@/lib/quito-data"

type MiningState = "idle" | "mining" | "found"

export function MiningView() {
  const [state, setState] = useState<MiningState>("idle")
  const [nonce, setNonce] = useState(0)
  const [hash, setHash] = useState(`0x${"0".repeat(40)}`)
  const tickRef = useRef<ReturnType<typeof setInterval> | null>(null)
  const doneRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const clearTimers = () => {
    if (tickRef.current) clearInterval(tickRef.current)
    if (doneRef.current) clearTimeout(doneRef.current)
    tickRef.current = null
    doneRef.current = null
  }

  const startMining = () => {
    setState("mining")
    tickRef.current = setInterval(() => {
      setNonce(Math.floor(Math.random() * 9_999_999))
      setHash(`0x${randomHex(40)}`)
    }, 60)
    // Simulate finding a block after a short delay
    doneRef.current = setTimeout(() => {
      if (tickRef.current) clearInterval(tickRef.current)
      tickRef.current = null
      setNonce(Math.floor(Math.random() * 9_999_999))
      setHash(`0x0000${randomHex(36)}`)
      setState("found")
    }, 3800)
  }

  const stopMining = () => {
    clearTimers()
    setState("idle")
  }

  const reset = () => {
    clearTimers()
    setState("idle")
    setHash(`0x${"0".repeat(40)}`)
    setNonce(0)
  }

  useEffect(() => () => clearTimers(), [])

  return (
    <div className="mx-auto max-w-5xl">
      <header className="mb-8">
        <p className="text-sm text-zinc-500">Consenso</p>
        <h1 className="text-2xl font-semibold tracking-tight text-zinc-50">Mineração</h1>
      </header>

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-5">
        {/* Mempool */}
        <div className="rounded-2xl border border-zinc-800 bg-zinc-900 p-5 lg:col-span-2">
          <div className="mb-1 flex items-center gap-2">
            <Layers className="size-4 text-yellow-400" />
            <h2 className="text-sm font-semibold text-zinc-100">Mempool</h2>
          </div>
          <p className="mb-4 text-xs text-zinc-500">
            <span className="text-zinc-300">{MEMPOOL.length} / 3</span> tx prontas para o bloco
          </p>
          <ul className="flex flex-col gap-2">
            {MEMPOOL.map((tx) => (
              <li
                key={tx.id}
                className="rounded-xl border border-zinc-800 bg-zinc-950/40 p-3"
              >
                <div className="flex items-center justify-between">
                  <span className="font-mono text-xs text-zinc-400">{truncate(tx.from)}</span>
                  <span className="text-sm font-semibold text-zinc-100">
                    {formatQtc(tx.amount)}
                  </span>
                </div>
                <div className="mt-1 flex items-center justify-between">
                  <span className="font-mono text-[11px] text-zinc-600">
                    → {truncate(tx.to)}
                  </span>
                  <span className="text-[11px] text-yellow-400">{tx.timestamp}</span>
                </div>
              </li>
            ))}
          </ul>
        </div>

        {/* Console */}
        <div className="lg:col-span-3">
          <div
            className={`flex flex-col rounded-2xl border bg-zinc-900 p-6 transition-all duration-300 ${
              state === "found"
                ? "border-emerald-400/60 shadow-[0_0_40px_-10px] shadow-emerald-400/40"
                : state === "mining"
                  ? "border-yellow-400/50 shadow-[0_0_40px_-12px] shadow-yellow-400/40"
                  : "border-zinc-800"
            }`}
          >
            <div className="flex items-center gap-2 text-sm text-zinc-400">
              <Cpu className="size-4 text-yellow-400" />
              Console de mineração
            </div>

            {/* Big action button */}
            <div className="my-6 flex justify-center">
              {state !== "found" ? (
                <button
                  type="button"
                  onClick={state === "mining" ? stopMining : startMining}
                  className={`flex size-40 flex-col items-center justify-center gap-2 rounded-full text-sm font-semibold transition-all ${
                    state === "mining"
                      ? "animate-pulse bg-zinc-800 text-zinc-200 ring-4 ring-yellow-400/30"
                      : "bg-yellow-400 text-zinc-950 shadow-[0_0_40px_-6px] shadow-yellow-400/70 hover:scale-105 hover:bg-yellow-300"
                  }`}
                >
                  {state === "mining" ? (
                    <>
                      <Square className="size-7" />
                      Parar
                    </>
                  ) : (
                    <>
                      <Pickaxe className="size-8" />
                      Iniciar mineração
                    </>
                  )}
                </button>
              ) : (
                <button
                  type="button"
                  onClick={reset}
                  className="flex size-40 flex-col items-center justify-center gap-2 rounded-full bg-emerald-500 text-sm font-semibold text-zinc-950 shadow-[0_0_40px_-6px] shadow-emerald-400/70 transition-all hover:scale-105"
                >
                  <Sparkles className="size-8" />
                  Minerar novamente
                </button>
              )}
            </div>

            {/* Readout */}
            <div className="rounded-xl border border-zinc-800 bg-black/50 p-4 font-mono text-xs">
              <div className="flex items-center justify-between border-b border-zinc-800 pb-2">
                <span className="text-zinc-600">nonce</span>
                <span
                  className={
                    state === "mining" ? "text-yellow-400" : "text-zinc-300"
                  }
                >
                  {nonce.toLocaleString("pt-BR")}
                </span>
              </div>
              <div className="flex items-center justify-between gap-2 pt-2">
                <span className="shrink-0 text-zinc-600">hash</span>
                <span
                  className={`truncate ${
                    state === "found"
                      ? "text-emerald-400"
                      : state === "mining"
                        ? "text-yellow-400"
                        : "text-zinc-500"
                  }`}
                >
                  {hash}
                </span>
              </div>
            </div>

            {/* Success badge */}
            {state === "found" && (
              <div className="mt-4 flex items-center justify-center gap-2 rounded-xl border border-emerald-400/40 bg-emerald-500/10 py-3 text-sm font-semibold text-emerald-400">
                <Sparkles className="size-4" />
                Bloco encontrado! +50 QTC
              </div>
            )}
            {state === "idle" && (
              <p className="mt-4 text-center text-xs text-zinc-600">
                Pressione iniciar para começar a buscar um hash válido.
              </p>
            )}
            {state === "mining" && (
              <p className="mt-4 text-center text-xs text-zinc-500">
                Buscando hash com dificuldade alvo…
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
