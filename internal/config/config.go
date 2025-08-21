package config

import (
	"os"

	"github.com/joho/godotenv"
)

func Load() {
	// Abaikan error jika file .env tidak ditemukan, gunakan nilai dari environment
	_ = godotenv.Load()
}

func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
