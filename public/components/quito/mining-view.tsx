"use client"

import { useCallback, useEffect, useRef, useState } from "react"
import useSWR, { useSWRConfig } from "swr"
import { Cpu, Layers, Pickaxe, RefreshCw, Sparkles, Square } from "lucide-react"
import { api, type Block } from "@/lib/api"
import { DIFFICULTY, TARGET_PREFIX, mineBlock } from "@/lib/mining"
import { formatQtc, relativeTime, truncate } from "@/lib/quito-data"
import { useAuth } from "./auth-provider"
import { EmptyState, ErrorState, LoadingRows } from "./states"

type MiningState = "idle" | "mining" | "submitting" | "found"

const EMPTY_HASH = "0".repeat(64)

export function MiningView() {
  const { refreshUser } = useAuth()
  const { mutate: globalMutate } = useSWRConfig()

  const [state, setState] = useState<MiningState>("idle")
  const [nonce, setNonce] = useState(0)
  const [hash, setHash] = useState(EMPTY_HASH)
  const [hashRate, setHashRate] = useState(0)
  const [attempts, setAttempts] = useState(0)
  const [error, setError] = useState<string | null>(null)
  const [minedBlock, setMinedBlock] = useState<(Block & { reward: number }) | null>(null)

  const abortRef = useRef<{ aborted: boolean }>({ aborted: false })

  const mempool = useSWR("pending-all", () => api.pending(), { revalidateOnFocus: false })
  const nextBlock = useSWR("next-block", () => api.nextBlock(), { revalidateOnFocus: false })

  const pendingTxs = mempool.data?.transactions ?? []

  useEffect(() => {
    return () => {
      abortRef.current.aborted = true
    }
  }, [])

  const start = useCallback(async () => {
    setError(null)
    setMinedBlock(null)
    setState("mining")
    setAttempts(0)
    setHashRate(0)

    abortRef.current = { aborted: false }
    const signal = abortRef.current

    try {
      // 1. Busca os dados do próximo bloco na API.
      const next = await api.nextBlock()
      if (signal.aborted) return

      if (next.mined) {
        setError("Este bloco já foi minerado por outro nó. Atualizando dados…")
        setState("idle")
        await Promise.all([nextBlock.mutate(), globalMutate("blocks")])
        return
      }

      // 2. Resolve a prova de trabalho no navegador: SHA-256(`${nonce}${data}`).
      const result = await mineBlock(next.data, {
        signal,
        onProgress: (p) => {
          setNonce(p.nonce)
          setHash(p.hash)
          setAttempts(p.attempts)
          if (p.hashRate) setHashRate(p.hashRate)
        },
      })

      if (!result || signal.aborted) return

      setNonce(result.nonce)
      setHash(result.hash)
      setAttempts(result.attempts)

      // 3. Envia o nonce encontrado para a API validar.
      setState("submitting")
      const block = await api.mine(result.nonce)
      if (signal.aborted) return

      setMinedBlock(block)
      setState("found")

      // 4. Atualiza saldo, mempool e cadeia.
      await Promise.all([
        refreshUser(),
        mempool.mutate(),
        nextBlock.mutate(),
        globalMutate("blocks"),
        globalMutate("ranking"),
      ])
    } catch (err) {
      if (signal.aborted) return
      setError(err instanceof Error ? err.message : "Falha ao minerar o bloco.")
      setState("idle")
    }
  }, [globalMutate, mempool, nextBlock, refreshUser])

  const stop = useCallback(() => {
    abortRef.current.aborted = true
    setState("idle")
  }, [])

  const reset = useCallback(() => {
    abortRef.current.aborted = true
    setState("idle")
    setHash(EMPTY_HASH)
    setNonce(0)
    setAttempts(0)
    setHashRate(0)
    setMinedBlock(null)
    setError(null)
  }, [])

  const busy = state === "mining" || state === "submitting"

  return (
    <div className="mx-auto max-w-5xl">
      <header className="mb-8 flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="text-sm text-zinc-500">Consenso</p>
          <h1 className="text-2xl font-semibold tracking-tight text-zinc-50">Mineração</h1>
        </div>
        <span className="rounded-full border border-zinc-800 bg-zinc-900 px-3 py-1 font-mono text-xs text-zinc-400">
          dificuldade {DIFFICULTY} · alvo {TARGET_PREFIX}…
        </span>
      </header>

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-5">
        {/* Mempool */}
        <div className="rounded-2xl border border-zinc-800 bg-zinc-900 p-5 lg:col-span-2">
          <div className="mb-1 flex items-center gap-2">
            <Layers className="size-4 text-yellow-400" aria-hidden="true" />
            <h2 className="text-sm font-semibold text-zinc-100">Mempool</h2>
            <button
              type="button"
              onClick={() => mempool.mutate()}
              aria-label="Atualizar mempool"
              className="ml-auto text-zinc-500 transition-colors hover:text-zinc-300"
            >
              <RefreshCw className={`size-3.5 ${mempool.isValidating ? "animate-spin" : ""}`} />
            </button>
          </div>
          <p className="mb-4 text-xs text-zinc-500">
            <span className="text-zinc-300">{pendingTxs.length}</span> tx aguardando confirmação
          </p>

          {mempool.isLoading ? (
            <LoadingRows rows={3} />
          ) : mempool.error ? (
            <ErrorState
              message={mempool.error instanceof Error ? mempool.error.message : "Erro ao carregar."}
              onRetry={() => mempool.mutate()}
            />
          ) : pendingTxs.length === 0 ? (
            <EmptyState message="Mempool vazio. Blocos podem ser minerados apenas com a recompensa." />
          ) : (
            <ul className="flex flex-col gap-2">
              {pendingTxs.map((tx, i) => (
                <li key={i} className="rounded-xl border border-zinc-800 bg-zinc-950/40 p-3">
                  <div className="flex items-center justify-between gap-2">
                    <span className="truncate font-mono text-xs text-zinc-400">
                      {truncate(tx.from as string)}
                    </span>
                    <span className="shrink-0 text-sm font-semibold text-zinc-100">
                      {formatQtc(typeof tx.amount === "number" ? tx.amount : 0)}
                    </span>
                  </div>
                  <div className="mt-1 flex items-center justify-between gap-2">
                    <span className="truncate font-mono text-[11px] text-zinc-600">
                      → {truncate(tx.to as string)}
                    </span>
                    <span className="shrink-0 text-[11px] text-yellow-400">
                      {relativeTime(tx.created_at)}
                    </span>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>

        {/* Console */}
        <div className="lg:col-span-3">
          <div
            className={`flex flex-col rounded-2xl border bg-zinc-900 p-6 transition-all duration-300 ${
              state === "found"
                ? "border-emerald-400/60 shadow-[0_0_40px_-10px] shadow-emerald-400/40"
                : busy
                  ? "border-yellow-400/50 shadow-[0_0_40px_-12px] shadow-yellow-400/40"
                  : "border-zinc-800"
            }`}
          >
            <div className="flex items-center gap-2 text-sm text-zinc-400">
              <Cpu className="size-4 text-yellow-400" aria-hidden="true" />
              Console de mineração
            </div>

            {/* Botão principal */}
            <div className="my-6 flex justify-center">
              {state !== "found" ? (
                <button
                  type="button"
                  onClick={busy ? stop : start}
                  disabled={state === "submitting"}
                  className={`flex size-40 flex-col items-center justify-center gap-2 rounded-full text-sm font-semibold transition-all disabled:cursor-not-allowed ${
                    busy
                      ? "animate-pulse bg-zinc-800 text-zinc-200 ring-4 ring-yellow-400/30"
                      : "bg-yellow-400 text-zinc-950 shadow-[0_0_40px_-6px] shadow-yellow-400/70 hover:scale-105 hover:bg-yellow-300"
                  }`}
                >
                  {state === "submitting" ? (
                    <>
                      <RefreshCw className="size-7 animate-spin" aria-hidden="true" />
                      Validando…
                    </>
                  ) : state === "mining" ? (
                    <>
                      <Square className="size-7" aria-hidden="true" />
                      Parar
                    </>
                  ) : (
                    <>
                      <Pickaxe className="size-8" aria-hidden="true" />
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
                  <Sparkles className="size-8" aria-hidden="true" />
                  Minerar novamente
                </button>
              )}
            </div>

            {/* Leitura ao vivo */}
            <div className="rounded-xl border border-zinc-800 bg-black/50 p-4 font-mono text-xs">
              <div className="flex items-center justify-between border-b border-zinc-800 pb-2">
                <span className="text-zinc-600">nonce</span>
                <span className={busy ? "text-yellow-400" : "text-zinc-300"}>
                  {nonce.toLocaleString("pt-BR")}
                </span>
              </div>
              <div className="flex items-center justify-between gap-2 border-b border-zinc-800 py-2">
                <span className="shrink-0 text-zinc-600">hash</span>
                <span
                  className={`truncate ${
                    state === "found" ? "text-emerald-400" : busy ? "text-yellow-400" : "text-zinc-500"
                  }`}
                >
                  {hash}
                </span>
              </div>
              <div className="flex items-center justify-between pt-2 text-zinc-600">
                <span>tentativas {attempts.toLocaleString("pt-BR")}</span>
                <span>{hashRate ? `${hashRate.toLocaleString("pt-BR")} H/s` : "—"}</span>
              </div>
            </div>

            {/* Mensagens de estado */}
            {error && (
              <p
                role="alert"
                className="mt-4 rounded-xl border border-red-900/50 bg-red-950/30 px-3 py-2 text-center text-xs leading-relaxed text-red-200"
              >
                {error}
              </p>
            )}
            {state === "found" && minedBlock && (
              <div className="mt-4 flex flex-col items-center gap-1 rounded-xl border border-emerald-400/40 bg-emerald-500/10 py-3 text-sm font-semibold text-emerald-400">
                <span className="flex items-center gap-2">
                  <Sparkles className="size-4" aria-hidden="true" />
                  Bloco #{minedBlock.index} minerado! +{formatQtc(minedBlock.reward)} QTC
                </span>
                <span className="font-mono text-[10px] font-normal text-emerald-300/70">
                  nonce {minedBlock.nonce} · {truncate(minedBlock.hash, 16, 8)}
                </span>
              </div>
            )}
            {state === "idle" && !error && (
              <p className="mt-4 text-center text-xs text-zinc-600">
                Pressione iniciar para buscar um hash com {DIFFICULTY} zeros à esquerda.
              </p>
            )}
            {state === "mining" && (
              <p className="mt-4 text-center text-xs text-zinc-500">
                Calculando SHA-256(nonce + data) no navegador…
              </p>
            )}
            {state === "submitting" && (
              <p className="mt-4 text-center text-xs text-zinc-500">
                Enviando o nonce para validação na API…
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
