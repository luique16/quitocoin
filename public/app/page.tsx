"use client"

import { useRouter } from "next/navigation"
import { AuthView } from "@/components/quito/auth-view"

export default function Page() {
  const router = useRouter()

  return <AuthView onEnter={() => router.replace("/dashboard")} />
}
