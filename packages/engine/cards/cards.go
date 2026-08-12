// Package cards เก็บข้อมูลการ์ดของเกม — DealCard (โอกาสซื้อสินทรัพย์)
// และ LifeEvent (เหตุการณ์ในชีวิตทุกประเภท) + ตัวจั่วการ์ด
//
// ออกแบบเนื้อหาให้หลากหลาย (ไม่ซ้ำจนเบื่อ) และสมดุล (บวก/ลบ/โอกาสผสม)
// ตัวเลขเป็นค่าประมาณแบบ game-simplified
package cards

import (
	"github.com/finance-boardgame/engine/domain"
	"github.com/finance-boardgame/engine/rng"
)

// le เป็นตัวช่วยสร้าง LifeEvent แบบสั้น (positional args ของ package ตัวเอง → vet-clean)
func le(category, title, detail string, amount domain.Money) domain.LifeEvent {
	return domain.LifeEvent{Category: category, Title: title, Detail: detail, Amount: amount}
}

// DealCards คืนชุดการ์ดดีล (จากช่อง Opportunity) — ผสมดีลดี/แยก สอนประเมินดีล
func DealCards() []domain.DealCard {
	return []domain.DealCard{
		{Title: "หุ้นกองทุนรวม", AssetType: domain.AssetStock, DownPayment: 5_000, Cost: 5_000, CashFlow: 50, LoanPayment: 0},
		{Title: "หุ้นปันผล", AssetType: domain.AssetStock, DownPayment: 10_000, Cost: 10_000, CashFlow: 150, LoanPayment: 0},
		{Title: "เครือข่ายพันธมิตร (ธุรกิจเครือข่าย)", AssetType: domain.AssetBusiness, DownPayment: 5_000, Cost: 5_000, CashFlow: 300, LoanPayment: 0},
		{Title: "ธุรกิจออนไลน์", AssetType: domain.AssetBusiness, DownPayment: 8_000, Cost: 30_000, CashFlow: 1_200, LoanPayment: 600},
		{Title: "ซักผ้าหยอดเหรียญ", AssetType: domain.AssetBusiness, DownPayment: 20_000, Cost: 80_000, CashFlow: 1_500, LoanPayment: 800},
		{Title: "ร้านกาแฟ", AssetType: domain.AssetBusiness, DownPayment: 40_000, Cost: 200_000, CashFlow: 3_500, LoanPayment: 1_500},
		{Title: "กิจการร้านอาหาร", AssetType: domain.AssetBusiness, DownPayment: 60_000, Cost: 400_000, CashFlow: 5_000, LoanPayment: 2_500},
		{Title: "ที่ดินนอกเมือง", AssetType: domain.AssetRealEstate, DownPayment: 15_000, Cost: 15_000, CashFlow: 0, LoanPayment: 0},
		{Title: "บ้านเก่าเช่าถูก", AssetType: domain.AssetRealEstate, DownPayment: 10_000, Cost: 150_000, CashFlow: 500, LoanPayment: 1_200},
		{Title: "คอนโด 1 ห้องให้เช่า", AssetType: domain.AssetRealEstate, DownPayment: 30_000, Cost: 300_000, CashFlow: 2_500, LoanPayment: 1_800},
		{Title: "ทาวน์เฮ้าส์ให้เช่า", AssetType: domain.AssetRealEstate, DownPayment: 50_000, Cost: 600_000, CashFlow: 4_000, LoanPayment: 3_500},
		{Title: "อาคารพาณิชย์ให้เช่า", AssetType: domain.AssetRealEstate, DownPayment: 80_000, Cost: 1_500_000, CashFlow: 9_000, LoanPayment: 8_000},
	}
}

