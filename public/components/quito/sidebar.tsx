"use client"

import { LayoutDashboard, Send, Boxes, Pickaxe, LogOut } from "lucide-react"
import type { View } from "@/lib/quito-data"
import { QuitoWordmark } from "./logo"

const NAV: { id: View; label: string; icon: typeof LayoutDashboard }[] = [
  { id: "dashboard", label: "Dashboard", icon: LayoutDashboard },
  { id: "transfer", label: "Transferir", icon: Send },
  { id: "explorer", label: "Explorer", icon: Boxes },
  { id: "mining", label: "Mineração", icon: Pickaxe },
]

export function Sidebar({
  current,
  onNavigate,
  onLogout,
}: {
  current: View
  onNavigate: (v: View) => void
  onLogout: () => void
}) {
  return (
    <aside className="fixed inset-y-0 left-0 z-30 hidden w-64 flex-col border-r border-zinc-800 bg-zinc-900 p-4 md:flex">
      <div className="px-2 py-3">
        <QuitoWordmark />
      </div>

      <nav className="mt-6 flex flex-1 flex-col gap-1">
        {NAV.map(({ id, label, icon: Icon }) => {
          const active = current === id
          return (
            <button
              key={id}
              type="button"
              onClick={() => onNavigate(id)}
              className={`group flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-all ${
                active
                  ? "bg-yellow-400 text-zinc-950 shadow-[0_0_20px_-4px] shadow-yellow-400/50"
                  : "text-zinc-400 hover:bg-zinc-800 hover:text-zinc-100"
              }`}
            >
              <Icon className="size-[18px]" strokeWidth={active ? 2.4 : 2} />
              {label}
            </button>
          )
        })}
      </nav>

      <button
        type="button"
        onClick={onLogout}
        className="flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium text-zinc-400 transition-colors hover:bg-zinc-800 hover:text-red-400"
      >
        <LogOut className="size-[18px]" />
        Sair da rede
      </button>
    </aside>
  )
}

export function MobileNav({
  current,
  onNavigate,
}: {
  current: View
  onNavigate: (v: View) => void
}) {
  return (
    <nav className="fixed inset-x-0 bottom-0 z-30 flex items-center justify-around border-t border-zinc-800 bg-zinc-900/95 px-2 py-2 backdrop-blur md:hidden">
      {NAV.map(({ id, label, icon: Icon }) => {
        const active = current === id
        return (
          <button
            key={id}
            type="button"
            onClick={() => onNavigate(id)}
            className={`flex flex-1 flex-col items-center gap-1 rounded-lg py-1.5 text-[10px] font-medium transition-colors ${
              active ? "text-yellow-400" : "text-zinc-500"
            }`}
          >
            <Icon className="size-5" strokeWidth={active ? 2.4 : 2} />
            {label}
          </button>
        )
      })}
    </nav>
  )
}
