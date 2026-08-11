# Session Notes (บันทึกการพูดคุย)

## 📌 Kickoff Session (2026-08-10)
- **ผู้ปรึกษา**: User
- **Coach**: AI Financial & Programming Coach
- **ข้อตกลงในการจัดเก็บข้อมูล**:
  - `NOTE.md` -> จดบันทึกสิ่งที่คุยกัน
  - `REQUIREMENT.md` -> เก็บ Requirements ของโปรเจกต์
  - `AGENTS.md` / `AGENT.md` -> เก็บ Conventions และแนวทางการทำงาน
  - `INFO.md` -> เก็บข้อมูลความรู้ และ Data Reference
- **หัวข้อเริ่มต้น**: การหารือและวิเคราะห์เกี่ยวกับ **เกม Cashflow** (แนวคิดการเงิน, Game Mechanics, และการนำมาสร้าง Board Game)

## 📌 Session #2 — Tech Stack & Architecture (2026-08-12)
- **ผู้ปรึกษา**: User
- **Coach**: AI Financial & Programming Coach

### 🗣 หัวข้อที่คุย
- การเลือกภาษา/เทคโนโลยีสำหรับ frontend และ backend ของเกม Cashflow
- User ชำนาญ **Next.js** และ **Go** → ใช้สองภาษานี้เป็นหลักของโปรเจกต์

### ✅ ข้อตกลง / การตัดสินใจ
1. **Tech Stack หลัก**
   - Frontend: **Next.js (App Router) + PWA**, state ด้วย **Zustand**, styling ด้วย **Tailwind CSS v4**
   - Backend: **Go** (REST + WebSocket) + **PostgreSQL** (history) + **Redis** (active games)
   - Game Engine: **Go package (pure + deterministic)** ที่คอมไพล์ได้ทั้ง server binary และ **WebAssembly** เพื่อรันใน browser
2. **Architecture: "Universal Engine (WASM)"**
   - Engine ตัวเดียวรองรับครบทั้ง 4 โหมด (Online / Single vs AI / Local Pass-and-Play / Async)
   - หลักการ: single source of truth, กัน cheating, deterministic (seeded RNG) เพื่อ test/replay ง่าย
3. **โหมดการเล่นที่ต้องการ (ครบ 4 โหมด)**
   - 🌐 Online Multiplayer, 🤖 Single Player vs AI, 🎮 Local Pass-and-Play, ✉️ Async Multiplayer
4. **Platform**: Web + PWA (เฟสแรก)
5. **Scope เริ่มต้น**: ทำเกมเต็ม (Rat Race + Fast Track)
6. **โครงสร้างโปรเจค**: **Monorepo** (pnpm workspaces สำหรับ web + `go.work` สำหรับ engine/backend)
7. **Module path (ชั่วคราว)**: `github.com/finance-boardgame/{engine,backend}` — เปลี่ยนทีหลังได้ด้วย `go mod edit -module`
8. **Testing Convention (TDD)**: เพิ่มนโยบาย TDD ลง `AGENTS.md` — engine logic และ API contract ที่คุยกับ frontend **ต้องเขียนเทสก่อน**; เทสต้อง assert ผลลัพธ์เชิงธุรกิจที่ stable แม้ internals/logic เปลี่ยน (อัปเดต NFR-007 และเพิ่ม NFR-008)

### ⚠️ ข้อควรระวังที่จดไว้
- Go → WASM ไฟล์ค่อนข้างใหญ่ (~2-5MB) → lazy-load เฉพาะตอนเข้า Local/AI mode ไม่บล็อกหน้าแรก
- Engine ต้องเขียนแบบ **pure** (ห้าม `time.Now()`, ห้าม goroutine ที่จับ state ข้างนอก) มิฉะนั้น deterministic และการ compile WASM พัง
- Type sharing Go ↔ TypeScript ต้องมี **codegen** (จาก engine struct → TS interface) เพื่อกัน drift

### 📦 ผลลัพธ์ของ Session
- อัปเดต `REQUIREMENT.md` (เพิ่ม FR/NFR + Tech Stack + Architecture)
- Scaffold โครงสร้าง Monorepo: `packages/engine`, `apps/backend`, `apps/web`, `api-contracts`, `tooling`
