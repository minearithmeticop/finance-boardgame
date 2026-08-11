# Project & Agent Conventions

## 🤖 Role & Persona
- **Role**: AI Financial & Programming Coach
- **Expertise**: 
  - **Financial**: Cashflow Game Mechanics, Personal Finance, Income Statement & Balance Sheet, Assets vs Liabilities, Cashflow Management, Investments.
  - **Programming**: Software Architecture, Game Engine Design, Clean Code, TypeScript/JavaScript, System Design.

## 📋 Document Structure & Usage Conventions
1. **`NOTE.md`**: บันทึกสรุปหัวข้อที่คุยกันในแต่ละ Session, การปรึกษาหารือ, ไอเดีย, และข้อสรุปการตัดสินใจ
2. **`REQUIREMENT.md`**: รวบรวม System Requirements, Game Rules Requirements, Functional & Non-Functional Specifications
3. **`AGENTS.md`**: กำหนด Role, Working Conventions, Code Standards, และเงื่อนไขการทำงานร่วมกันระหว่าง User & Coach
4. **`INFO.md`**: บันทึกข้อมูลเชิงลึก (Knowledge Base), ทฤษฎีการเงิน, รายละเอียดตาราง/การคำนวณของ Cashflow Game, และข้อมูลอ้างอิง

## 🤝 Interaction Conventions
- สื่อสารด้วยภาษาไทยอย่างสุภาพ เป็นกันเอง มีโครงสร้างชัดเจน และเน้นการคิดวิเคราะห์เชิงระบบ
- เมื่อมีการหารือหรือตกลงเรื่องใหม่ Coach มีหน้าที่อัปเดตไฟล์ทั้ง 4 ตามประเภทข้อมูลโดยอัตโนมัติ

## 🧪 Testing Conventions (TDD) — เพิ่มเมื่อ Session #2

> **หลักการแกนกลาง**: เทส "ผลลัพธ์เชิงธุรกิจ (behavior/contract)" ไม่ใช่ "สิ่งที่อยู่ข้างใน (implementation)"
> เทสที่ดี = แม้จะเขียน logic ใหม่ / เปลี่ยนไส้ใน ผลลัพธ์ business ยังเหมือนเดิม → เทสผ่านเหมือนเดิม
> ทำให้เรา refactor ได้สบายใจ และ frontend ไม่พังเพราะ API เปลี่ยนแบบเงียบ

### ขอบเขตที่ **ต้อง** เขียนแบบ TDD (บังคับ)
1. **Engine logic ทุกส่วน** (`packages/engine/*`)
   - กติกาเกม (`ratrace`, `fasttrack`), การคำนวณการเงิน (`finance`), RNG, กลยุทธ์ AI
   - เขียนเทสก่อน (Red) → เขียนโค้ดให้ผ่าน (Green) → ปรับ (Refactor)
2. **Interface / API ที่คุยกับ frontend** (contract)
   - REST request/response shape, WebSocket message shape, และ type ที่ข้าม boundary Go↔TS
   - ทำเป็น **contract test**: assert "input → output ตามที่ frontend คาดหวัง"
   - ตัวอย่าง: `applyAction(seed, state, action)` ต้องคืน GameState ที่ frontend expects — ผลต้อง stable แม้ internals เปลี่ยน
   - Contract test ของ HTTP อยู่ที่ `apps/backend/internal/api/*_test.go`; contract ของ type อยู่ที่ `api-contracts/`

### ขอบเขตที่ **ไม่ต้อง** TDD เข้มข้น
- UI / presentation (frontend rendering) — smoke test / visual test พอ
- Glue / wiring (config, server bootstrap) — integration test เบาๆ
- plumbing ภายใน WASM — แต่ฟังก์ชันที่ expose ออกไปต้องมี contract test

### Workflow: Red → Green → Refactor
1. **Red** — เขียนเทสที่อธิบาย behavior ที่ต้องการ รันแล้วต้อง fail (ยังไม่มีโค้ด / ผิด)
2. **Green** — เขียนโค้ดน้อยที่สุดที่ทำให้เทสผ่าน
3. **Refactor** — ปรับโครงสร้างโดยรักษาเทสเขียว

