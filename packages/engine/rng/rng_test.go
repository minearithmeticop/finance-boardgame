package rng

import "testing"

// TestDeterminism ยืนยันว่า seed เดียวกัน → ลำดับผลเหมือนกันทุกตัว (จำเป็นสำหรับ replay/test)
func TestDeterminism(t *testing.T) {
	a := New(42)
	b := New(42)

	for i := 0; i < 100; i++ {
		ra := a.RollDice(6)
		rb := b.RollDice(6)
		if ra != rb {
			t.Fatalf("roll %d diverged: %d vs %d (determinism broken)", i, ra, rb)
		}
	}
}

func TestRollDiceInRange(t *testing.T) {
	g := New(7)
	for i := 0; i < 1000; i++ {
		r := g.RollDice(6)
		if r < 1 || r > 6 {
			t.Fatalf("roll out of range [1,6]: %d", r)
		}
	}
}

func TestRollDiceN(t *testing.T) {
	g := New(1)
	// ทอย 2 ลูก d6 → ผลรวมต้องอยู่ใน [2, 12]
	for i := 0; i < 100; i++ {
		s := g.RollDiceN(2, 6)
		if s < 2 || s > 12 {
			t.Fatalf("2d6 sum out of range [2,12]: %d", s)
		}
	}
}
