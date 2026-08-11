# Finance Boardgame 🎲💰

เกมกระดานการเงินดิจิทัล อ้างอิงแนวคิดจาก **Cashflow Game** (Robert Kiyosaki) — สอน Financial Literacy ผ่านการจำลองงบการเงินและการบริหาร Cashflow

## 📐 Architecture: "Universal Engine (WASM)"
Game Engine เขียนด้วย Go เป็น **pure + deterministic package** ตัวเดียว คอมไพล์ได้ทั้ง:
- 🟢 **Server binary** → รันใน `apps/backend` สำหรับ Online/Async mode
- 🟢 **WebAssembly** → รันใน browser สำหรับ Local/AI mode (offline ได้)

ดูรายละเอียดเต็มใน `NOTE.md` (Session #2) และ `REQUIREMENT.md`

## 📁 Monorepo Structure
```
finance-boardgame/
├── packages/
│   └── engine/              # 🎯 Game Engine (Go: pure + deterministic)
│       ├── domain/          # Player, GameState, Board, Tile, Card
│       ├── finance/         # คำนวณ Income Statement / Balance Sheet
│       ├── ratrace/         # กติกาวงนอก
│       ├── fasttrack/       # กติกาวงใน
│       ├── rng/             # seeded deterministic RNG
│       ├── ai/              # กลยุทธ์ AI player
│       └── cmd/wasm/        # build target → WASM สำหรับ browser
├── apps/
│   ├── backend/             # Go service: REST + WebSocket + persistence
│   └── web/                 # Next.js + PWA (UI และ engine host)
├── api-contracts/           # type sharing (Go struct → TS)
├── tooling/                 # build scripts (WASM, codegen)
└── docs/                    # NOTE / REQUIREMENT / AGENTS / INFO
```

## 🚀 Quick Start

### Prerequisites
- Go 1.24+
- Node.js 20+, pnpm 10+

### Game Engine & Backend (Go)
```bash
# ทุกคำสั่งรันจาก root (ใช้ go.work)
go build ./packages/engine/... ./apps/backend/...      # build engine + backend
go test ./packages/engine/...                          # run engine tests
cd apps/backend && go run ./cmd/server                 # start API server (:8080)
```

### Web (Next.js)
```bash
pnpm install
pnpm dev                           # start Next.js dev (build wasm ก่อนครั้งแรก)
pnpm build:wasm                    # compile engine → apps/web/public/wasm
```

## 📚 Documents
- [`NOTE.md`](./NOTE.md) — บันทึกการพูดคุย/การตัดสินใจ
- [`REQUIREMENT.md`](./REQUIREMENT.md) — System & Game Requirements
- [`AGENTS.md`](./AGENTS.md) — Conventions การทำงาน
- [`INFO.md`](./INFO.md) — Knowledge Base (ทฤษฎี Cashflow)
