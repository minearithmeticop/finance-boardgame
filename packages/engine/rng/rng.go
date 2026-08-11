// Package rng คือ deterministic random number generator สำหรับเกม
//
// ทำไมต้อง deterministic?
//   - replay: ใช้ seed เดียวกัน → เดินและทอยเต๋าเหมือนกันทุกอย่าง
//   - testing: สามารถเขียน unit test ที่ assert ผลลัพธ์ได้
//   - anti-cheat: ฝั่ง server เป็นผู้กำหนด seed เดียว
//
// สำคัญ: ใช้ *local* source (rand.NewSource) เท่านั้น เพราะ global math/rand
// ใน Go 1.20+ ถูก auto-seed แบบสุ่ม (non-deterministic)
package rng

import "math/rand"

// RNG wraps a local source so sequences are fully reproducible.
type RNG struct {
	r *rand.Rand
}

// New สร้าง RNG จาก seed ที่กำหนด
func New(seed int64) *RNG {
	return &RNG{r: rand.New(rand.NewSource(seed))}
}

// Seed คืนค่า seed ปัจจุบันของ source (เอาไว้ snapshot/replay)
func (g *RNG) Seed() int64 {
	// NOTE: rand.Source ไม่ expose seed; engine เก็บ seed แยกใน GameState อยู่แล้ว
	// ฟังก์ชันนี้ placeholder — ใช้ GameState.Seed เป็นหลัก
	return 0
}

// Intn คืนค่า random ในช่วง [0, n)
func (g *RNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	return g.r.Intn(n)
}

// RollDice ทอยลูกเต๋าลูกเดียว คืนค่าในช่วง [1, sides]
func (g *RNG) RollDice(sides int) int {
	if sides <= 0 {
		return 0
	}
	return g.r.Intn(sides) + 1
}

// RollDiceN ทอยลูกเต๋า count ลูก แต่ละลูก sides ด้าน คืนผลรวม
func (g *RNG) RollDiceN(count, sides int) int {
	total := 0
	for i := 0; i < count; i++ {
		total += g.RollDice(sides)
	}
	return total
}