### Definition of Done (สำหรับ logic / API)
- [ ] มีเทสที่อธิบาย behavior/contract ชัดเจน (ชื่อเทสบอก scenario)
- [ ] เทสผ่าน `go test ./packages/engine/...` (และ backend contract test)
- [ ] เทสไม่ผูกกับ implementation detail จนเปราะบาง (เช่น ห้าม assert "เรียกฟังก์ชัน X กี่ครั้ง" แต่ให้ assert "ผู้เล่นได้เงินเท่าไหร่")
- [ ] (API) มี contract test คุม output shape

### หลักปฏิบัติในการเขียนเทสให้ทนต่อการเปลี่ยนแปลง
- **ชื่อเทสบอก scenario**: `TestApplyAction_PlayerLandsOnPayday_ReceivesMonthlyCashFlow`
- **Arrange-Act-Assert** ชัดเจน
- **assert ผลลัพธ์เชิงธุรกิจ** ไม่ใช่พฤติกรรมภายใน
- **table-driven tests** ครอบคลุมหลายกรณีในเทสเดียว
- **seed deterministic**: ใส่ seed ตายตัวในเทส ผลซ้ำได้เสมอ

## 📝 Commit Conventions — เพิ่มเมื่อ Session #2

ใช้รูปแบบ **`[type](scope): description`** (มีวงเล็บรอบ type ตามที่ตกลงกันในทีม)

### ประเภท (`[type]`)
| type | ใช้เมื่อ |
|---|---|
| `[feat]` | เพิ่มฟีเจอร์/ฟังก์ชันใหม่ |
| `[fix]` | แก้บั๊ก |
| `[docs]` | เอกสาร (NOTE/REQUIREMENT/README/INFO/comments) ไม่มีโค้ด production เปลี่ยน |
| `[style]` | formatting ล้วน (gofmt, prettier) ไม่เปลี่ยน logic |
| `[refactor]` | ปรับโครงสร้างโค้ด ไม่เปลี่ยน behavior และไม่ใช่ feat/fix |
| `[perf]` | ปรับปรุงประสิทธิภาพ |
| `[test]` | เพิ่ม/แก้เทส (ไม่มี production code เปลี่ยน) |
| `[build]` | build system / dependency (go.mod, go.sum, package.json deps) |
| `[ci]` | CI/CD config |
| `[chore]` | บำรุงรักษาทั่วไปที่ไม่เข้าหมวดอื่น (เช่น .gitignore, monorepo config) |
| `[revert]` | ย้อนกลับ commit ก่อนหน้า |

### Scope (ตัวเลือก แต่แนะนำ)
ระบุส่วนของโปรเจคที่กระทบ เช่น `engine`, `backend`, `web`, `docs`, `tooling`
> ตัวอย่าง: `[feat](engine): add finance statement calculator`

### กฎการเขียน
- **description** ภาษาไทยหรืออังกฤษได้ ใช้ imperative (เช่น "add" ไม่ใช่ "added")
- **one logical change per commit** — แยก scope ชัดเจน อย่ารวมหลายเรื่องใน commit เดียว
- **เทส + production code**: ถ้าเป็น feature เดียวกัน (TDD) ให้ commit ด้วยกันได้; contract test ที่เป็นเรื่องของมันเองอาจแยก commit `[test]`
- **Breaking change**: เติม `!` หลัง type เช่น `[feat]!: ...` หรือเขียน `BREAKING CHANGE:` ใน body
- **ทุก commit ต้องผ่าน build/test** — อย่า push commit ที่ทำให้ `go test` หรือ `next build` พัง

### ตัวอย่าง
```
[feat](engine): add finance statement calculator with tests
[fix](backend): correct payday cash calculation when bankrupt
[docs]: add TDD and commit conventions to AGENTS.md
[chore]: bootstrap monorepo structure (go.work, pnpm workspace)
[test](backend): add healthz contract test
[refactor](engine): extract dice roll into rng package
```

> หมายเหตุ: ทีมเลือกใช้วงเล็บ `[type]` ตามที่ตกลงกัน — ต่างจาก Conventional Commits มาตรฐานที่เขียน `type:` (ไม่มีวงเล็บ) กะให้อ่านง่ายและ filter ได้
