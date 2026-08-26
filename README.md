# QuitoCoin

Educational cryptocurrency with its own blockchain — Go backend, Next.js frontend, and browser-based Proof-of-Work mining. A hands-on project to demonstrate how blocks, transactions, UTXO, mempool, and consensus work in practice.

---

## Table of Contents

1. [Overview](#overview)
2. [How the Currency Works](#how-the-currency-works)
3. [Cryptocurrency Numbers](#cryptocurrency-numbers)
4. [Getting Started](#getting-started)
5. [Technologies — Basic Overview](#technologies--basic-overview)
6. [Architecture and Advanced Technologies](#architecture-and-advanced-technologies)
7. [Folder Structure](#folder-structure)
8. [API](#api)
9. [Tests](#tests)

---

## Overview

QuitoCoin is a simplified blockchain where each block is mined via **Proof-of-Work (SHA-256)**, transactions sit in a **mempool** until included in the next block, and balances are tracked with a **simplified UTXO model** (one balance per user). There is no halving or supply cap — the focus is educational clarity, not real monetary economics.

```
User creates transfer → Mempool (Redis List) → Miner finds valid nonce → Block persisted (Postgres) → Balances updated (Redis UTXO) → History indexed (Redis Sorted Set)
```

---

## How the Currency Works

### Block

Each block contains `index`, `previousHash`, `nonce`, `miner` (publicID), `reward`, `transactions`, and `hash`. The hash is computed as:

```
data = "<nonce>:<index>:<miner>:<reward>:<previousHash>:<txs>"
hash = SHA256(data)  // hex
```

Where `<txs>` is the concatenation of `"<from>:<amount>:<to>;"` for each transaction with 6 decimal places. A block is only valid if `hash` starts with `0000000` (difficulty 7).

The **genesis block** (`index 0`, `miner "system"`, `reward 0`, `previousHash ""`) is mined the same way — iterating `nonce` until the difficulty is met — and is created automatically on startup if the chain is empty.

### Transaction and Mempool

A transfer is created via `POST /transfer` and goes to the mempool (Redis List `mempool`). Each entry stores `from`, `to`, and `amount` (net amount, excluding fee). During mining, up to **3 transactions** are dequeued FIFO and included in the block.

### Simplified UTXO

Unlike Bitcoin (per-output UTXO), each user has **a single balance** at `utxo:<userId>` in Redis. Operations are atomic `Credit` and `Debit` with insufficient-balance and negative-amount validation.

### Mining

The `POST /blockchain/mine` flow (authenticated):

1. `PullFirst(3)` from mempool
2. `reward = 50 + 1 * len(txs)`
3. `TryToMineBlock` — fetches the last block, computes `SHA256(nonce + FormatBlockInput(...))`, validates the `0000000` prefix
4. On success: persists the block to Postgres, `DeleteFirst`, credits the miner, debits senders (`amount + 1`), credits recipients (`amount`), indexes each participant's history in `user_blocks:<publicID>`

The frontend mirrors the same computation in `public/lib/mining.ts` (`SHA-256(nonce + data)`) and mines in the browser with `batchSize 500` and cooperative yielding via `setTimeout(0)`. There is also a concurrent standalone CLI miner in `cmd/miner` using `runtime.NumCPU()` workers and `atomic.Int64`.

### Chain Validation

`ValidateChain` iterates all ordered blocks, recomputes each hash, and verifies chaining `previousHash == prev.hash`.

---

## Cryptocurrency Numbers

| Parameter | Value | Notes |
|---|---|---|
| **Difficulty (PoW)** | `7` leading zeros | Fixed, no retargeting |
| **Target prefix** | `"0000000"` | `SHA256` hex |
| **Base reward** | `50.0` | Per block |
| **Per-transaction reward** | `1.0` | Also the fee charged to the sender |
| **Total reward** | `50 + 1 × n` | `n` = tx count in block (0→50, 3→53) |
| **Transactions per block** | `3` | FIFO from mempool |
| **Transfer fee** | `1.0` | Sender pays `amount + 1`, recipient gets `amount`, miner collects fee via reward |
| **Minimum transfer amount** | `> 1.0` | `amount <= 1` is rejected (`ErrAmountBelowFee`) |
| **Max supply** | ∞ | No cap defined |
| **Halving** | Not implemented | Constant reward |
| **Ranking** | Top 10 | Highest UTXO balances |
| **Default history** | 100 blocks | Paginated |
| **Block listing** | 20 per page | `limit`/`offset` |
| **JWT expiry** | 24 hours | HS256 |
| **bcrypt cost** | `DefaultCost` (10) |  |
| **PublicID format** | `XXX-XXXX-XXXX-XXXX-XXXX` | 19 chars, uppercase hex derived from UUID v4 |
| **Password rules** | ≥ 8 chars + uppercase + lowercase + digit + special character |  |

---

## Getting Started

### Prerequisites

- Go 1.22.2+
- Node 18+ and pnpm (or npm/yarn)
- PostgreSQL and Redis running locally

### Configuration

```bash
cp .env.example .env
# edit DATABASE_URL, REDIS_URL and JWT_SECRET
```

```env
SERVER_PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/quitocoin?sslmode=disable
REDIS_URL=localhost:6379
JWT_SECRET=change-me-to-a-random-secret
```

### Backend

```bash
make dev        # go run cmd/main.go
make build && make run
make test       # go test ./test/*/* -v
```

The server starts at `http://localhost:8080` with Swagger at `/swagger/*any`. On startup it runs `ent.Schema.Create` (auto-migrate) and rebuilds Redis caches if needed.

### Frontend

```bash
cd public
pnpm install
pnpm dev        # http://localhost:3000
pnpm build && pnpm start
```

Relevant env vars: `NEXT_PUBLIC_API_URL` (default `http://localhost:8080`) and `NEXT_PUBLIC_MINING_DIFFICULTY` (default `7`).

## Technologies — Basic Overview

| Layer | Technology | Version |
|---|---|---|
| **Language (backend)** | Go | 1.22.2 |
| **HTTP framework** | Gin | 1.10.0 |
| **ORM / Migrations** | Ent | 0.13.1 |
| **Relational DB** | PostgreSQL (`lib/pq`) | 1.12.3 |
| **Cache / Mempool / UTXO** | Redis (`go-redis/v9`) | 9.3.0 |
| **Authentication** | JWT (`golang-jwt/jwt/v5`) | 5.3.1 |
| **Password hashing** | `x/crypto/bcrypt` | 0.23.0 |
| **IDs** | `google/uuid` v4 | 1.3.0 |
| **Config** | `godotenv` | 1.5.1 |
| **API docs** | Swagger (`swaggo/gin-swagger`) | 1.6.0 |
| **Language (frontend)** | TypeScript | 5.7.3 |
| **Framework (frontend)** | Next.js | 16.2.6 |
| **UI** | React 19, Tailwind 4.2, shadcn, lucide-react |  |
| **State / Fetch** | SWR | 2.5.0 |
| **Tests (fake Redis)** | miniredis/v2 | 2.38.0 |

---

## Architecture and Advanced Technologies

### Architectural Pattern — Clean / Hexagonal with Manual DI

```
cmd/main.go  (composition root)
  → internal/config          — env loading
  → ent + go-redis           — infrastructure (Postgres, Redis)
  → internal/provider        — JWT, PasswordHasher, IdGenerator
  → internal/domain/{block,transaction,utxo,user,userblock}
        model.go             — entities and domain constants
        repository.go        — interfaces + implementations (Ent/Redis)
        service.go           — pure business logic
  → internal/usecase/*       — multi-domain orchestration (application layer)
  → internal/handler/*       — HTTP adapters (Gin) + Swagger annotations
  → internal/middleware      — Auth (Bearer JWT)
  → internal/error           — typed sentinel errors
```

Strict dependency rule: `handler → usecase → service → repository → infra`. No cycles. Dependency injection is done manually via constructors (`NewService`, `NewRepository`, `NewXUseCase`) with explicit wiring in `cmd/main.go` — no container (wire/google).

### Domains

- **block** — `Block` entity (Ent schema with `hash` unique immutable, `index`, `previousHash`, `nonce`, `miner`, `reward float64`, `transactions JSON`, `createdAt`). `Service` with `CalculateHash` (`SHA256`), `CreateGenesisBlock` (nonce loop), `TryToMineBlock` (prefix validation), `ValidateChain`.
- **transaction / Mempool** — `Transaction{From, To, Amount}`. `MemPool` interface with `CreateTransaction`, `Push/Pop` FIFO via Redis `LPush/RPop` and an in-memory alternative for tests. JSON serialization.
- **utxo** — `map[userId]float32` over Redis `GET/SET utxo:<id>`. `GetAll` via `KEYS utxo:*`. `Service.Credit/Debit` with negative and insufficient-balance checks; `NotFound` is treated as zero balance on credit.
- **user** — `User` entity (id UUID, name, email unique, password bcrypt, publicID `XXX-XXXX-...`, createdAt). Email validation, password strength, and uniqueness checks.
- **userblock** — Secondary index in Redis Sorted Set `user_blocks:<publicID>` with member `"<role>:<index>"` and score `blockIndex` (`miner`/`sender`/`receiver`). Enables paginated per-user history via `ZRevRange`.

### Ent (ORM) and Migrations

Schemas in `ent/schema/block.go` and `ent/schema/user.go` generate code in `ent/` (`client.go`, `mutation.go`, `block_*.go`, `user_*.go`, etc.) via `ent generate`. Tables `blocks` and `users` with types mapped in `ent/migrate/schema.go`. Auto-migration on boot via `client.Schema.Create(ctx)` — no manual `.sql` files.

### UseCase Layer

Each use case is a struct with injected dependencies and an `Execute(ctx, ...)` method. Examples:

- `MineBlock` — coordinates `block`, `utxo`, `transaction.MemPool`, and `userblock` in sequence; the most complex example (pull mempool → compute reward → mine → persist → clear mempool → update balances → index history).
- `CreateTransfer` — validates rules (`To` exists, `Amount > 1`, `From != To`, etc.), computes `net = Amount - 1`, and enqueues.
- `Initializer` — idempotent: if `chainLength == 0` creates genesis; if any Redis cache is empty, clears everything and rebuilds UTXO + history by replaying blocks from Postgres.

### Handler and Middleware Layer

- **Gin** with `gin.Default()` + inline CORS middleware (mirrors `Origin`, `Allow-Credentials: true`, `204` for `OPTIONS`).
- **Auth middleware** (`internal/middleware/auth.go`) extracts `Authorization: Bearer <token>`, validates via `JWTProvider`, injects `*JWTClaims{UserID, PublicID}` into `c.Set("claims", ...)`.
- Each handler follows `func HandleX(uc *UseCase) gin.HandlerFunc { ShouldBindJSON → GetClaims → uc.Execute → mapError → JSON }`.
- **Error mapping** (`handler/helpers.go:mapError`): sentinels → HTTP status (`400` validation/nonce/hash, `401` credentials, `404` not found, `409` conflict, `422` insufficient balance, `500` fallback).
- **Swagger** generated by `swaggo/swag` in `docs/` and served at `GET /swagger/*any`.

### Providers

- `JWTProvider` — HS256, `GenerateToken(userID, publicID)` with `24h` expiry, `ValidateToken` via `jwt.ParseWithClaims`.
- `IdGenerator` — `Generate()` = `uuid.NewString()`; `GeneratePublic()` = hyphenless UUID, truncated and formatted as uppercase `XXX-XXXX-XXXX-XXXX-XXXX`.
- `PasswordHasher` — `bcrypt.GenerateFromPassword(DefaultCost)` and `CompareHashAndPassword`.

### Error Handling

Centralized sentinels in `internal/error/` (`block.go`, `utxo.go`, `user.go`, `general.go`) — e.g., `ErrInvalidNonce`, `ErrInsufficientBalance`, `ErrEmailExists`. Repositories map `ent.IsNotFound → ErrXNotFound` and `ent.IsConstraintError → ErrXExists`; services propagate with `errors.Is`; handlers translate to HTTP. No `panic` on the happy path.

### Frontend — Details

- **Next.js 16 App Router** (`public/app/`), components in `public/components/quito/` and primitives in `public/components/ui/`.
- **API client** (`public/lib/api.ts`): `API_URL` via env, token in cookie `quito.token` (`Max-Age 7 days`, `SameSite=Lax`), `fetch` wrapper with CORS handling and `ApiError` class.
- **Browser mining** (`public/lib/mining.ts`): `DIFFICULTY` via `NEXT_PUBLIC_MINING_DIFFICULTY`, `hashAttempt(nonce, data) = SHA256(nonce+data)` hex, batched loop of 500 with cooperative `setTimeout(0)` and `hashRate` metric every 80ms.
- **Styling**: Tailwind 4.2 + `@tailwindcss/postcss`, `class-variance-authority`, `clsx`/`tailwind-merge`, `tw-animate-css`.

### Testing

- Command: `make test` → `go test ./test/*/* -v`.
- Framework: standard `testing` + `assertNoError`/`assertErrorIs` helpers.
- **Hand-rolled mocks** per domain (`test/block/mock_repository.go`, `test/utxo/mock_repository.go`, `test/user/mock_repository.go`) with `CreateFn`, `GetByHashFn` fields and call counters — no `testify/mock` or `mockgen`.
- **Redis repository tests** with `miniredis/v2` (in-memory Redis) in `test/utxo/repository_test.go`.
- **Mempool tests** with `transaction.NewInMemoryRepository` and `test/transaction/mempool_test.go`.
- **Extensive service coverage** (`test/block/service_test.go` ~710 lines): genesis, mining, hash computation, chain validation (tampered hash, broken link), repository error cases; uses `block.DifficultyPrefix = ""` trick to bypass PoW in unit tests.

### Infrastructure

Postgres and Redis are expected locally (no `Dockerfile`/`docker-compose`). Redis keys: `mempool` (List), `utxo:<id>` (String float), `user_blocks:<publicID>` (Sorted Set). Ent uses `sql` dialect via `lib/pq`.

---

## Folder Structure

```
.
├── cmd/
│   ├── main.go          # composition root, DI wiring, boot
│   └── miner/main.go    # concurrent PoW miner (CLI)
├── internal/
│   ├── config/          # .env loading
│   ├── domain/
│   │   ├── block/       # model, service (hash/PoW), repository (Ent)
│   │   ├── transaction/ # model, mempool (Redis + in-memory)
│   │   ├── utxo/        # model, service (Credit/Debit), repository (Redis)
│   │   ├── user/        # model, DTOs, service, repository (Ent)
│   │   └── userblock/   # history index (Redis Sorted Set)
│   ├── error/           # typed sentinels
│   ├── handler/         # Gin handlers + router + helpers + Swagger annotations
│   ├── middleware/      # Auth (JWT Bearer)
│   ├── provider/        # JWT, ID, Password
│   └── usecase/         # use cases (auth, transfer, mine, ranking, etc.)
├── ent/
│   ├── schema/          # Block and User schemas
│   ├── migrate/         # generated DDL
│   └── *.go             # generated code (client, mutations, queries)
├── docs/                # generated Swagger (docs.go, swagger.json/yaml)
├── public/              # Next.js frontend
│   ├── app/             # App Router (page, layout, dashboard)
│   ├── components/      # quito/* + ui/*
│   └── lib/             # api.ts, mining.ts, utils
├── test/
│   ├── block/           # service_test + mock
│   ├── transaction/     # mempool_test
│   ├── utxo/            # service_test, repository_test + mock
│   └── user/            # service_test + mock
├── Makefile
├── go.mod / go.sum
└── .env.example
```

---

## API

Base `http://localhost:8080` · Swagger at `/swagger/index.html` · Auth via `Authorization: Bearer <JWT>`.

| Method | Route | Auth | Description |
|---|---|---|---|
| `POST` | `/auth/register` | — | Create user, returns JWT |
| `POST` | `/auth/login` | — | Login, returns JWT |
| `GET` | `/me` | ✓ | Authenticated user profile |
| `PUT` | `/me` | ✓ | Update name/email |
| `PUT` | `/me/password` | ✓ | Change password |
| `DELETE` | `/me` | ✓ | Delete account |
| `POST` | `/blockchain/mine` | ✓ | Mine next block (requires valid `nonce`) |
| `GET` | `/blockchain/next-block` | ✓ | Next block data for mining |
| `GET` | `/blockchain/blocks` | ✓ | List paginated blocks |
| `GET` | `/blockchain/blocks/:index` | ✓ | Block by index |
| `GET` | `/blockchain/history` | ✓ | Authenticated user's history |
| `GET` | `/blockchain/history/:public_id` | ✓ | Another user's history |
| `POST` | `/transfer` | ✓ | Create transfer (goes to mempool) |
| `GET` | `/transfer/pending` | ✓ | Pending transactions (mempool) |
| `GET` | `/transfer/pending/me` | ✓ | Authenticated user's pending txs |
| `GET` | `/ranking` | ✓ | Top 10 richest (UTXO) |

---

## Tests

```bash
make test
# or
go test ./test/*/* -v
```

Mocks are hand-rolled per domain; Redis repository tests use `miniredis` with no real Redis required.

---

License: educational use.
