# System & Game Requirements

## 🎯 High-Level Objective
พัฒนา **เกมกระดานการเงิน (Finance Boardgame)** แบบดิจิทัล โดยอ้างอิงและประยุกต์ใช้แนวคิดจาก **Cashflow Game** รองรับการเล่นหลายโหมดทั้ง online และ offline บน Web + PWA

## 🎮 Game Modes (Functional Scope)
| โหมด | รายละเอียด | ความต้องการด้านโครงสร้าง |
|---|---|---|
| 🌐 Online Multiplayer | ผู้เล่นหลายคนเล่นเรียลไทม์ผ่านเน็ต | WebSocket server + lobby/room |
| 🤖 Single Player (vs AI) | เล่นคนเดียวกับ AI | Engine + AI logic (รันฝั่ง browser ผ่าน WASM ได้) |
| 🎮 Local Pass-and-Play | หลายคนสลับเล่นที่จอเดียว | รัน engine ฝั่ง browser ได้โดยไม่ต้องต่อ server |
| ✉️ Async Multiplayer | เล่นคนละเวลากันได้ | Persistence + turn notification |

## 📑 Functional Requirements

### Game Mechanics
- [ ] **FR-001**: ระบบจำลองงบการเงิน (Income Statement & Balance Sheet)
- [ ] **FR-002**: วงจร Rat Race (หนูถีบจักร) และ Fast Track (ทางด่วน)
- [ ] **FR-003**: ระบบคำนวณ Cash Flow (รายได้ − รายจ่าย) และ Passive Income
- [ ] **FR-004**: ระบบลูกเต๋า (seeded/deterministic เพื่อรองรับ replay และ test)
- [ ] **FR-005**: ระบบเดินบนกระดาน (Rat Race board + Fast Track board)
- [ ] **FR-006**: ระบบการ์ด (Profession, Small Deal, Big Deal, Doodad, Market, Opportunity)
- [ ] **FR-007**: เงื่อนไขชนะ — Escape Rat Race (`Passive Income ≥ Total Expenses`) และ Fast Track goal (Dream)

### Game Modes / Multiplayer
- [ ] **FR-101**: Lobby/Room สำหรับ Online Multiplayer (สร้าง/เข้าร่วม/เริ่มเกม)
- [ ] **FR-102**: ระบบ AI Player (กลยุทธ์พื้นฐาน สามารถเล่นแทนคนได้)
- [ ] **FR-103**: Local Pass-and-Play (หมุนเวียนผู้เล่นบนจอเดียว, รัน offline ได้)
- [ ] **FR-104**: Async Multiplayer (บันทึก state, แจ้งเตือนถึงตา)
- [ ] **FR-105**: ระบบ replay / history ของเกมที่จบแล้ว

### UI/UX
- [ ] **FR-201**: หน้าจองบการเงินประจำตัว (interactive)
- [ ] **FR-202**: กระดานเกม + ตำแหน่งผู้เล่น
- [ ] **FR-203**: หน้ามือการ์ด / หน้าจอตัดสินใจดีล
- [ ] **FR-204**: PWA — ติดตั้งบนมือถือได้, รองรับ offline สำหรับ Local/AI mode

## ⚙️ Non-Functional Requirements
- [ ] **NFR-001 ความถูกต้อง**: Engine เป็น deterministic (seeded RNG) — ผลเดียวกันเสมอเมื่อ input เหมือนกัน
- [ ] **NFR-002 Single Source of Truth**: กติกาเขียนที่เดียว (engine), ห้าม duplicate logic ที่ frontend
- [ ] **NFR-003 Anti-cheat**: ใน Online mode server เป็นผู้ validate action ทุกอย่าง
- [ ] **NFR-004 Portability**: Engine คอมไพล์ได้ทั้ง server binary และ WebAssembly
- [ ] **NFR-005 Performance**: WASM lazy-load ไม่บล็อกหน้าแรก
- [ ] **NFR-006 Maintainability**: Type sharing Go ↔ TypeScript ผ่าน codegen เพื่อกัน drift
- [ ] **NFR-007 Testability (TDD)**: Engine logic + API contract ที่คุยกับ frontend เขียนแบบ TDD (test ก่อน → โค้ดผ่าน) — รายละเอียดใน `AGENTS.md` § Testing Conventions
- [ ] **NFR-008 Contract Stability**: ผลลัพธ์ business และ output shape ของ API ต้อง stable แม้ internals/logic เปลี่ยน — คุมด้วย contract tests (engine-level + HTTP)

## 🛠 Tech Stack (ตกลงแล้วใน Session #2)
| Layer | Technology |
|---|---|
| Frontend UI | Next.js (App Router) + React 19 + TypeScript |
| Styling | Tailwind CSS v4 |
| State (Client) | Zustand |
| PWA | Next.js PWA (manifest + service worker) |
| Backend Service | Go (REST + WebSocket) |
| Game Engine | Go (pure, deterministic) — compile เป็น server binary + WASM |
| Database | PostgreSQL (game history) + Redis (active games) |
| AI Player | Go (อยู่ใน engine package) |
| Type Sharing | struct → TS codegen |
| Monorepo | pnpm workspaces (web) + `go.work` (engine/backend) |

## 🏗 Architecture (ตกลงแล้ว)
**"Universal Engine (WASM)"** — ดูไดอะแกรมและเหตุผลเต็มใน `NOTE.md` (Session #2)

- `packages/engine` (Go): pure + deterministic, คอมไพล์ทั้ง server binary และ WASM
- `apps/backend` (Go): REST + WS hub + persistence + AI host
- `apps/web` (Next.js): เลือก host engine ระหว่าง backend (online/async) หรือ WASM (local/AI)

---

## 🔄 Proposed Evolution — Session #3 (Lifetime Financial Simulation)

> ⚠️ **ยังไม่ล็อก** — เป็นวิสัยทัศน์ที่กำลังหารือ รอตัดสินใจ (ดู open questions ใน `NOTE.md` Session #3)
> FR/NFR ด้านบนยังเป็น baseline ระยะใกล้; ส่วนนี้คือทิศทางระยะกลาง-ยาว

### Epics ที่เสนอ
- [ ] **EPIC-L1 Lifetime Timeline**: ผู้เล่นมีอายุ เวลาเดินพร้อมกัน เกมจบที่อายุ 80+ พร้อม life review
- [ ] **EPIC-L2 Vitals System**: สุขภาพ/ความเครียด/พลังงาน/เวลา กระทบการตัดสินใจ (อายุมากใช้พลังงานน้อย)
- [ ] **EPIC-L3 Insurance Systems**: ประกันสังคม/ชีวิต/อุบัติเหตุ/สุขภาพ สอนจัดการความเสี่ยง
- [ ] **EPIC-L4 Realistic Tax**: ลดหย่อนภาษี จ่ายภาษีทุกปี วางแผนภาษี
- [ ] **EPIC-L5 Real Investing**: อ่านหุ้นจริง วางแผนการเงินจริง
- [ ] **EPIC-L6 Childhood Mode (เฟส 2)**: ย้อนวัยเด็ก/ประถม — เงินจากพ่อแม่ ความแตกต่างของการเติบโต

### ผลกระทบต่อ Architecture (ต้องคุยให้ชัด)
- engine ปัจจุบัน = turn-based + seeded RNG ยังไม่มี concept "อายุ/เวลา" → อาจต้องเพิ่ม
- "การหลุดวงหนู" เปลี่ยนจาก win-condition เป็น milestone หนึ่งในชีวิต
- ระบบ vitals/insurance/tax = domain types ใหม่ใน `packages/engine/domain`
