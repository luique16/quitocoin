"use client"

import { createContext, useCallback, useContext, useEffect, useState } from "react"
import useSWR from "swr"
import { api, ApiError, clearToken, getToken, setToken, type User } from "@/lib/api"

type AuthContextValue = {
  user: User | null
  loading: boolean
  error: string | null
  login: (email: string, password: string) => Promise<void>
  register: (name: string, email: string, password: string) => Promise<void>
  logout: () => void
  refreshUser: () => Promise<User | undefined>
  setUser: (u: User) => void
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [hasToken, setHasToken] = useState<boolean | null>(null)

  useEffect(() => {
    setHasToken(!!getToken())
  }, [])

  const {
    data: user,
    error,
    isLoading,
    mutate,
  } = useSWR<User>(hasToken ? "/me" : null, () => api.me(), {
    revalidateOnFocus: false,
    shouldRetryOnError: false,
    onError: (err) => {
      if (err instanceof ApiError && err.status === 401) {
        clearToken()
        setHasToken(false)
      }
    },
  })

  const login = useCallback(
    async (email: string, password: string) => {
      const { token } = await api.login(email, password)
      setToken(token)
      setHasToken(true)
      await mutate()
    },
    [mutate],
  )

  const register = useCallback(
    async (name: string, email: string, password: string) => {
      const { token } = await api.register(name, email, password)
      setToken(token)
      setHasToken(true)
      await mutate()
    },
    [mutate],
  )

  const logout = useCallback(() => {
    clearToken()
    setHasToken(false)
    mutate(undefined, { revalidate: false })
  }, [mutate])

  const refreshUser = useCallback(() => mutate(), [mutate])
  const setUserLocal = useCallback((u: User) => mutate(u, { revalidate: false }), [mutate])

  return (
    <AuthContext.Provider
      value={{
        user: user ?? null,
        loading: hasToken === null || (!!hasToken && isLoading),
        error: error instanceof Error ? error.message : null,
        login,
        register,
        logout,
        refreshUser,
        setUser: setUserLocal,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error("useAuth must be used within AuthProvider")
  return ctx
}
