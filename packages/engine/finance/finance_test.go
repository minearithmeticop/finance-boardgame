package finance

import (
	"testing"

	"github.com/finance-boardgame/engine/domain"
	"github.com/finance-boardgame/engine/payroll"
	"github.com/finance-boardgame/engine/profession"
)

// helper สร้าง player พื้นฐานสำหรับทดสอบ
func samplePlayer(salary, taxes, other domain.Money, mortgagePayment, mortgageBalance domain.Money) domain.Player {
	return domain.Player{
		Cash: 1000,
		Profession: domain.Profession{
			Name:          "Tester",
			Salary:        salary,
			Taxes:         taxes,
			OtherExpenses: other,
			HomeMortgage:  domain.Liability{Payment: mortgagePayment, Balance: mortgageBalance},
		},
	}
}

func TestStatement_BasicCashFlow(t *testing.T) {
	p := samplePlayer(2000, 500, 200, 100, 3000)
	// เพิ่มสินทรัพย์อสังหา cashflow 100/เดือน
	p.Assets = []domain.Asset{
		{Type: domain.AssetRealEstate, CashFlow: 100, Cost: 5000},
	}

	fs := Statement(p)

	// Income
	if fs.EarnedIncome != 2000 {
		t.Errorf("EarnedIncome = %d, want 2000", fs.EarnedIncome)
	}
	if fs.PassiveIncome != 100 {
		t.Errorf("PassiveIncome = %d, want 100", fs.PassiveIncome)
	}
	if fs.TotalIncome != 2100 {
		t.Errorf("TotalIncome = %d, want 2100", fs.TotalIncome)
	}

	// Expenses = taxes(500) + other(200) + mortgage payment(100) = 800
	if fs.TotalExpenses != 800 {
		t.Errorf("TotalExpenses = %d, want 800", fs.TotalExpenses)
	}
	if fs.MonthlyCashFlow != 1300 {
		t.Errorf("MonthlyCashFlow = %d, want 1300", fs.MonthlyCashFlow)
	}

	// Balance sheet
	if fs.TotalAssets != 6000 { // cash 1000 + asset cost 5000
		t.Errorf("TotalAssets = %d, want 6000", fs.TotalAssets)
	}
	if fs.TotalLiabilities != 3000 {
		t.Errorf("TotalLiabilities = %d, want 3000", fs.TotalLiabilities)
	}
}

func TestStatement_StockGoesToPortfolioIncome(t *testing.T) {
	p := samplePlayer(1000, 0, 0, 0, 0)
	p.Assets = []domain.Asset{
		{Type: domain.AssetStock, CashFlow: 50},
		{Type: domain.AssetRealEstate, CashFlow: 200},
	}
	fs := Statement(p)

	if fs.PortfolioIncome != 50 {
		t.Errorf("PortfolioIncome = %d, want 50", fs.PortfolioIncome)
	}
	if fs.PassiveIncome != 200 {
		t.Errorf("PassiveIncome = %d, want 200", fs.PassiveIncome)
	}
}

func TestCanEscapeRatRace(t *testing.T) {
	// กรณี 1: passive น้อยกว่า expenses → ยังหนีไม่ได้
	p1 := samplePlayer(2000, 500, 200, 100, 3000)
	p1.Assets = []domain.Asset{{Type: domain.AssetRealEstate, CashFlow: 100}}
	if CanEscapeRatRace(p1) {
		t.Error("expected cannot escape (passive 100 < expenses 800)")
	}

	// กรณี 2: passive มากพอ → หนีได้
	p2 := samplePlayer(2000, 0, 0, 100, 3000) // expenses = 100
	p2.Assets = []domain.Asset{{Type: domain.AssetRealEstate, CashFlow: 150}}
	if !CanEscapeRatRace(p2) {
		t.Error("expected can escape (passive 150 >= expenses 100)")
	}
}

// TestStatement_IncludesSocialSecurityInExpenses — SS ต้องถูกนำรวมเข้า TotalExpenses
func TestStatement_IncludesSocialSecurityInExpenses(t *testing.T) {
	p := samplePlayer(20_000, 0, 0, 0, 0) // taxes/other/mortgage = 0
	p.Profession.SocialSecurity = 750

	fs := Statement(p)

	if fs.SocialSecurity != 750 {
		t.Errorf("SocialSecurity = %d, want 750", fs.SocialSecurity)
	}
	if fs.TotalExpenses != 750 {
		t.Errorf("TotalExpenses = %d, want 750 (เกิดจาก SS เพียงตัวเดียว)", fs.TotalExpenses)
	}
}

// TestStatement_RealProfessionPaycheck — ใช้อาชีพจริง (วิศวกร) ตรวจ breakdown เงินเดือน
// เป็น teaching artifact: เห็นเลขจริงว่าเงินเดือนถูกหัก ณ ที่จ่ายยังไง
func TestStatement_RealProfessionPaycheck(t *testing.T) {
	var eng domain.Profession
	for _, p := range profession.All() {
		if p.Name == "วิศวกร" {
			eng = p
			break
		}
	}
	if eng.Name == "" {
		t.Fatal("ไม่พบ 'วิศวกร' ใน profession dataset")
	}

	player := domain.Player{ID: "p1", Cash: eng.Savings, Profession: eng}
	fs := Statement(player)

	// breakdown ตรงตาม formula
	if fs.Tax != payroll.MonthlyTax(eng.Salary) {
		t.Errorf("Tax = %d, want %d", fs.Tax, payroll.MonthlyTax(eng.Salary))
	}
	if fs.SocialSecurity != payroll.SocialSecurity(eng.Salary) {
		t.Errorf("SocialSecurity = %d, want %d", fs.SocialSecurity, payroll.SocialSecurity(eng.Salary))
	}

	// MonthlyCashFlow = เงินเดือน − ภาษี − SS − other − ผ่อนหนี้ทั้งหมด
	payments := eng.HomeMortgage.Payment + eng.CarLoan.Payment + eng.CreditCard.Payment + eng.SchoolLoan.Payment
	want := eng.Salary - fs.Tax - fs.SocialSecurity - eng.OtherExpenses - payments
	if fs.MonthlyCashFlow != want {
		t.Errorf("MonthlyCashFlow = %d, want %d", fs.MonthlyCashFlow, want)
	}

	// แสดง paycheck จริงใน test output (สอนผู้สร้าง)
	t.Logf("วิศวกร: เงินเดือน %d − ภาษี %d − SS %d − อื่นๆ %d − ผ่อน %d = สุทธิ %d/เดือน",
		eng.Salary, fs.Tax, fs.SocialSecurity, eng.OtherExpenses, payments, fs.MonthlyCashFlow)
}
