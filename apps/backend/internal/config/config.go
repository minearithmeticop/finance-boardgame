// Package config โหลดการตั้งค่า backend จาก environment
package config

import "os"

// Config เก็บการตั้งค่า runtime ของ backend
type Config struct {
	HTTPAddr   string // address ที่ HTTP server รับ เช่น ":8080"
	PostgresDSN string // TODO(Session#3): เชื่อม PostgreSQL
	RedisAddr  string // TODO(Session#3): เชื่อม Redis
}

// Load อ่าน config จาก env (มี default)
func Load() Config {
	return Config{
		HTTPAddr:    getenv("HTTP_ADDR", ":8080"),
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
		RedisAddr:   getenv("REDIS_ADDR", "localhost:6379"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
