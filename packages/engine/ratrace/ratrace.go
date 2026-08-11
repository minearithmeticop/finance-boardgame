// Package ratrace implements Rat Race (วงจรหนูถีบจักร) game rules
// เป็นวงเล็กที่ผู้เล่นเริ่มเกม ต้องสะสม passive income จนเท่ากับหรือมากกว่า expenses
package ratrace

import (
	"github.com/finance-boardgame/engine/domain"
	"github.com/finance-boardgame/engine/rng"
)

// BoardSize — จำนวนช่องบนกระดาน Rat Race
// TODO(Session#3): ยืนยันตัวเลขและเลย์เอาต์ช่องจริงตาม Cashflow 101
const BoardSize = 24

// TileType — ประเภทช่องบนกระดาน Rat Race
type TileType int

const (
	TilePayday TileType = iota
	TileOpportunity
	TileDoodad
	TileMarket
	TileDownsizing
	TileBaby
	TileCharity
	TileBlank
)

// MovePlayer เดินผู้เล่นตามจำนวนลูกเต๋า แล้วคืนตำแหน่งปลายทาง (wraps รอบกระดาน)
//
// TODO(Session#3): trigger เหตุการณ์ของช่องปลายทาง (Payday ผ่าน/Opportunity/Doodad ฯลฯ)
func MovePlayer(p *domain.Player, steps int, _ *rng.RNG) int {
	if p == nil {
		return 0
	}
	p.Position = (p.Position + steps) % BoardSize
	if p.Position < 0 {
		p.Position += BoardSize
	}
	return p.Position
}
