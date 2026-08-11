package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/finance-boardgame/backend/internal/config"
)

// TestHealthz_ReturnsExpectedContract เป็นตัวอย่าง "contract test" สำหรับ API ที่ frontend พึ่งพา
//
// ทำไมต้องมี: ฝั่ง web (lib/api/client.ts) คาดหวังว่า /healthz จะคืน { status, service, engine }
// ถ้าวันนึงมีคนเปลี่ยนชื่อ field "engine" → "version" หรือเปลี่ยน shape เทสนี้จะแดงทันที
// → กัน frontend พังแบบเงียบ และกันมานั่งไล่แก้ "ผลผิดไม่จบไม่สิ้น"
//
// หลักการ: assert "contract/output ที่ frontend ต้องการ" ไม่ใช่ internals ของ handler
func TestHealthz_ReturnsExpectedContract(t *testing.T) {
	// Arrange — สร้าง server จริง (wired ครบ) แต่ไม่ต้อง bind port
	srv := NewServer(config.Config{HTTPAddr: ":0"})
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	// Act
	srv.mux.ServeHTTP(rec, req)

	// Assert — contract ที่ frontend คาดหวัง
	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	// ทุก field ที่ frontend ใช้ต้องอยู่ครบ — ถ้าหายไป = contract broken
	requiredKeys := []string{"status", "service", "engine"}
	for _, key := range requiredKeys {
		if _, ok := body[key]; !ok {
			t.Errorf("response missing required field %q — frontend contract broken", key)
		}
	}

	// ค่า business ที่ frontend ตั้งใจตรวจ
	if body["status"] != "ok" {
		t.Errorf("status = %q, want %q", body["status"], "ok")
	}
}
