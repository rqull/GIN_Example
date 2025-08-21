package main

import (
	"log"

	"github.com/rqull/GIN_Example/internal/config"
	"github.com/rqull/GIN_Example/internal/db"
	"github.com/rqull/GIN_Example/internal/router"
)

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
	port := config.GetEnv("APP_PORT", "8080")
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
