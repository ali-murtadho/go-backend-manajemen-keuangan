package controllers

import (
	"backend-api/database"
	"backend-api/helpers"
	"backend-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Get All Waktu
// @Description Mendapatkan semua data waktu
// @Tags Waktu
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Router /times [get]
// @Param Authorization header string true "Authorization token" example("Bearer <token>")
func GetListWaktu(c *gin.Context) {

	// Inisialisasi slice untuk menampung data waktu
	var times []models.Waktu

	// Ambil data waktu dari database
	database.DB.Find(&times)

	// Kirimkan response sukses dengan data waktu
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Lists Data Waktu",
		Data:    times,
	})
}

// @Summary Get Waktu By Bulan Tahun
// @Description Mendapatkan data waktu berdasarkan bulan dan tahun
// @Tags Waktu
// @Accept json
// @Produce json
// @Param bulan_tahun path string true "Bulan dan Tahun"
// @Success 200 {object} models.SuccessResponse
// @Success 404 {object} models.ErrorResponse
// @Router /times/{bulan_tahun} [get]
// @Param Authorization header string true "Authorization token" example("Bearer <token>")
func GetWaktuByBulanTahun(c *gin.Context) {
	bulanTahun := c.Param("bulan_tahun")
	var waktu models.Waktu

	// Cari data waktu berdasarkan bulan dan tahun
	err := database.DB.Where("bulan_tahun = ?", bulanTahun).First(&waktu).Error
	if err != nil {
		// Jika data tidak ditemukan, kirimkan response error
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success: false,
				Message: "Data tidak ditemukan",
			})
			return
		}
		// Jika terjadi error lainnya, kirimkan response error
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Kirimkan response sukses dengan data waktu
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Lists Data Waktu",
		Data:    waktu,
	})
}

// @Summary Create Waktu
// @Description Membuat data waktu baru
// @Tags Waktu
// @Accept json
// @Produce json
// @Param request body models.WaktuRequest true "Request body"
// @Success 201 {object} models.SuccessResponse
// @Success 422 {object} models.ErrorResponse
// @Router /times [post]
// @Param Authorization header string true "Authorization token" example("Bearer <token>")
func CreateTimes(c *gin.Context) {

	//struct user request
	var req = models.WaktuRequest{}

	// Bind JSON request ke struct UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Inisialisasi user baru
	waktu := models.Waktu{
		BulanTahun: req.BulanTahun,
	}

	// Simpan user ke database
	database.DB.Create(&waktu)

	// Kirimkan response sukses
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Success Create Data Waktu",
		Data:    waktu,
	})
}

// @Summary Update Waktu By ID
// @Description Memperbarui data waktu berdasarkan ID
// @Tags Waktu
// @Accept json
// @Produce json
// @Param id path string true "ID Waktu"
// @Param request body models.WaktuRequest true "Request body"
// @Success 200 {object} models.SuccessResponse
// @Success 422 {object} models.ErrorResponse
// @Success 404 {object} models.ErrorResponse
// @Router /times/{id} [put]
// @Param Authorization header string true "Authorization token" example("Bearer <token>")
func UpdateTimes(c *gin.Context) {

	// Ambil ID user dari parameter URL
	id := c.Param("id")

	// Inisialisasi user
	var waktu models.Waktu

	// Cari user berdasarkan ID
	if err := database.DB.First(&waktu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Bind JSON request ke struct UserRequest
	if err := c.ShouldBindJSON(&waktu); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Simpan user ke database
	database.DB.Save(&waktu)

	// Kirimkan response sukses
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Update Data Waktu",
		Data:    waktu,
	})
}

// @Summary Delete Waktu By ID
// @Description Menghapus data waktu berdasarkan ID
// @Tags Waktu
// @Accept json
// @Produce json
// @Param id path string true "ID Waktu"
// @Success 200 {object} models.SuccessResponse
// @Success 404 {object} models.ErrorResponse
// @Router /times/{id} [delete]
// @Param Authorization header string true "Authorization token" example("Bearer <token>")
func DeleteTimes(c *gin.Context) {
	id := c.Param("id")
	var waktu models.Waktu
	err := database.DB.First(&waktu, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success: false,
				Message: "Data tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	database.DB.Delete(&waktu)
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Delete Data Waktu",
		Data:    waktu,
	})
}
