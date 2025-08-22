package main

import (
	"log"
	"os"

	"github.com/rqull/GIN_Example/internal/config"
	"github.com/rqull/GIN_Example/internal/db"
	"github.com/rqull/GIN_Example/internal/router"
)

func httpPort() string {
	if p := os.Getenv("PORT"); p != "" { // prioritas env dari Railway
		return p
	}
	if p := os.Getenv("APP_PORT"); p != "" {
		return p
	}
	return "8080"
}

func main() {
	config.Load()

	// Koneksi database
	database := db.ConnectDB()
	defer database.Close()

	// Jalankan migrasi database
	if err := db.RunMigrations(database, "migrations"); err != nil {
		log.Fatal("Gagal menjalankan migrasi:", err)
	}

	// Setup router
	r := router.SetupRouter(database)

	// Jalankan server dengan port dari environment
	port := httpPort()
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
