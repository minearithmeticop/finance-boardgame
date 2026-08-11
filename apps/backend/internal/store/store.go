// Package store จัดการ persistence ของ backend
//   - PostgreSQL: เก็บ game history / ผู้เล่น / replay
//   - Redis: เก็บ active game state (fast read/write ระหว่างเล่น)
//
// TODO(Session#3): define interfaces + concrete repos (สำหรับ mock ใน test)
package store

// Store aggregates repositories ทั้งหมด
type Store struct {
	// TODO: Games GameRepo, Players PlayerRepo, ActiveGames ActiveGameRepo
}
