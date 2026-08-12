package profession

import (
	"testing"

	"github.com/finance-boardgame/engine/payroll"
	"github.com/finance-boardgame/engine/rng"
)

// TestAll_NotEmptyAndValid — ทุกอาชีพต้องมีชื่อ/เงินเดือน และ tax/SS ตรง formula (กัน drift)
func TestAll_NotEmptyAndValid(t *testing.T) {
	all := All()
	if len(all) == 0 {
		t.Fatal("All() returned no professions")
	}
	for _, p := range all {
		if p.Name == "" {
			t.Error("found profession with empty name")
		}
		if p.Salary <= 0 {
			t.Errorf("%s: salary = %d, want > 0", p.Name, p.Salary)
		}
		// tax/SS ต้องตรงกับ payroll calc (data ต้องไม่ drift จาก formula)
		if p.Taxes != payroll.MonthlyTax(p.Salary) {
			t.Errorf("%s: Taxes = %d, want %d (payroll.MonthlyTax)", p.Name, p.Taxes, payroll.MonthlyTax(p.Salary))
		}
		if p.SocialSecurity != payroll.SocialSecurity(p.Salary) {
			t.Errorf("%s: SocialSecurity = %d, want %d (payroll.SocialSecurity)", p.Name, p.SocialSecurity, payroll.SocialSecurity(p.Salary))
		}
	}
}

// TestAll_ReturnsCopy — แก้ผลลัพธ์ต้องไม่กระทบชุดต้นฉบับ
func TestAll_ReturnsCopy(t *testing.T) {
	a := All()
	a[0].Name = "ถูกแก้"
	if allProfessions[0].Name == "ถูกแก้" {
		t.Error("All() returned the underlying slice (mutation leaked)")
	}
}

// TestRandom_ReturnsMemberOfAll — Random ต้องคืนอาชีพที่อยู่ในชุด และหลากหลาย
func TestRandom_ReturnsMemberOfAll(t *testing.T) {
	all := All()
	names := map[string]bool{}
	for _, p := range all {
		names[p.Name] = true
	}

	r := rng.New(7)
	seen := map[string]bool{}
	for i := 0; i < 300; i++ {
		p := Random(r)
		if !names[p.Name] {
			t.Errorf("Random returned unknown profession %q", p.Name)
		}
		seen[p.Name] = true
	}
	if len(seen) < 2 {
		t.Errorf("Random looks biased: only saw %d distinct professions", len(seen))
	}
}
