package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Bioskop struct {
	ID     int     `json:"id"`
	Nama   string  `json:"nama"`
	Lokasi string  `json:"lokasi"`
	Rating float64 `json:"rating"`
}

var db *sql.DB

func ConnectDB() *sql.DB {
	dsn := "host=localhost port=5432 user=postgres password=12345 dbname=bioskopdb sslmode=disable"

	database, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal membuka koneksi database:", err)
	}

	if err := database.Ping(); err != nil {
		log.Fatal("Tidak bisa ping database:", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS bioskop (
		id SERIAL PRIMARY KEY,
		nama VARCHAR(100) NOT NULL,
		lokasi VARCHAR(100) NOT NULL,
		rating REAL
	);
	`
	if _, err := database.Exec(createTableQuery); err != nil {
		log.Fatal("Gagal membuat tabel bioskop:", err)
	}

	fmt.Println("Database connected & tabel ready")
	return database
}

func CreateBioskop(c *gin.Context) {
	var input Bioskop

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
		return
	}

	if input.Nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama tidak boleh kosong"})
		return
	}

	if input.Lokasi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lokasi tidak boleh kosong"})
	}

	query := `INSERT INTO bioskop (nama, lokasi, rating) VALUES ($1, $2, $3) RETURNING id`
	if err := db.QueryRow(query, input.Nama, input.Lokasi, input.Rating).Scan(&input.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func GetBioskop(c *gin.Context) {
	rows, err := db.Query("SELECT id, nama, lokasi, rating FROM bioskop")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	var result []Bioskop
	for rows.Next() {
		var b Bioskop
		if err := rows.Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data"})
			return
		}
		result = append(result, b)
	}

	c.JSON(http.StatusOK, result)
}

func main() {
	db = ConnectDB()
	defer db.Close()

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.POST("/bioskop", CreateBioskop)
		api.GET("/bioskop", GetBioskop)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
