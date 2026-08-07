"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import type { View } from "@/lib/quito-data"
import { DashboardView } from "./dashboard-view"
import { TransferView } from "./transfer-view"
import { ExplorerView } from "./explorer-view"
import { MiningView } from "./mining-view"
import { AccountView } from "./account-view"
import { MobileNav, Sidebar } from "./sidebar"
import { useAuth } from "./auth-provider"

export function DashboardShell() {
  const { logout } = useAuth()
  const router = useRouter()
  const [view, setView] = useState<View>("dashboard")

  const handleLogout = () => {
    logout()
    router.replace("/")
  }

  return (
    <div className="min-h-screen bg-zinc-950">
      <Sidebar current={view} onNavigate={setView} onLogout={handleLogout} />
      <MobileNav current={view} onNavigate={setView} />
      <main className="px-4 pb-24 pt-6 md:ml-64 md:px-8 md:pb-8 md:pt-10">
        {view === "dashboard" && <DashboardView />}
        {view === "transfer" && <TransferView />}
        {view === "explorer" && <ExplorerView />}
        {view === "mining" && <MiningView />}
        {view === "account" && <AccountView onLogout={handleLogout} />}
      </main>
    </div>
  )
}
