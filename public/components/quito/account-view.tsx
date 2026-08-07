"use client"

import { useEffect, useState } from "react"
import { Check, KeyRound, Loader2, Mail, Trash2, User, X } from "lucide-react"
import { api } from "@/lib/api"
import { formatDateTime } from "@/lib/quito-data"
import { useAuth } from "./auth-provider"

export function AccountView({ onLogout }: { onLogout?: () => void }) {
  const { user, setUser, logout } = useAuth()

  const handleLogout = () => {
    logout()
    onLogout?.()
  }
  const [name, setName] = useState(user?.name ?? "")
  const [email, setEmail] = useState(user?.email ?? "")
  const [savedFlash, setSavedFlash] = useState(false)
  const [saving, setSaving] = useState(false)
  const [profileError, setProfileError] = useState<string | null>(null)
  const [passwordOpen, setPasswordOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)

  // Syncs the fields when the user is loaded from the API.
  useEffect(() => {
    if (user) {
      setName(user.name)
      setEmail(user.email)
    }
  }, [user])

  async function handleProfileSave(e: React.FormEvent) {
    e.preventDefault()
    setProfileError(null)
    setSaving(true)
    try {
      const updated = await api.updateMe({ name, email })
      setUser(updated)
      setSavedFlash(true)
      setTimeout(() => setSavedFlash(false), 2200)
    } catch (err) {
      setProfileError(err instanceof Error ? err.message : "Failed to update the profile.")
    } finally {
      setSaving(false)
    }
  }

  const dirty = user ? name !== user.name || email !== user.email : false

  return (
    <div className="mx-auto max-w-3xl">
      <header className="mb-8">
        <p className="text-sm text-zinc-500">Settings</p>
        <h1 className="text-2xl font-semibold tracking-tight text-zinc-50">Your account</h1>
      </header>

      {/* Profile */}
      <section className="mb-6 rounded-2xl border border-zinc-800 bg-zinc-900 p-6">
        <h2 className="mb-1 text-sm font-semibold text-zinc-100">Profile</h2>
        <p className="mb-5 text-xs text-zinc-500">Update your public identifying information.</p>

        <form onSubmit={handleProfileSave} className="flex flex-col gap-4">
          <label className="flex flex-col gap-1.5">
            <span className="text-xs font-medium text-zinc-400">Name</span>
            <div className="relative">
              <User className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-zinc-500" />
              <input
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full rounded-xl border border-zinc-700 bg-zinc-950/60 py-2.5 pl-10 pr-3 text-sm text-zinc-100 outline-none transition-shadow focus:border-yellow-400 focus:ring-2 focus:ring-yellow-400/30"
              />
            </div>
          </label>

          <label className="flex flex-col gap-1.5">
            <span className="text-xs font-medium text-zinc-400">Email</span>
            <div className="relative">
              <Mail className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-zinc-500" />
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full rounded-xl border border-zinc-700 bg-zinc-950/60 py-2.5 pl-10 pr-3 text-sm text-zinc-100 outline-none transition-shadow focus:border-yellow-400 focus:ring-2 focus:ring-yellow-400/30"
              />
            </div>
          </label>

          {profileError && (
            <p
              role="alert"
              className="rounded-lg border border-red-900/50 bg-red-950/30 px-3 py-2 text-xs leading-relaxed text-red-200"
            >
              {profileError}
            </p>
          )}

          <div className="flex items-center gap-3">
            <button
              type="submit"
              disabled={saving || !dirty}
              className="flex items-center gap-2 rounded-xl bg-yellow-400 px-5 py-2.5 text-sm font-semibold text-zinc-950 transition-colors hover:bg-yellow-300 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {saving && <Loader2 className="size-4 animate-spin" aria-hidden="true" />}
              {saving ? "Saving…" : "Save changes"}
            </button>
            {savedFlash && (
              <span className="flex items-center gap-1.5 text-xs font-medium text-emerald-400">
                <Check className="size-4" aria-hidden="true" /> Profile updated
              </span>
            )}
          </div>
        </form>

        <dl className="mt-6 grid grid-cols-1 gap-3 border-t border-zinc-800 pt-5 sm:grid-cols-2">
          <div>
            <dt className="text-[10px] uppercase tracking-widest text-zinc-600">Public code</dt>
            <dd className="mt-1 truncate font-mono text-xs text-zinc-400">{user?.public_id ?? "—"}</dd>
          </div>
          <div>
            <dt className="text-[10px] uppercase tracking-widest text-zinc-600">Account created on</dt>
            <dd className="mt-1 font-mono text-xs text-zinc-400">{formatDateTime(user?.created_at)}</dd>
          </div>
        </dl>
      </section>

      {/* Security */}
      <section className="mb-6 rounded-2xl border border-zinc-800 bg-zinc-900 p-6">
        <h2 className="mb-1 text-sm font-semibold text-zinc-100">Security</h2>
        <p className="mb-5 text-xs text-zinc-500">Manage the credentials to access your wallet.</p>

        <div className="flex items-center justify-between rounded-xl border border-zinc-800 bg-zinc-950/40 px-4 py-3.5">
          <div className="flex items-center gap-3">
            <span className="flex size-9 items-center justify-center rounded-lg bg-zinc-800 text-yellow-400">
              <KeyRound className="size-4" />
            </span>
            <div>
              <p className="text-sm font-medium text-zinc-100">Password</p>
              <p className="text-xs text-zinc-500">Use a strong, unique password</p>
            </div>
          </div>
          <button
            type="button"
            onClick={() => setPasswordOpen(true)}
            className="rounded-xl border border-zinc-700 px-4 py-2 text-sm font-medium text-zinc-100 transition-colors hover:border-yellow-400 hover:text-yellow-400"
          >
            Change password
          </button>
        </div>
      </section>

      {/* Danger zone */}
      <section className="rounded-2xl border border-red-500/30 bg-red-500/5 p-6">
        <h2 className="mb-1 text-sm font-semibold text-red-400">Danger zone</h2>
        <p className="mb-5 text-xs text-zinc-400">
          Account deletion is permanent and removes your wallet from the network. This action cannot be undone.
        </p>
        <button
          type="button"
          onClick={() => setDeleteOpen(true)}
          className="flex items-center gap-2 rounded-xl bg-red-500 px-5 py-2.5 text-sm font-semibold text-zinc-50 transition-colors hover:bg-red-600"
        >
          <Trash2 className="size-4" />
          Delete account
        </button>
      </section>

      {passwordOpen && <ChangePasswordModal onClose={() => setPasswordOpen(false)} />}
      {deleteOpen && <DeleteAccountModal onClose={() => setDeleteOpen(false)} onDeleted={handleLogout} />}
    </div>
  )
}

