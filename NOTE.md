# Session Notes (บันทึกการพูดคุย)

## 📌 Kickoff Session (2026-08-10)
- **ผู้ปรึกษา**: User
- **Coach**: AI Financial & Programming Coach
- **ข้อตกลงในการจัดเก็บข้อมูล**:
  - `NOTE.md` -> จดบันทึกสิ่งที่คุยกัน
  - `REQUIREMENT.md` -> เก็บ Requirements ของโปรเจกต์
  - `AGENTS.md` -> เก็บ Conventions และแนวทางการทำงาน
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

## 📌 Session #3 — Glossary + Lifetime-Simulation Vision (2026-08-12)
- **ผู้ปรึกษา**: User
- **Coach**: AI Financial & Programming Coach

### 🎯 สิ่งที่ทำ (deliverable)
- สร้าง **คัมภีร์การเงิน** เป็นหน้าเว็บ `/glossary` (data-driven, ค้นหา/กรองหมวดได้)
  - รวบรวมคำศัพท์ ~26 คำ พร้อมที่มาที่ไป + ตัวอย่าง + mapping ไป engine field
  - เป้าหมาย: ทีมและผู้เล่นเข้าใจตรงกัน; ผู้สร้างเรียนรู้ไปพร้อมสร้าง (User บอกเองว่ายังไม่เข้าใจ เช่น Taxes/ประกัน)

### 🗣 วิสัยทัศน์ใหม่ที่หารือ (ยังไม่ล็อก)
User แชร์วิสัยทัศน์ที่ขยายจาก Cashflow ไปสู่ **"จำลองชีวิตการเงินตลอดทั้งชีวิต"**:
- ผู้เล่น = จบใหม่สุ่มอาชีพ, อายุเริ่มต้นต่างกัน
- เวลาเดินพร้อมกัน แก่ไปด้วยกัน; **จบเมื่ออายุเกิน 80** แล้วรีวิวชีวิต
- ระบบใหม่: Vitals (สุขภาพ/ความเครียด/พลังงาน/เวลา), ประกัน (สังคม/ชีวิต/อุบัติเหตุ/สุขภาพ), ภาษีสมจริง (ลดหย่อน/จ่ายทุกปี), อ่านหุ้นจริง
- เฟส 2 (ทีหลัง): ระบบวัยเด็ก/ประถม (เงินจากพ่อแม่ ความแตกต่างของการเติบโต)
- เฟสตั้งใจเริ่มที่ผู้ใหญ่ก่อน: "ให้คุ้นชินระบบเกมและเข้าใจชีวิตก่อน แล้วค่อยย้อนวัยเด็ก"

### ❓ Open Questions (รอตัดสินใจ)
1. **ทิศทางหลัก**: ล็อกเป็น "Lifetime Sim" เลย หรือค่อย ๆ เพิ่ม realism บนฐาน Cashflow?
2. **โมเดลเวลาใน engine**: ปัจจุบัน engine เป็น turn-based + seeded RNG; lifetime-sim อาจต้องมี "นาฬิกา/อายุ" → กระทบ architecture
3. **Vertical slice แรก**: อยากได้ playable แบบไหนก่อน (Rat Race loop / lifetime skeleton / systems แยก)?
4. **จุดจบเกม**: รีวิวชีวิตที่อายุ 80+ ยังไง? คะแนน? บอกเล่า? เปรียบเทียบ?

### 📦 ผลลัพธ์ของ Session
- อัปเดต `INFO.md` (Game Design Vision), `REQUIREMENT.md` (Proposed Evolution), ไฟล์นี้
- สร้าง `apps/web/data/glossary.ts` + `apps/web/app/glossary/page.tsx` + nav ใน layout
