"use client"

import { useState } from "react"
import { ArrowRight, Lock, Mail, User } from "lucide-react"
import { QuitoWordmark } from "./logo"
import { useAuth } from "./auth-provider"
import { Spinner } from "./states"
import { API_URL } from "@/lib/api"

function NetworkIllustration() {
  const nodes = [
    { x: 80, y: 90 },
    { x: 210, y: 60 },
    { x: 330, y: 130 },
    { x: 140, y: 220 },
    { x: 290, y: 260 },
    { x: 60, y: 330 },
    { x: 360, y: 350 },
    { x: 200, y: 400 },
    { x: 110, y: 470 },
    { x: 320, y: 470 },
  ]
  const links: [number, number][] = [
    [0, 1],
    [1, 2],
    [0, 3],
    [1, 4],
    [3, 4],
    [3, 5],
    [4, 6],
    [4, 7],
    [5, 8],
    [7, 8],
    [7, 9],
    [6, 9],
    [2, 4],
  ]
  return (
    <svg viewBox="0 0 420 540" className="h-full w-full" aria-hidden="true">
      <defs>
        <radialGradient id="glow" cx="50%" cy="40%" r="70%">
          <stop offset="0%" stopColor="rgb(250 204 21 / 0.12)" />
          <stop offset="100%" stopColor="rgb(250 204 21 / 0)" />
        </radialGradient>
      </defs>
      <rect width="420" height="540" fill="url(#glow)" />
      {links.map(([a, b], i) => (
        <line
          key={i}
          x1={nodes[a].x}
          y1={nodes[a].y}
          x2={nodes[b].x}
          y2={nodes[b].y}
          stroke="rgb(250 204 21 / 0.25)"
          strokeWidth="1"
        />
      ))}
      {nodes.map((n, i) => (
        <g key={i}>
          <circle cx={n.x} cy={n.y} r={i % 3 === 0 ? 6 : 4} fill="rgb(250 204 21)">
            <animate
              attributeName="opacity"
              values="0.4;1;0.4"
              dur={`${2 + (i % 4)}s`}
              repeatCount="indefinite"
            />
          </circle>
          <circle
            cx={n.x}
            cy={n.y}
            r={i % 3 === 0 ? 12 : 9}
            fill="none"
            stroke="rgb(250 204 21 / 0.2)"
            strokeWidth="1"
          />
        </g>
      ))}
    </svg>
  )
}

function Field({
  icon: Icon,
  label,
  type,
  placeholder,
  value,
  onChange,
  autoComplete,
  required = true,
  minLength,
}: {
  icon: typeof Mail
  label: string
  type: string
  placeholder: string
  value: string
  onChange: (v: string) => void
  autoComplete?: string
  required?: boolean
  minLength?: number
}) {
  return (
    <label className="block">
      <span className="mb-1.5 block text-xs font-medium text-zinc-400">{label}</span>
      <div className="relative">
        <Icon className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-zinc-500" />
        <input
          type={type}
          placeholder={placeholder}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          autoComplete={autoComplete}
          required={required}
          minLength={minLength}
          className="w-full rounded-xl border border-zinc-700 bg-zinc-950/60 py-2.5 pl-10 pr-3 text-sm text-zinc-100 placeholder:text-zinc-600 outline-none transition-shadow focus:border-yellow-400 focus:ring-2 focus:ring-yellow-400/30"
        />
      </div>
    </label>
  )
}

export function AuthView({ onEnter }: { onEnter: () => void }) {
  const { login, register } = useAuth()
  const [mode, setMode] = useState<"login" | "register">("login")
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    setSubmitting(true)
    try {
      if (mode === "login") {
        await login(email, password)
      } else {
        await register(name, email, password)
      }
      onEnter()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Authentication failed.")
    } finally {
      setSubmitting(false)
    }
  }

  function switchMode(m: "login" | "register") {
    setMode(m)
    setError(null)
  }

  return (
    <div className="grid min-h-screen grid-cols-1 lg:grid-cols-2">
      {/* Left: illustration */}
      <div className="relative hidden overflow-hidden border-r border-zinc-800 bg-zinc-900 lg:block">
        <div className="absolute inset-0 opacity-90">
          <NetworkIllustration />
        </div>
        <div className="absolute bottom-0 left-0 right-0 p-12">
          <h2 className="max-w-sm text-balance text-3xl font-semibold leading-tight text-zinc-50">
            An educational, clean and transparent blockchain.
          </h2>
          <p className="mt-3 max-w-sm text-pretty text-sm leading-relaxed text-zinc-400">
            Learn how transfers, blocks and mining work in practice with the QuitoCoin
            network.
          </p>
        </div>
      </div>

      {/* Right: form */}
      <div className="flex items-center justify-center bg-zinc-950 p-6">
        <div className="w-full max-w-sm">
          <div className="mb-8 flex justify-center lg:justify-start">
            <QuitoWordmark />
          </div>

          <div className="rounded-2xl border border-zinc-800 bg-zinc-900/60 p-6 shadow-2xl backdrop-blur-xl">
            <div className="mb-6 flex rounded-xl bg-zinc-950/60 p-1">
              {(["login", "register"] as const).map((m) => (
                <button
                  key={m}
                  type="button"
                  onClick={() => switchMode(m)}
                  className={`flex-1 rounded-lg py-2 text-sm font-medium transition-colors ${
                    mode === m
                      ? "bg-zinc-800 text-zinc-50"
                      : "text-zinc-500 hover:text-zinc-300"
                  }`}
                >
                  {m === "login" ? "Sign in" : "Register"}
                </button>
              ))}
            </div>

            <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
              {mode === "register" && (
                <Field
                  icon={User}
                  label="Name"
                  type="text"
                  placeholder="Your name"
                  value={name}
                  onChange={setName}
                  autoComplete="name"
                />
              )}
              <Field
                icon={Mail}
                label="Email"
                type="email"
                placeholder="you@example.com"
                value={email}
                onChange={setEmail}
                autoComplete="email"
              />
              <Field
                icon={Lock}
                label="Password"
                type="password"
                placeholder="••••••••"
                value={password}
                onChange={setPassword}
                autoComplete={mode === "login" ? "current-password" : "new-password"}
                minLength={mode === "register" ? 6 : undefined}
              />

              {error && (
                <p role="alert" className="rounded-lg border border-red-900/50 bg-red-950/30 px-3 py-2 text-xs leading-relaxed text-red-200">
                  {error}
                </p>
              )}

              <button
                type="submit"
                disabled={submitting}
                className="mt-2 flex items-center justify-center gap-2 rounded-xl bg-yellow-400 py-3 text-sm font-semibold text-zinc-950 shadow-[0_0_24px_-6px] shadow-yellow-400/60 transition-all hover:bg-yellow-300 hover:shadow-yellow-400/80 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {submitting ? (
                  <>
                    <Spinner />
                    {mode === "login" ? "Signing in…" : "Creating account…"}
                  </>
                ) : (
                  <>
                    {mode === "login" ? "Sign in to the network" : "Create account"}
                    <ArrowRight className="size-4" />
                  </>
                )}
              </button>
            </form>
          </div>

          <p className="mt-6 text-center text-xs text-zinc-600">
            Educational project — no real value is moved.
          </p>
          <p className="mt-1 text-center font-mono text-[10px] text-zinc-700">API: {API_URL}</p>
        </div>
      </div>
    </div>
  )
}
