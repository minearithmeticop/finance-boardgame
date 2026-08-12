package payroll

import (
	"testing"

	"github.com/finance-boardgame/engine/domain"
)

// TestSocialSecurity_CappedAt750 — 5% ของเงินเดือน เพดานฐาน 15,000 → สูงสุด 750/เดือน
func TestSocialSecurity_CappedAt750(t *testing.T) {
	cases := []struct {
		salary domain.Money
		want   domain.Money
	}{
		{0, 0},
		{10_000, 500},  // 5% ของ 10,000
		{15_000, 750},  // ขอบเพดานฐาน → 750
		{20_000, 750},  // เกินฐาน → cap 750
		{45_000, 750},  // วิศวกร
		{120_000, 750}, // แพทย์
	}
	for _, c := range cases {
		if got := SocialSecurity(c.salary); got != c.want {
			t.Errorf("SocialSecurity(%d) = %d, want %d", c.salary, got, c.want)
		}
	}
}

// TestAnnualTax_ProgressiveBrackets — progressive ตามขั้นบันได
func TestAnnualTax_ProgressiveBrackets(t *testing.T) {
	cases := []struct {
		name    string
		taxable domain.Money
		want    domain.Money
	}{
		{"zero", 0, 0},
		{"within 0% band", 71_000, 0},          // < 150k → 0 (พนักงานออฟฟิศ 20k)
		{"top of 0% band", 150_000, 0},          // ขอบบน 0%
		{"into 5% band", 300_000, 7_500},        // 150k @5%
		{"วิศวกร 371k", 371_000, 14_600},        // 150k@5 + 71k@10
		{"500k", 500_000, 27_500},               // +200k@10 = 7500+20000
		{"1M", 1_000_000, 115_000},              // +250k@15 +250k@20
	}
	for _, c := range cases {
		if got := AnnualTax(c.taxable); got != c.want {
			t.Errorf("%s: AnnualTax(%d) = %d, want %d", c.name, c.taxable, got, c.want)
		}
	}
}

// TestTaxableIncome_พนักงานออฟฟิศ — salary 20k → taxable 71k (หัก ลดหย่อน 60k + ค่าใช้จ่าย 100k + SS 9k)
func TestTaxableIncome_OfficeWorker20k(t *testing.T) {
	// annual 240k − 60k − 100k − 9k = 71k
	if got := TaxableIncome(20_000); got != 71_000 {
		t.Errorf("TaxableIncome(20000) = %d, want 71000", got)
	}
}

// TestTaxableIncome_วิศวกร45k — salary 45k → taxable 371k
func TestTaxableIncome_Engineer45k(t *testing.T) {
	// annual 540k − 60k − 100k − 9k = 371k
	if got := TaxableIncome(45_000); got != 371_000 {
		t.Errorf("TaxableIncome(45000) = %d, want 371000", got)
	}
}

// TestMonthlyTax_วิศวกร — salary 45k → ภาษี ~1,216/เดือน (14,600/12)
func TestMonthlyTax_Engineer45k(t *testing.T) {
	// 14,600 / 12 = 1,216 (integer truncation)
	if got := MonthlyTax(45_000); got != 1_216 {
		t.Errorf("MonthlyTax(45000) = %d, want 1216", got)
	}
}

// TestMonthlyTax_LowIncomeZero — salary 20k → ภาษี 0 (รายได้ต่ำ ไม่เสียภาษี)
func TestMonthlyTax_LowIncomeZero(t *testing.T) {
	if got := MonthlyTax(20_000); got != 0 {
		t.Errorf("MonthlyTax(20000) = %d, want 0 (low income exempt)", got)
	}
}
