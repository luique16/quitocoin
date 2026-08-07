"use client"

import useSWR from "swr"
import { Box, Crown, RefreshCw, Radio } from "lucide-react"
import { api } from "@/lib/api"
import { formatQtc, relativeTime, truncate } from "@/lib/quito-data"
import { useAuth } from "./auth-provider"
import { ComingSoon, EmptyState, ErrorState, LoadingRows } from "./states"

export function ExplorerView() {
  const blocks = useSWR("blocks", () => api.blocks(12, 0), { revalidateOnFocus: false })
  const list = blocks.data?.blocks ?? []
  const total = blocks.data?.total_count ?? 0

  return (
    <div className="mx-auto max-w-6xl">
      <header className="mb-8 flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="text-sm text-zinc-500">Network</p>
          <h1 className="text-2xl font-semibold tracking-tight text-zinc-50">
            Blockchain Explorer
          </h1>
        </div>
        <div className="flex items-center gap-3">
          <span className="rounded-full border border-zinc-800 bg-zinc-900 px-3 py-1 text-xs text-zinc-400">
            {total} {total === 1 ? "block" : "blocks"} on the chain
          </span>
          <button
            type="button"
            onClick={() => blocks.mutate()}
            className="inline-flex items-center gap-1.5 rounded-full border border-zinc-800 bg-zinc-900 px-3 py-1 text-xs text-zinc-300 transition-colors hover:border-zinc-700 hover:text-zinc-100"
          >
            <RefreshCw
              className={`size-3 ${blocks.isValidating ? "animate-spin" : ""}`}
              aria-hidden="true"
            />
            Refresh
          </button>
        </div>
      </header>

      {/* Timeline de blocos */}
      <section className="mb-6">
        <div className="mb-3 flex items-center gap-2">
          <Box className="size-4 text-yellow-400" aria-hidden="true" />
          <h2 className="text-sm font-semibold text-zinc-100">Recent blocks</h2>
        </div>

        {blocks.isLoading ? (
          <LoadingRows rows={3} />
        ) : blocks.error ? (
          <ErrorState
            message={blocks.error instanceof Error ? blocks.error.message : "Error loading blocks."}
            onRetry={() => blocks.mutate()}
          />
        ) : list.length === 0 ? (
          <EmptyState message="No blocks on the chain yet. Mine the first block in the Mining tab." />
        ) : (
          <div className="flex flex-nowrap gap-3 overflow-x-auto pb-3">
            {list.map((b) => (
              <div
                key={b.hash || b.index}
                className="w-52 shrink-0 rounded-2xl border border-zinc-800 bg-zinc-900 p-4 transition-colors hover:border-yellow-400/50"
              >
                <div className="flex items-center justify-between">
                  <span className="text-lg font-bold text-yellow-400">#{b.index}</span>
                  <span className="rounded-full bg-zinc-800 px-2 py-0.5 text-[10px] text-zinc-400">
                    {b.tx_count ?? b.transactions?.length ?? 0} tx
                  </span>
                </div>
                <p className="mt-3 text-[10px] uppercase tracking-widest text-zinc-600">Hash</p>
                <p className="truncate font-mono text-xs text-zinc-300">{truncate(b.hash, 12, 6)}</p>
                <p className="mt-2 text-[10px] uppercase tracking-widest text-zinc-600">Miner</p>
                <p className="truncate font-mono text-xs text-zinc-400">{truncate(b.miner, 8, 4)}</p>
                <p className="mt-3 font-mono text-[11px] text-zinc-600">{relativeTime(b.created_at)}</p>
              </div>
            ))}
          </div>
        )}
      </section>

      {/* Split: rich list + logs */}
      <div className="grid grid-cols-1 items-stretch gap-4 lg:grid-cols-2">
        <RichList />
        <NetworkLogsSoon />
      </div>
    </div>
  )
}

function RichList() {
  const { user } = useAuth()
  const ranking = useSWR("ranking", () => api.ranking(), { revalidateOnFocus: false })
  const entries = ranking.data?.richest ?? []

  return (
    <div className="rounded-2xl border border-zinc-800 bg-zinc-900 p-5">
      <div className="mb-4 flex items-center gap-2">
        <Crown className="size-4 text-yellow-400" aria-hidden="true" />
        <h2 className="text-sm font-semibold text-zinc-100">Rich List</h2>
        <span className="ml-auto text-xs text-zinc-500">Highest balances</span>
      </div>

      {ranking.isLoading ? (
        <LoadingRows rows={5} />
      ) : ranking.error ? (
        <ErrorState
          message={ranking.error instanceof Error ? ranking.error.message : "Error loading the ranking."}
          onRetry={() => ranking.mutate()}
        />
      ) : entries.length === 0 ? (
        <EmptyState message="The ranking is currently empty." />
      ) : (
        <ul className="flex flex-col gap-1">
          {entries.map((entry, i) => {
            const isUser = entry.public_id === user?.public_id
            return (
              <li
                key={entry.public_id}
                className={`flex items-center gap-3 rounded-lg px-2 py-2 ${
                  isUser ? "bg-yellow-400/10 ring-1 ring-yellow-400/30" : ""
                }`}
              >
                <span
                  className={`flex size-6 shrink-0 items-center justify-center rounded-md text-xs font-bold ${
                    i < 3 ? "bg-yellow-400 text-zinc-950" : "bg-zinc-800 text-zinc-400"
                  }`}
                >
                  {i + 1}
                </span>
                <span className="min-w-0 flex-1 truncate font-mono text-xs text-zinc-300">
                  {entry.public_id}
                  {isUser && <span className="ml-2 text-[10px] text-yellow-400">you</span>}
                </span>
                <span className="shrink-0 font-mono text-xs font-semibold text-zinc-100">
                  {formatQtc(entry.balance)}
                </span>
              </li>
            )
          })}
        </ul>
      )}
    </div>
  )
}

function NetworkLogsSoon() {
  return (
    <div className="flex h-full min-h-[420px] flex-col overflow-hidden rounded-2xl border border-zinc-800 bg-black/50">
      <div className="flex shrink-0 items-center gap-2 border-b border-zinc-800 px-4 py-3">
        <span className="inline-flex size-2.5 rounded-full bg-zinc-600" />
        <span className="text-xs font-semibold text-zinc-300">Network logs</span>
        <ComingSoon />
        <span className="ml-auto font-mono text-[10px] uppercase tracking-widest text-zinc-700">
          websocket · offline
        </span>
      </div>
      <div className="flex min-h-0 flex-1 flex-col items-center justify-center gap-3 p-8">
        <span className="flex size-12 items-center justify-center rounded-full border border-zinc-800 bg-zinc-900">
          <Radio className="size-5 text-zinc-600" aria-hidden="true" />
        </span>
        <p className="max-w-xs text-pretty text-center text-sm leading-relaxed text-zinc-500">
          Real-time event streaming via WebSocket will be enabled once the endpoint is available
          in the API.
        </p>
      </div>
    </div>
  )
}
