"use client"

import { useEffect, useState } from "react"
import { Check, KeyRound, Loader2, Mail, Trash2, User, X } from "lucide-react"
import { USER } from "@/lib/quito-data"

export function AccountView({ onLogout }: { onLogout: () => void }) {
  const [name, setName] = useState(USER.name)
  const [email, setEmail] = useState("joao@quitocoin.dev")
  const [savedFlash, setSavedFlash] = useState(false)
  const [passwordOpen, setPasswordOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)

  function handleProfileSave(e: React.FormEvent) {
    e.preventDefault()
    setSavedFlash(true)
    setTimeout(() => setSavedFlash(false), 2200)
  }

  return (
    <div className="mx-auto max-w-3xl">
      <header className="mb-8">
        <p className="text-sm text-zinc-500">Configurações</p>
        <h1 className="text-2xl font-semibold tracking-tight text-zinc-50">Sua conta</h1>
      </header>

      {/* Perfil */}
      <section className="mb-6 rounded-2xl border border-zinc-800 bg-zinc-900 p-6">
        <h2 className="mb-1 text-sm font-semibold text-zinc-100">Perfil</h2>
        <p className="mb-5 text-xs text-zinc-500">Atualize suas informações públicas de identificação.</p>

        <form onSubmit={handleProfileSave} className="flex flex-col gap-4">
          <label className="flex flex-col gap-1.5">
            <span className="text-xs font-medium text-zinc-400">Nome</span>
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

          <div className="flex items-center gap-3">
            <button
              type="submit"
              className="rounded-xl bg-yellow-400 px-5 py-2.5 text-sm font-semibold text-zinc-950 transition-colors hover:bg-yellow-300"
            >
              Salvar alterações
            </button>
            {savedFlash && (
              <span className="flex items-center gap-1.5 text-xs font-medium text-emerald-400">
                <Check className="size-4" /> Perfil atualizado
              </span>
            )}
          </div>
        </form>
      </section>

      {/* Segurança */}
      <section className="mb-6 rounded-2xl border border-zinc-800 bg-zinc-900 p-6">
        <h2 className="mb-1 text-sm font-semibold text-zinc-100">Segurança</h2>
        <p className="mb-5 text-xs text-zinc-500">Gerencie as credenciais de acesso à sua carteira.</p>

        <div className="flex items-center justify-between rounded-xl border border-zinc-800 bg-zinc-950/40 px-4 py-3.5">
          <div className="flex items-center gap-3">
            <span className="flex size-9 items-center justify-center rounded-lg bg-zinc-800 text-yellow-400">
              <KeyRound className="size-4" />
            </span>
            <div>
              <p className="text-sm font-medium text-zinc-100">Senha</p>
              <p className="text-xs text-zinc-500">Última alteração há 3 meses</p>
            </div>
          </div>
          <button
            type="button"
            onClick={() => setPasswordOpen(true)}
            className="rounded-xl border border-zinc-700 px-4 py-2 text-sm font-medium text-zinc-100 transition-colors hover:border-yellow-400 hover:text-yellow-400"
          >
            Alterar senha
          </button>
        </div>
      </section>

      {/* Zona de perigo */}
      <section className="rounded-2xl border border-red-500/30 bg-red-500/5 p-6">
        <h2 className="mb-1 text-sm font-semibold text-red-400">Zona de perigo</h2>
        <p className="mb-5 text-xs text-zinc-400">
          A exclusão da conta é permanente e remove sua carteira da rede. Esta ação não pode ser desfeita.
        </p>
        <button
          type="button"
          onClick={() => setDeleteOpen(true)}
          className="flex items-center gap-2 rounded-xl bg-red-500 px-5 py-2.5 text-sm font-semibold text-zinc-50 transition-colors hover:bg-red-600"
        >
          <Trash2 className="size-4" />
          Deletar conta
        </button>
      </section>

      {passwordOpen && <ChangePasswordModal onClose={() => setPasswordOpen(false)} />}
      {deleteOpen && (
        <DeleteAccountModal onClose={() => setDeleteOpen(false)} onConfirm={onLogout} />
      )}
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

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (next.length < 6) {
      setStatus("error")
      setError("A nova senha deve ter ao menos 6 caracteres.")
      return
    }
    if (next !== confirm) {
      setStatus("error")
      setError("A confirmação não coincide com a nova senha.")
      return
    }
    setError("")
    setStatus("loading")
    setTimeout(() => {
      onClose()
    }, 1000)
  }

  return (
    <Modal onClose={onClose}>
      <div className="mb-5 flex items-start justify-between">
        <div>
          <h3 className="text-base font-semibold text-zinc-50">Alterar senha</h3>
          <p className="mt-0.5 text-xs text-zinc-500">Informe a senha atual e escolha uma nova.</p>
        </div>
        <button
          type="button"
          onClick={onClose}
          aria-label="Fechar"
          className="rounded-lg p-1 text-zinc-500 transition-colors hover:bg-zinc-800 hover:text-zinc-100"
        >
          <X className="size-4" />
        </button>
      </div>

      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <label className="flex flex-col gap-1.5">
          <span className="text-xs font-medium text-zinc-400">Senha atual</span>
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
          <span className="text-xs font-medium text-zinc-400">Nova senha</span>
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
          <span className="text-xs font-medium text-zinc-400">Confirmar nova senha</span>
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
            Cancelar
          </button>
          <button
            type="submit"
            disabled={status === "loading"}
            className="flex items-center gap-2 rounded-xl bg-yellow-400 px-5 py-2.5 text-sm font-semibold text-zinc-950 transition-colors hover:bg-yellow-300 disabled:opacity-60"
          >
            {status === "loading" && <Loader2 className="size-4 animate-spin" />}
            {status === "loading" ? "Salvando..." : "Confirmar"}
          </button>
        </div>
      </form>
    </Modal>
  )
}

