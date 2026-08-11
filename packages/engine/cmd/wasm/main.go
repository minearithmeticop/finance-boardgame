//go:build js && wasm

// Package main คือ WebAssembly entrypoint ของ engine
//
// คอมไพล์ด้วย:
//   GOOS=js GOARCH=wasm go build -o engine.wasm ./cmd/wasm
//
// รันใน browser ผ่าน wasm_exec.js (ดู apps/web/lib/engine-wasm)
package main

import (
	"syscall/js"

	"github.com/finance-boardgame/engine"
)

func main() {
	// ลงทะเบียนฟังก์ชัน engine ให้เรียกจาก JavaScript ได้ผ่าน globalThis.engine.*
	g := js.Global()
	g.Set("engineVersion", js.FuncOf(func(this js.Value, args []js.Value) any {
		return engine.Version
	}))

	// TODO(Session#3): expose engine.New(), Apply(), State() ผ่าน JS bindings
	//   - serialize GameState ↔ JSON เพื่อข้าม boundary Go↔JS
	//   - สร้าง wrapper object ที่ถือ engine instance ต่อหนึ่งเกม

	// WASM main ต้อง block ตลอดไป มิฉะนั้นโปรแกรมจะจบทันที
	select {}
}
