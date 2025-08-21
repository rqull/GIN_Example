package controller

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rqull/GIN_Example/internal/models"
)

type BioskopController struct {
	DB *sql.DB
}

func NewBioskopController(db *sql.DB) *BioskopController {
	return &BioskopController{DB: db}
}

func respondError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

func parseIDParam(c *gin.Context) (int, bool) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "ID tidak valid")
		return 0, false
	}
	return id, true
}

func validateBioskopInput(input models.BioskopInput) (models.BioskopInput, string) {
	input.Nama = strings.TrimSpace(input.Nama)
	input.Lokasi = strings.TrimSpace(input.Lokasi)

	if input.Nama == "" {
		return input, "Nama tidak boleh kosong"
	}
	if input.Lokasi == "" {
		return input, "Lokasi tidak boleh kosong"
	}
	if input.Rating < 0 || input.Rating > 5 {
		return input, "Rating harus antara 0 sampai 5"
	}
	return input, ""
}

func (ctl *BioskopController) Create(c *gin.Context) {
	var input models.BioskopInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	input, errMsg := validateBioskopInput(input)
	if errMsg != "" {
		respondError(c, http.StatusBadRequest, errMsg)
		return
	}

	var created models.Bioskop
	query := `INSERT INTO bioskop (nama, lokasi, rating) VALUES ($1, $2, $3) RETURNING id, nama, lokasi, rating`
	if err := ctl.DB.QueryRow(query, input.Nama, input.Lokasi, input.Rating).
		Scan(&created.ID, &created.Nama, &created.Lokasi, &created.Rating); err != nil {
		respondError(c, http.StatusInternalServerError, "Gagal menyimpan data")
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (ctl *BioskopController) GetAll(c *gin.Context) {
	rows, err := ctl.DB.Query("SELECT id, nama, lokasi, rating FROM bioskop ORDER BY id")
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Gagal mengambil data")
		return
	}
	defer rows.Close()

	var result []models.Bioskop
	for rows.Next() {
		var b models.Bioskop
		if err := rows.Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating); err != nil {
			respondError(c, http.StatusInternalServerError, "Gagal membaca data")
			return
		}
		result = append(result, b)
	}
	if err := rows.Err(); err != nil {
		respondError(c, http.StatusInternalServerError, "Gagal membaca data")
		return
	}

	c.JSON(http.StatusOK, result)
}

func (ctl *BioskopController) GetByID(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var b models.Bioskop
	err := ctl.DB.QueryRow("SELECT id, nama, lokasi, rating FROM bioskop WHERE id = $1", id).
		Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(c, http.StatusNotFound, "Bioskop tidak ditemukan")
			return
		}
		respondError(c, http.StatusInternalServerError, "Gagal mengambil data")
		return
	}

	c.JSON(http.StatusOK, b)
}

func (ctl *BioskopController) Update(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var input models.BioskopInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	input, errMsg := validateBioskopInput(input)
	if errMsg != "" {
		respondError(c, http.StatusBadRequest, errMsg)
		return
	}

	var updated models.Bioskop
	query := `UPDATE bioskop SET nama = $1, lokasi = $2, rating = $3 WHERE id = $4 RETURNING id, nama, lokasi, rating`
	err := ctl.DB.QueryRow(query, input.Nama, input.Lokasi, input.Rating, id).
		Scan(&updated.ID, &updated.Nama, &updated.Lokasi, &updated.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(c, http.StatusNotFound, "Bioskop tidak ditemukan")
			return
		}
		respondError(c, http.StatusInternalServerError, "Gagal memperbarui data")
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (ctl *BioskopController) Delete(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	res, err := ctl.DB.Exec("DELETE FROM bioskop WHERE id = $1", id)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Gagal menghapus data")
		return
	}
	affected, err := res.RowsAffected()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Gagal menghapus data")
		return
	}
	if affected == 0 {
		respondError(c, http.StatusNotFound, "Bioskop tidak ditemukan")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bioskop berhasil dihapus"})
}