function DeleteAccountModal({ onClose, onConfirm }: { onClose: () => void; onConfirm: () => void }) {
  const [text, setText] = useState("")
  const canDelete = text.trim().toUpperCase() === "DELETAR"

  return (
    <Modal onClose={onClose}>
      <div className="mb-4 flex items-start justify-between">
        <div className="flex items-center gap-3">
          <span className="flex size-10 items-center justify-center rounded-xl bg-red-500/15 text-red-400">
            <Trash2 className="size-5" />
          </span>
          <div>
            <h3 className="text-base font-semibold text-zinc-50">Deletar conta</h3>
            <p className="mt-0.5 text-xs text-zinc-500">Esta ação é permanente.</p>
          </div>
        </div>
        <button
          type="button"
          onClick={onClose}
          aria-label="Fechar"
          className="rounded-lg p-1 text-zinc-500 transition-colors hover:bg-zinc-800 hover:text-zinc-100"
        >
          <X className="size-4" />
        </button>
      </div>

      <p className="mb-4 text-sm text-zinc-400">
        Digite <span className="font-mono font-semibold text-red-400">DELETAR</span> para confirmar a
        remoção definitiva da sua carteira.
      </p>
      <input
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="DELETAR"
        className="mb-5 w-full rounded-xl border border-zinc-700 bg-zinc-950/60 px-3 py-2.5 font-mono text-sm text-zinc-100 outline-none transition-shadow focus:border-red-500 focus:ring-2 focus:ring-red-500/30"
      />

      <div className="flex items-center justify-end gap-3">
        <button
          type="button"
          onClick={onClose}
          className="rounded-xl px-4 py-2.5 text-sm font-medium text-zinc-400 transition-colors hover:text-zinc-100"
        >
          Cancelar
        </button>
        <button
          type="button"
          disabled={!canDelete}
          onClick={onConfirm}
          className="rounded-xl bg-red-500 px-5 py-2.5 text-sm font-semibold text-zinc-50 transition-colors hover:bg-red-600 disabled:cursor-not-allowed disabled:opacity-40"
        >
          Deletar permanentemente
        </button>
      </div>
    </Modal>
  )
}