function Modal({ children, onClose }: { children: React.ReactNode; onClose: () => void }) {
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose()
    }
    document.addEventListener("keydown", onKey)
    document.body.style.overflow = "hidden"
    return () => {
      document.removeEventListener("keydown", onKey)
      document.body.style.overflow = ""
    }
  }, [onClose])

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4"
      role="dialog"
      aria-modal="true"
    >
      <div className="absolute inset-0 bg-black/70 backdrop-blur-sm" onClick={onClose} />
      <div className="relative z-10 w-full max-w-md rounded-2xl border border-zinc-800 bg-zinc-900 p-6 shadow-2xl">
        {children}
      </div>
    </div>
  )
}

function ChangePasswordModal({ onClose }: { onClose: () => void }) {
  const [current, setCurrent] = useState("")
  const [next, setNext] = useState("")
  const [confirm, setConfirm] = useState("")
  const [status, setStatus] = useState<"idle" | "loading" | "error">("idle")
  const [error, setError] = useState("")

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (next.length < 6) {
      setStatus("error")
      setError("The new password must have at least 6 characters.")
      return
    }
    if (next !== confirm) {
      setStatus("error")
      setError("The confirmation does not match the new password.")
      return
    }
    if (next === current) {
      setStatus("error")
      setError("The new password must be different from the current one.")
      return
    }

    setError("")
    setStatus("loading")
    try {
      await api.updatePassword(current, next)
      onClose()
    } catch (err) {
      setStatus("error")
      setError(err instanceof Error ? err.message : "Could not change the password.")
    }
  }

  return (
    <Modal onClose={onClose}>
      <div className="mb-5 flex items-start justify-between">
        <div>
          <h3 className="text-base font-semibold text-zinc-50">Change password</h3>
          <p className="mt-0.5 text-xs text-zinc-500">Enter your current password and choose a new one.</p>
        </div>
        <button
          type="button"
          onClick={onClose}
          aria-label="Close"
          className="rounded-lg p-1 text-zinc-500 transition-colors hover:bg-zinc-800 hover:text-zinc-100"
        >
          <X className="size-4" />
        </button>
      </div>

      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <label className="flex flex-col gap-1.5">
          <span className="text-xs font-medium text-zinc-400">Current password</span>
          <input
            type="password"
            value={current}
            onChange={(e) => setCurrent(e.target.value)}
            required
            placeholder="••••••••"
            className="w-full rounded-xl border border-zinc-700 bg-zinc-950/60 px-3 py-2.5 text-sm text-zinc-100 outline-none transition-shadow focus:border-yellow-400 focus:ring-2 focus:ring-yellow-400/30"
          />
        </label>
        <label className="flex flex-col gap-1.5">
          <span className="text-xs font-medium text-zinc-400">New password</span>
          <input
            type="password"
            value={next}
            onChange={(e) => setNext(e.target.value)}
            required
            placeholder="••••••••"
            className="w-full rounded-xl border border-zinc-700 bg-zinc-950/60 px-3 py-2.5 text-sm text-zinc-100 outline-none transition-shadow focus:border-yellow-400 focus:ring-2 focus:ring-yellow-400/30"
          />
        </label>
        <label className="flex flex-col gap-1.5">
          <span className="text-xs font-medium text-zinc-400">Confirm new password</span>
          <input
            type="password"
            value={confirm}
            onChange={(e) => setConfirm(e.target.value)}
            required
            placeholder="••••••••"
            className="w-full rounded-xl border border-zinc-700 bg-zinc-950/60 px-3 py-2.5 text-sm text-zinc-100 outline-none transition-shadow focus:border-yellow-400 focus:ring-2 focus:ring-yellow-400/30"
          />
        </label>

        {error && <p className="text-xs font-medium text-red-400">{error}</p>}

        <div className="mt-1 flex items-center justify-end gap-3">
          <button
            type="button"
            onClick={onClose}
            className="rounded-xl px-4 py-2.5 text-sm font-medium text-zinc-400 transition-colors hover:text-zinc-100"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={status === "loading"}
            className="flex items-center gap-2 rounded-xl bg-yellow-400 px-5 py-2.5 text-sm font-semibold text-zinc-950 transition-colors hover:bg-yellow-300 disabled:opacity-60"
          >
            {status === "loading" && <Loader2 className="size-4 animate-spin" />}
            {status === "loading" ? "Saving…" : "Confirm"}
          </button>
        </div>
      </form>
    </Modal>
  )
}

