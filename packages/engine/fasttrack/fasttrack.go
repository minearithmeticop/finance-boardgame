// Package fasttrack implements Fast Track (ทางด่วน) game rules
// เป็นวงใหญ่ที่ผู้เล่นเข้าหลังหลุดจาก Rat Race แล้ว ลงทุนดีลใหญ่เพื่อทำเป้าหมาย (Dream)
package fasttrack

import "github.com/finance-boardgame/engine/domain"

// BoardSize — จำนวนช่องบนกระดาน Fast Track
// TODO(Session#3): ยืนยันเลย์เอาต์จริง
const BoardSize = 48

// EnterFastTrack ส่งผู้เล่นเข้าสู่ Fast Track
// เกม Cashflow: ผู้เล่นจะได้รับเงินเริ่มต้น 50,000 ในวง Fast Track (TODO: ยืนยันกติกา)
func EnterFastTrack(p *domain.Player) {
	if p == nil {
		return
	}
	p.OnFastTrack = true
	p.Position = 0
}

// HasWon เช็คเงื่อนไขชนะใน Fast Track
// TODO(Session#3): implement จริง — ปกติคือ (a) ซื้อ Dream ที่เลือกไว้ หรือ (b) เพิ่ม passive income 50,000/เดือน
func HasWon(p domain.Player) bool {
	_ = p
	return false
}
