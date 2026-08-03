"use client"

import { AlertTriangle, Loader2, RefreshCw } from "lucide-react"

export function Spinner({ className = "size-4" }: { className?: string }) {
  return <Loader2 className={`animate-spin ${className}`} aria-hidden="true" />
}

export function LoadingRows({ rows = 4, className = "" }: { rows?: number; className?: string }) {
  return (
    <div className={`flex flex-col gap-2 ${className}`} aria-busy="true" aria-live="polite">
      <span className="sr-only">Carregando dados</span>
      {Array.from({ length: rows }).map((_, i) => (
        <div key={i} className="h-12 animate-pulse rounded-lg bg-zinc-800/60" />
      ))}
    </div>
  )
}

export function ErrorState({ message, onRetry }: { message: string; onRetry?: () => void }) {
  return (
    <div className="flex flex-col items-start gap-3 rounded-lg border border-red-900/50 bg-red-950/30 p-4">
      <div className="flex gap-2.5">
        <AlertTriangle className="mt-0.5 size-4 shrink-0 text-red-400" aria-hidden="true" />
        <p className="text-sm leading-relaxed text-red-200">{message}</p>
      </div>
      {onRetry && (
        <button
          type="button"
          onClick={onRetry}
          className="ml-6 inline-flex items-center gap-1.5 rounded-md border border-red-900/60 px-2.5 py-1 text-xs font-medium text-red-200 transition-colors hover:bg-red-950/60"
        >
          <RefreshCw className="size-3" aria-hidden="true" />
          Tentar novamente
        </button>
      )}
    </div>
  )
}

export function EmptyState({ message }: { message: string }) {
  return (
    <div className="flex items-center justify-center rounded-lg border border-dashed border-zinc-800 px-4 py-8">
      <p className="text-center text-sm text-zinc-500">{message}</p>
    </div>
  )
}

export function ComingSoon({ label = "Em breve" }: { label?: string }) {
  return (
    <span className="rounded-full border border-yellow-500/30 bg-yellow-500/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider text-yellow-500">
      {label}
    </span>
  )
}
