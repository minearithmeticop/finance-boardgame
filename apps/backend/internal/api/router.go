// Package api คือ HTTP/WebSocket layer ของ backend
// รับ request จาก Next.js แล้วส่งต่อเข้า engine/room/store
package api

import (
	"encoding/json"
	"net/http"

	"github.com/finance-boardgame/backend/internal/config"
	"github.com/finance-boardgame/backend/internal/room"
	"github.com/finance-boardgame/engine"
)

// Server คือ HTTP server หลัก
type Server struct {
	cfg  config.Config
	hub  *room.Hub
	mux  *http.ServeMux
}

// NewServer สร้าง server พร้อม wire dependencies
func NewServer(cfg config.Config) *Server {
	s := &Server{
		cfg: cfg,
		hub: room.NewHub(),
		mux: http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealth)
	// TODO(Session#3):
	//   POST   /api/v1/games          — สร้างเกม/ห้อง
	//   GET    /api/v1/games/{id}     — ดึงสถานะเกม
	//   POST   /api/v1/games/{id}/action — ส่ง action (async mode)
	//   GET    /ws/{roomId}           — WebSocket (real-time mode)
}

// ListenAndServe เริ่มรับ HTTP connection
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.cfg.HTTPAddr, s.mux)
}

// handleHealth — endpoint ตรวจสอบว่า service พร้อมใช้งาน
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "finance-boardgame-backend",
		"engine":  engine.Version,
	})
}
