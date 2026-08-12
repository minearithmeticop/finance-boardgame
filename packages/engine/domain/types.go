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

// TileType — ประเภทช่องบนกระดาน (ใช้ร่วมระหว่าง Rat Race และ Fast Track)
//
// หัวใจของบทพลิกเศรษฐกิจในเกม = สามเสา:
//   - Opportunity (โอกาส)  → บวก (ซื้อสินทรัพย์)
//   - Shopping (ช้อป)       → ล่อใจ (ใช้จ่ายฟุ่มเฟือย)
//   - Crisis (วิกฤต)         → ลบ (เหตุการณ์ร้าย)
type TileType int

const (
	TilePayday TileType = iota
	TileOpportunity // 🃏 โอกาสลงทุน (บวก — ตัดสินใจซื้อสินทรัพย์)
	TileShopping    // 🛍️ ช้อป/ใช้จ่าย (ลบ)
	TileCrisis      // ⚠️ วิกฤต (ลบ)
	TileMarket      // 📈 เหตุการณ์ตลาด (สำรอง — ยังไม่ลง board)
	TileDownsizing  // 📉 ตกงาน (สำรอง)
	TileFamily      // 👨‍👩‍👧 ครอบครัว
	TileDonate      // ❤️ บริจาค
	TileBlank       // · ช่องว่าง (สำรอง)
	TileNews        // 📰 ข่าวสาร (flavor — ไม่เปลี่ยนเงิน)
	TileWindfall    // 🎁 ได้รับเงิน (บวก)
	TileSideJob     // 💼 งานเสริม (บวก)
	TileLearn       // 📚 เรียนรู้/พัฒนาตัวเอง (flavor)
	TileHealth      // 🩺 ตรวจสุขภาพ (ลบเล็ก)
)

// หมวด LifeEvent — key สำหรับเลือก deck ใน cards.DrawLifeEvent
const (
	CatNews     = "news"
	CatWindfall = "windfall"
	CatSideJob  = "sidejob"
	CatShopping = "shopping"
	CatFamily   = "family"
	CatDonate   = "donate"
	CatLearn    = "learn"
	CatHealth   = "health"
	CatCrisis   = "crisis"
)

// LifeCategory คืนหมวด LifeEvent ของ tile นี้ ("" ถ้าไม่ใช่ life-event tile
// เช่น Payday/Opportunity หรือ tile สำรอง)
func (t TileType) LifeCategory() string {
	switch t {
	case TileNews:
		return CatNews
	case TileWindfall:
		return CatWindfall
	case TileSideJob:
		return CatSideJob
	case TileShopping:
		return CatShopping
	case TileFamily:
		return CatFamily
	case TileDonate:
		return CatDonate
	case TileLearn:
		return CatLearn
	case TileHealth:
		return CatHealth
	case TileCrisis:
		return CatCrisis
	default:
		return ""
	}
}

// Tile — ช่องหนึ่งช่องบนกระดาน
type Tile struct {
	Type TileType
	Name string
}

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

// DealCard — การ์ดดีลจากช่อง Opportunity (เสนอให้ซื้อสินทรัพย์)
type DealCard struct {
	Title       string
	AssetType   AssetType
	DownPayment Money // เงินดาวน์ที่ต้องจ่ายตอนซื้อ
	Cost        Money // ราคาเต็มของสินทรัพย์
	CashFlow    Money // รายได้ต่อเดือนที่สินทรัพย์นี้จะสร้าง
	LoanPayment Money // ค่าผ่อนหนี้ต่อเดือน (ถ้ากู้ = Cost − DownPayment)
}

// NetCashFlow คืนผลกระทบต่อเงินสดสุทธิ/เดือน ถ้าซื้อดีลนี้ (CashFlow − LoanPayment)
func (c DealCard) NetCashFlow() Money { return c.CashFlow - c.LoanPayment }

// DoodadCard/CrisisCard ถูกแทนที่ด้วย LifeEvent (ด้านล่าง) — รวมทุก life-event ใน type เดียว

// LifeEvent — การ์ดเหตุการณ์ในชีวิต ใช้กับทุก tile ที่ไม่ใช่ Opportunity/Payday
// resolve = Cash += Amount (บวก=ได้เงิน, ลบ=เสียเงิน, 0=แค่ข่าวสาร/flavor)
type LifeEvent struct {
	Category string // หมวด (CatNews/CatWindfall/...) — ใช้เลือก deck
	Title    string
	Detail   string
	Amount   Money // บวก/ลบ/ศูนย์
}

// PendingDecision — ดีลที่กำลังรอผู้เล่นตัดสินใจ (ซื้อ/ผ่าน) หลังตกช่อง Opportunity
type PendingDecision struct {
	PlayerID string
	DealCard DealCard
}

// Profession — อาชีพเริ่มต้น กำหนดเงินเดือนและรายจ่าย/หนี้สินพื้นฐาน
type Profession struct {
	Name           string
	Salary         Money
	Taxes          Money // ภาษีเงินได้/เดือน (คำนวณจากขั้นบันได ตอนสร้างอาชีพ)
	SocialSecurity Money // ประกันสังคมหัก ณ ที่จ่าย (5% capped)
	OtherExpenses  Money
	HomeMortgage   Liability
	SchoolLoan     Liability
	CarLoan        Liability
	CreditCard     Liability
	Savings        Money // เงินสดเริ่มต้น
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

	// รายจ่ายหัก ณ ที่จ่าย (breakdown ให้เห็นชัด)
	Tax            Money // ภาษีเงินได้/เดือน
	SocialSecurity Money // ประกันสังคม/เดือน

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
	CurrentTurn int              // index ของผู้เล่นที่กำลังเล่น
	Round       int              // รอบที่เท่าไหร่
	Seed        int64            // seed ของ RNG (เพื่อ replay)
	Pending     *PendingDecision // ดีลที่รอตัดสินใจ (ไม่ nil = เกมจอดอยู่ที่ decision phase)
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
	ActionDecline // ผ่านดีลที่กำลังตัดสินใจ (resolve PendingDecision)
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
	// Data คือ payload ยืดหยุ่น (เช่น {title, amount} ของการ์ด) — UI ใช้แสดงเรื่องราว
	Data map[string]any
}
