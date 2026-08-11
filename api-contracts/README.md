# API Contracts

โฟลเดอร์นี้เก็บนิยามของชนิดข้อมูลที่แชร์ระหว่าง **Go** (engine/backend) และ **TypeScript** (Next.js)
เพื่อกัน type drift (NFR-006)

## 🎯 เป้าหมาย
เนื่องจาก Go กับ TypeScript ไม่ share type อัตโนมัติ เราจึงต้องมี single source สำหรับ
data shape ที่ข้าม boundary สำคัญๆ ได้แก่:
- `GameState` (snapshot ที่ serialize ระหว่าง server ↔ browser)
- `Action` (คำสั่งที่ผู้เล่นส่งเข้า engine)
- `Event` (ผลลัพธ์ที่ engine broadcast)
- REST API request/response shapes

## 🔜 แนวทางที่จะใช้ (Session #3)
เลือกวิธีใดวิธีหนึ่ง:
1. **Go struct → TypeScript** (แนะนำ): ใช้เครื่องมือ generate เช่น
   [`gotype`](https://github.com/gotyped) หรือ custom generator จาก `go/ast`
2. **OpenAPI schema**: นิยาม REST API ที่นี่ แล้ว generate both Go server stub และ TS client
3. **Protobuf**: กรณีต้องการ binary protocol + strict schema (อาจมากเกินไปสำหรับ web board game)

> TODO(Session#3): ตัดสินใจวิธี แล้ววาง schema files ในโฟลเดอร์นี้ + เพิ่ม generator ใน `tooling/`

## 📦 ปัจจุบัน
ยังไม่มี contract ที่ generate แล้ว — type ฝั่ง TS ยังเป็น `unknown`/hand-written ชั่วคราว
จนกว่าจะเลือก generator