function DeleteAccountModal({ onClose, onDeleted }: { onClose: () => void; onDeleted: () => void }) {
  const [text, setText] = useState("")
  const [deleting, setDeleting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const canDelete = text.trim().toUpperCase() === "DELETE" && !deleting

  async function handleDelete() {
    setError(null)
    setDeleting(true)
    try {
      await api.deleteMe()
      onDeleted()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not delete the account.")
      setDeleting(false)
    }
  }

  return (
    <Modal onClose={onClose}>
      <div className="mb-4 flex items-start justify-between">
        <div className="flex items-center gap-3">
          <span className="flex size-10 items-center justify-center rounded-xl bg-red-500/15 text-red-400">
            <Trash2 className="size-5" />
          </span>
          <div>
            <h3 className="text-base font-semibold text-zinc-50">Delete account</h3>
            <p className="mt-0.5 text-xs text-zinc-500">This action is permanent.</p>
          </div>
        </div>
        <button
          type="button"
          onClick={onClose}
          aria-label="Close"
          className="rounded-lg p-1 text-zinc-500 transition-colors hover:bg-zinc-800 hover:text-zinc-100"
        >
          <X className="size-4" />
        </button>
      </div>

      <p className="mb-4 text-sm text-zinc-400">
        Type <span className="font-mono font-semibold text-red-400">DELETE</span> to confirm the
        permanent removal of your wallet.
      </p>
      <input
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="DELETE"
        className="mb-5 w-full rounded-xl border border-zinc-700 bg-zinc-950/60 px-3 py-2.5 font-mono text-sm text-zinc-100 outline-none transition-shadow focus:border-red-500 focus:ring-2 focus:ring-red-500/30"
      />

      <div className="flex items-center justify-end gap-3">
        <button
          type="button"
          onClick={onClose}
          className="rounded-xl px-4 py-2.5 text-sm font-medium text-zinc-400 transition-colors hover:text-zinc-100"
        >
          Cancel
        </button>
        <button
          type="button"
          disabled={!canDelete}
          onClick={handleDelete}
          className="rounded-xl bg-red-500 px-5 py-2.5 text-sm font-semibold text-zinc-50 transition-colors hover:bg-red-600 disabled:cursor-not-allowed disabled:opacity-40"
        >
          Delete permanently
        </button>
      </div>
    </Modal>
  )
}
