// Package ai implements AI player strategy สำหรับเกม Cashflow
// ใช้ engine API เดียวกับผู้เล่นจริง → รันได้ทั้งฝั่ง server และ browser (WASM)
package ai

import "github.com/finance-boardgame/engine/domain"

// Strategy คือ interface ของกลยุทธ์ AI — เลือก action ถัดไปจาก game state
type Strategy interface {
	Decide(state domain.GameState, player domain.Player) domain.Action
}

// GreedyStrategy เป็น baseline: พยายามซื้อ asset ที่ให้ cashflow ดีที่สุดเท่าที่เงินจะพอ
// TODO(Session#3): implement จริง (ประเมิน ROI, เก็บ emergency cash, หลีก doodad ที่ไม่คุ้ม)
type GreedyStrategy struct{}

var _ Strategy = (*GreedyStrategy)(nil)

// Decide คืน action ถัดไป (ตอนนี้เป็น stub — จบเทิร์นเสมอ)
func (s *GreedyStrategy) Decide(_ domain.GameState, player domain.Player) domain.Action {
	return domain.Action{PlayerID: player.ID, Type: domain.ActionEndTurn}
}
