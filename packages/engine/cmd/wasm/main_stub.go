//go:build !(js && wasm)

// Package main คือ placeholder สำหรับ platform อื่น (non-WASM)
// ไฟล์จริงที่ใช้รันใน browser อยู่ใน main.go (build tag: js && wasm)
// ไฟล์นี้ทำให้ `go build ./...` ไม่ error เรื่อง "no Go files" บน platform ปกติ
package main

func main() {}
