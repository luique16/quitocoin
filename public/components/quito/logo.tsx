export function QuitoLogo({ className = "size-8" }: { className?: string }) {
  return (
    <svg
      viewBox="0 0 48 48"
      className={className}
      role="img"
      aria-label="Logo QuitoCoin"
      xmlns="http://www.w3.org/2000/svg"
    >
      {/* Coin */}
      <circle cx="24" cy="24" r="21" className="fill-yellow-400" />
      <circle
        cx="24"
        cy="24"
        r="21"
        className="fill-none stroke-yellow-200"
        strokeWidth="1.5"
      />
      <circle
        cx="24"
        cy="24"
        r="16.5"
        className="fill-none stroke-zinc-900/30"
        strokeWidth="1"
      />
      {/* Cockatiel crest */}
      <path
        d="M22 9 C20 12 19 14 20 17 M25 8.5 C24 12 23.5 14 24.5 17 M28 9.5 C27.5 12.5 27 14.5 27.5 17.5"
        className="fill-none stroke-zinc-900"
        strokeWidth="2"
        strokeLinecap="round"
      />
      {/* Head */}
      <path
        d="M17 20 C17 15 21 13 25 13 C31 13 34 17 34 22 C34 28 30 33 24 33 C19 33 16 29 16 25 Z"
        className="fill-zinc-900"
      />
      {/* Beak */}
      <path d="M34 22 L40 24 L34 26 Z" className="fill-zinc-900" />
      {/* Cheek */}
      <circle cx="27.5" cy="25" r="3.2" className="fill-red-500" />
      {/* Eye */}
      <circle cx="24" cy="20" r="1.6" className="fill-yellow-400" />
    </svg>
  )
}

export function QuitoWordmark({ className = "" }: { className?: string }) {
  return (
    <div className={`flex items-center gap-2.5 ${className}`}>
      <QuitoLogo className="size-9" />
      <div className="flex flex-col leading-none">
        <span className="text-base font-semibold tracking-tight text-zinc-50">
          QuitoCoin
        </span>
        <span className="font-mono text-[10px] uppercase tracking-widest text-yellow-400/80">
          QTC Network
        </span>
      </div>
    </div>
  )
}
