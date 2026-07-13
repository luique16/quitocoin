"use client"

import { useState } from "react"
import type { View } from "@/lib/quito-data"
import { AuthView } from "@/components/quito/auth-view"
import { DashboardView } from "@/components/quito/dashboard-view"
import { TransferView } from "@/components/quito/transfer-view"
import { ExplorerView } from "@/components/quito/explorer-view"
import { MiningView } from "@/components/quito/mining-view"
import { AccountView } from "@/components/quito/account-view"
import { MobileNav, Sidebar } from "@/components/quito/sidebar"

export default function Page() {
  const [view, setView] = useState<View>("auth")

  if (view === "auth") {
    return <AuthView onEnter={() => setView("dashboard")} />
  }

  return (
    <div className="min-h-screen bg-zinc-950">
      <Sidebar current={view} onNavigate={setView} onLogout={() => setView("auth")} />
      <MobileNav current={view} onNavigate={setView} />
      <main className="px-4 pb-24 pt-6 md:ml-64 md:px-8 md:pb-8 md:pt-10">
        {view === "dashboard" && <DashboardView />}
        {view === "transfer" && <TransferView />}
        {view === "explorer" && <ExplorerView />}
        {view === "mining" && <MiningView />}
        {view === "account" && <AccountView onLogout={() => setView("auth")} />}
      </main>
    </div>
  )
}
