import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { TOKEN_KEY } from "@/lib/api"
import { DashboardShell } from "@/components/quito/dashboard-shell"

export default async function DashboardPage() {
  const store = await cookies()
  if (!store.get(TOKEN_KEY)) redirect("/")

  return <DashboardShell />
}
