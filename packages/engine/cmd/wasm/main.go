//go:build js && wasm

// Package main คือ WebAssembly entrypoint ของ engine
//
// คอมไพล์ด้วย: GOOS=js GOARCH=wasm go build -o engine.wasm ./cmd/wasm
// (หรือ `pnpm build:wasm` จาก root)
//
// expose engine API ให้ JavaScript เรียกผ่าน globalThis.engine*
// ทุกฟังก์ชันคืน JSON string รูปแบบ { "data": <any>, "error": "<string>" }
// ฝั่ง TS แกะ envelope แล้ว throw ถ้ามี error
package main

import (
	"encoding/json"
	"errors"
	"sync"
	"syscall/js"

	"github.com/finance-boardgame/engine"
	"github.com/finance-boardgame/engine/domain"
	"github.com/finance-boardgame/engine/ratrace"
)

// single global engine instance — เพียงพอสำหรับ local/pass-and-play demo
// (multi-game registry ไว้ทำทีหลัง)
var (
	mu      sync.Mutex
	current *engine.Engine
)

func main() {
	g := js.Global()
	g.Set("engineVersion", js.FuncOf(func(this js.Value, args []js.Value) any {
		return engine.Version // คืน string plain (back-compat กับ loader เดิม)
	}))
	g.Set("engineCreate", js.FuncOf(engineCreate))
	g.Set("engineState", js.FuncOf(engineState))
	g.Set("engineApply", js.FuncOf(engineApply))
	g.Set("engineBoard", js.FuncOf(engineBoard))

	// ส่งสัญญาณ "พร้อม" — loader poll จนกว่าจะเจอ flag นี้ ก่อนใช้งาน
	// เพื่อกัน race ตอน startup (callbacks ยังไม่ถูก register)
	g.Set("__engineWasmReady", true)

	// WASM main ต้อง block ตลอดไป มิฉะนั้นโปรแกรมจะจบทันที
	select {}
}

// engineCreate(seed, playersJSON) → สร้างเกมใหม่
func engineCreate(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return envelope(nil, errors.New("engineCreate needs (seed, playersJSON)"))
	}
	seed := int64(args[0].Int())
	var players []domain.Player
	if err := json.Unmarshal([]byte(args[1].String()), &players); err != nil {
		return envelope(nil, err)
	}
	mu.Lock()
	current = engine.New(engine.Config{Seed: seed, Players: players})
	mu.Unlock()
	return envelope(nil, nil)
}

// engineState() → snapshot สถานะเกมปัจจุบัน
func engineState(this js.Value, args []js.Value) any {
	mu.Lock()
	defer mu.Unlock()
	if current == nil {
		return envelope(nil, errors.New("no game created"))
	}
	return envelope(current.State(), nil)
}

// engineApply(actionJSON) → ประมวลผล action คืน events
func engineApply(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return envelope(nil, errors.New("engineApply needs (actionJSON)"))
	}
	var action domain.Action
	if err := json.Unmarshal([]byte(args[0].String()), &action); err != nil {
		return envelope(nil, err)
	}
	mu.Lock()
	defer mu.Unlock()
	if current == nil {
		return envelope(nil, errors.New("no game created"))
	}
	events, err := current.Apply(action)
	return envelope(events, err)
}

// engineBoard() → layout กระดาน Rat Race (จาก ratrace.DefaultBoard — single source)
func engineBoard(this js.Value, args []js.Value) any {
	return envelope(ratrace.DefaultBoard(), nil)
}

// envelope ห่อ data+error เป็น JSON string ตาม protocol ของ boundary
func envelope(data any, err error) string {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	b, _ := json.Marshal(struct {
		Data  any    `json:"data"`
		Error string `json:"error"`
	}{Data: data, Error: errMsg})
	return string(b)
}
