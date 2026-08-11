// Package domain เก็บ core types ของเกม (entities, value objects, state)
// ทุก type ในนี้ต้องเป็น plain data — ไม่มี logic และไม่พึ่งพา I/O
// เพื่อให้ serialize ข้าม boundary (server ↔ WASM ↔ JSON) ได้สะดวก
package domain

// Money — จำนวนเงินในเกม ใช้ integer เพื่อกัน floating-point error
// (หน่วย = เงินตรา เช่นดอลลาร์; ถ้าต้องการความละเอียดวินาทีให้เปลี่ยนเป็น cents ภายหลัง)
type Money int64

// Phase — ระยะของเกม
type Phase int

const (
	PhaseRatRace Phase = iota // วงจรหนูถีบจักร
	PhaseFastTrack            // ทางด่วน (ออกจาก Rat Race แล้ว)
	PhaseEnded                // เกมจบ
)

func (p Phase) String() string {
	switch p {
	case PhaseRatRace:
		return "RatRace"
	case PhaseFastTrack:
		return "FastTrack"
	case PhaseEnded:
		return "Ended"
	default:
		return "Unknown"
	}
}

// AssetType — ประเภทสินทรัพย์
type AssetType int

const (
	AssetStock AssetType = iota
	AssetRealEstate
	AssetBusiness
	AssetOther
)

// Asset — สินทรัพย์ (สิ่งที่นำเงินเข้ากระเป๋า)
type Asset struct {
	ID            string
	Type          AssetType
	Name          string
	CashFlow      Money // รายได้ต่อเดือนที่สินทรัพย์นี้สร้าง
	Cost          Money // มูลค่ารวม
	DownPayment   Money // เงินดาวน์ที่จ่าย
	LoanRemaining Money // หนี้ที่ผูกกับสินทรัพย์ (ถ้ามี)
}

// Liability — หนี้สิน (สิ่งที่นำเงินออกจากกระเป๋า)
type Liability struct {
	ID      string
	Name    string
	Payment Money // รายจ่ายต่อเดือน
	Balance Money // เงินต้นคงเหลือ
}

// Profession — อาชีพเริ่มต้น กำหนดเงินเดือนและรายจ่าย/หนี้สินพื้นฐาน
type Profession struct {
	Name          string
	Salary        Money
	Taxes         Money
	OtherExpenses Money
	HomeMortgage  Liability
	SchoolLoan    Liability
	CarLoan       Liability
	CreditCard    Liability
	Savings       Money // เงินสดเริ่มต้น
}

// Player — ผู้เล่นหนึ่งคนในเกม
type Player struct {
	ID          string
	Name        string
	IsAI        bool
	Cash        Money
	Profession  Profession
	Assets      []Asset
	Liabilities []Liability // หนี้สินที่เพิ่มระหว่างเกม (นอกเหนือจาก profession)
	Position    int         // ตำแหน่งบนกระดาน (index ของช่อง)
	OnFastTrack bool
	Bankrupt    bool
}

// FinancialStatement — มุมมองงบการเงินที่คำนวณแล้ว (immutable view, derive จาก Player)
type FinancialStatement struct {
	// Income Statement
	EarnedIncome    Money // เงินเดือน
	PassiveIncome   Money // รายได้จากอสังหา/ธุรกิจ
	PortfolioIncome Money // รายได้จากหุ้น
	TotalIncome     Money

	TotalExpenses   Money
	MonthlyCashFlow Money // TotalIncome - TotalExpenses

	// Balance Sheet
	TotalAssets      Money
	TotalLiabilities Money
}

// GameState — snapshot สถานะเกมทั้งหมด (serialize ได้ทั้งเก็บ DB / ส่งผ่าน WS / WASM)
type GameState struct {
	Phase       Phase
	Players     []Player
	CurrentTurn int   // index ของผู้เล่นที่กำลังเล่น
	Round       int   // รอบที่เท่าไหร่
	Seed        int64 // seed ของ RNG (เพื่อ replay)
}

// ActionType — ประเภทคำสั่งที่ผู้เล่นส่งเข้า engine (command)
type ActionType int

const (
	ActionRoll ActionType = iota
	ActionBuyAsset
	ActionSellAsset
	ActionPayOffLiability
	ActionTakeLoan
	ActionEndTurn
)

// Action — คำสั่งที่ส่งเข้า engine
// TODO(Session#3): เปลี่ยน Payload จาก map[string]any เป็น typed struct ต่อ ActionType
type Action struct {
	PlayerID string
	Type     ActionType
	Payload  map[string]any
}

// EventType — เหตุการณ์ที่เกิดจาก action (event-sourcing style)
type EventType int

const (
	EventMoved EventType = iota
	EventLanded
	EventCashChanged
	EventAssetBought
	EventAssetSold
	EventPayday
	EventEscapedRatRace
	EventGameWon
)

// Event — ผลลัพธ์ของ action ที่ broadcast ให้ทุกคนในห้อง
type Event struct {
	Type     EventType
	PlayerID string
	// TODO(Session#3): typed payload
}