// lifeEvents — เหตุการณ์ในชีวิตแบ่งตามหมวด (key = domain.Cat*)
// Amount: บวก=ได้เงิน, ลบ=เสียเงิน, 0=ข่าวสาร/flavor
var lifeEvents = map[string][]domain.LifeEvent{
	domain.CatNews: {
		le(domain.CatNews, "CEO ประกาศปิดบริษัทลูก — หุ้นเทคโนโลยีร่วง", "", 0),
		le(domain.CatNews, "ความขัดแย้งภูมิภาคกดดันราคาพลังงาน", "", 0),
		le(domain.CatNews, "โควิดระลอกใหม่ — ธุรกิจท่องเที่ยวชะลอตัว", "", 0),
		le(domain.CatNews, "ประกาศเลือกตั้งใหญ่ — นักลงทุนรอความชัดเจน", "", 0),
		le(domain.CatNews, "ตลาดคริปโตเร่งระวังฟองสบู่", "", 0),
		le(domain.CatNews, "AI เข้ามาแทนงานบางอาชีพ — ปรับตัวหรือถูกแทน", "", 0),
		le(domain.CatNews, "รัฐประกาศนโยบายลดภาษีเงินได้", "", 0),
		le(domain.CatNews, "ดอกเบี้ยขึ้น — การกู้ยืมแพงขึ้น", "", 0),
		le(domain.CatNews, "ทองคำราคาพุ่งทะลุสถิติ", "", 0),
		le(domain.CatNews, "ราคาน้ำมันลด — ค่าขนส่งถูกลง", "", 0),
		le(domain.CatNews, "สินค้าเกษตรราคาดี — เกษตรกรมีรายได้", "", 0),
		le(domain.CatNews, "พายุเข้าทำลายพื้นที่เกษตร — ของแพง", "", 0),
		le(domain.CatNews, "บริษัทยักษ์ใหญ่ปลดพนักงานรอบใหม่", "", 0),
		le(domain.CatNews, "ดัชนีค่าครองชีพปรับตัวสูงขึ้น", "", 0),
		le(domain.CatNews, "รัฐเปิดให้ซื้อหวยผ่านแอปอย่างถูกกฎหมาย", "", 0),
	},
	domain.CatWindfall: {
		le(domain.CatWindfall, "เงินเดือนพิเศษปลายปี", "", 15_000),
		le(domain.CatWindfall, "คืนภาษี", "", 8_000),
		le(domain.CatWindfall, "ขายของเก่าออนไลน์ได้กำไร", "", 3_000),
		le(domain.CatWindfall, "ของขวัญวันเกิด", "", 2_000),
		le(domain.CatWindfall, "เงินรางวัลประกวด", "", 5_000),
		le(domain.CatWindfall, "เบี้ยขยัน", "", 1_500),
		le(domain.CatWindfall, "ส่วนลดค่าน้ำมันรถ", "", 500),
		le(domain.CatWindfall, "เงินคืนเดบิต", "", 1_000),
		le(domain.CatWindfall, "ขายผลผลิตในสวนหลังบ้าน", "", 4_000),
		le(domain.CatWindfall, "มรดกจากญาติ", "", 30_000),
	},
	domain.CatSideJob: {
		le(domain.CatSideJob, "รับงานฟรีแลนซ์วันหยุด", "", 3_000),
		le(domain.CatSideJob, "สอนพิเศษ", "", 4_000),
		le(domain.CatSideJob, "ขับรถรับจ้างช่วงค่ำ", "", 2_500),
		le(domain.CatSideJob, "ขายของออนไลน์", "", 3_500),
		le(domain.CatSideJob, "ออกตลาดนัดสุดสัปดาห์", "", 2_000),
		le(domain.CatSideJob, "รับจ้างทำบัญชี", "", 5_000),
		le(domain.CatSideJob, "งานแปล/เขียน", "", 3_000),
		le(domain.CatSideJob, "รับจ้างซ่อมคอมพิวเตอร์", "", 2_000),
	},
	domain.CatShopping: {
		le(domain.CatShopping, "นัดเพื่อนกินข้าวนอกบ้าน", "", -800),
		le(domain.CatShopping, "ดูหนัง + ขนม", "", -600),
		le(domain.CatShopping, "ช้อปเสื้อผ้าใหม่", "", -3_000),
		le(domain.CatShopping, "ซื้อรองเท้าใหม่", "", -1_500),
		le(domain.CatShopping, "อัปเกรดโทรศัพท์", "", -12_000),
		le(domain.CatShopping, "ซื้อของใช้ในบ้าน", "", -2_000),
		le(domain.CatShopping, "ซื้อเครื่องใช้ไฟฟ้า", "", -5_000),
		le(domain.CatShopping, "นัดหมายทำผม", "", -800),
		le(domain.CatShopping, "ซื้อของเล่นให้ลูก", "", -1_200),
		le(domain.CatShopping, "จองทริปเที่ยว", "", -8_000),
	},
	domain.CatFamily: {
		le(domain.CatFamily, "ลูกป่วยเข้าโรงพยาบาล", "", -5_000),
		le(domain.CatFamily, "เปิดเทอม — ค่าเทอม", "", -3_000),
		le(domain.CatFamily, "งานบวชญาติ", "", -3_000),
		le(domain.CatFamily, "ค่านมและของใช้เด็ก", "", -1_500),
		le(domain.CatFamily, "ทำบุญวันเกิดลูก", "", -2_000),
		le(domain.CatFamily, "ครอบครัวทานข้าวนอกบ้าน", "", -1_200),
	},
	domain.CatDonate: {
		le(domain.CatDonate, "ทอดกฐิน", "", -2_000),
		le(domain.CatDonate, "บริจาคโรงพยาบาล", "", -1_000),
		le(domain.CatDonate, "ช่วยเหลือผู้ประสบภัย", "", -1_500),
		le(domain.CatDonate, "บริจาคเพื่อการศึกษา", "", -1_000),
		le(domain.CatDonate, "ใส่บาตรทำบุญ", "", -500),
	},
	domain.CatLearn: {
		le(domain.CatLearn, "ซื้อหนังสือพัฒนาตัวเอง", "เรียนรู้เรื่องการเงิน", -500),
		le(domain.CatLearn, "คอร์สออนไลน์ฟรี", "เสริมทักษะใหม่", 0),
		le(domain.CatLearn, "สัมมนาการลงทุน", "ได้ไอเดียลงทุน", -1_500),
		le(domain.CatLearn, "อ่านหนังสือจบเล่ม", "ได้มุมมองใหม่", 0),
		le(domain.CatLearn, "สมัครเรียนภาษา", "ลงทุนในตัวเอง", -2_000),
		le(domain.CatLearn, "ฟังพอดแคสต์การเงิน", "เพิ่มความรู้", 0),
	},
	domain.CatHealth: {
		le(domain.CatHealth, "ตรวจสุขภาพประจำปี", "", -2_000),
		le(domain.CatHealth, "ทำฟัน", "", -3_000),
		le(domain.CatHealth, "ฉีดวัคซีนป้องกัน", "", -1_500),
		le(domain.CatHealth, "สมัครฟิตเนส", "", -1_500),
		le(domain.CatHealth, "ซื้อยาประจำบ้าน", "", -500),
		le(domain.CatHealth, "ตรวจสายตา + แว่นใหม่", "", -3_500),
	},
	domain.CatCrisis: {
		le(domain.CatCrisis, "รถเสีย", "", -8_000),
		le(domain.CatCrisis, "ค่าปรับจราจร", "", -2_000),
		le(domain.CatCrisis, "เครื่องใช้ในบ้านพัง", "", -5_000),
		le(domain.CatCrisis, "ป่วยหนักเข้าโรงพยาบาล", "", -25_000),
		le(domain.CatCrisis, "ตกงาน 3 เดือน", "", -15_000),
		le(domain.CatCrisis, "อุบัติเหตุ", "", -20_000),
		le(domain.CatCrisis, "บาดเจ็บจากการกีฬา", "", -3_000),
		le(domain.CatCrisis, "ทำฟันฉุกเฉิน", "", -7_000),
	},
}

// LifeEvents คืนสำเนาชุดเหตุการณ์ของหมวดที่กำหนด
func LifeEvents(category string) []domain.LifeEvent {
	deck := lifeEvents[category]
	out := make([]domain.LifeEvent, len(deck))
	copy(out, deck)
	return out
}

// LifeCategories คืนรายชื่อหมวด LifeEvent ทั้งหมด
func LifeCategories() []string {
	return []string{
		domain.CatNews, domain.CatWindfall, domain.CatSideJob, domain.CatShopping,
		domain.CatFamily, domain.CatDonate, domain.CatLearn, domain.CatHealth, domain.CatCrisis,
	}
}

// DrawDealCard จั่วการ์ดดีลแบบสุ่ม (deterministic ผ่าน rng)
func DrawDealCard(r *rng.RNG) domain.DealCard {
	deck := DealCards()
	if r == nil {
		return deck[0]
	}
	return deck[r.Intn(len(deck))]
}

// DrawLifeEvent จั่วเหตุการณ์ของหมวดที่กำหนดแบบสุ่ม
// (ถ้าหมวดไม่มี deck จะคืนเหตุการณ์ flavor ปลอม ๆ)
func DrawLifeEvent(r *rng.RNG, category string) domain.LifeEvent {
	deck := lifeEvents[category]
	if len(deck) == 0 {
		return domain.LifeEvent{Category: category, Title: "ไม่มีอะไรเกิดขึ้น", Amount: 0}
	}
	if r == nil {
		return deck[0]
	}
	return deck[r.Intn(len(deck))]
}
