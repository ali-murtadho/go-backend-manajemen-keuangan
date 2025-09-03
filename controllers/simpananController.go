package controllers

import (
	"backend-api/database"
	"backend-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Get All Simpanan
// @Description Mendapatkan semua data simpanan
// @Tags Simpanan
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Router /simpanan [get]
// @Param Authorization header string true "Authorization token" example("Bearer <token>")
func GetListSimpanan(c *gin.Context) {

	// Inisialisasi slice untuk menampung data simpanan
	var simpanan []models.Simpanan

	// Ambil data simpanan dari database
	database.DB.Find(&simpanan)

	// Kirimkan response sukses dengan data simpanan
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Lists Data Simpanan",
		Data:    simpanan,
	})
}

// @Summary Get Simpanan By ID
// @Description Mendapatkan data simpanan berdasarkan ID
// @Tags Simpanan
// @Accept json
// @Produce json
// @Param id path string true "ID Simpanan"
// @Success 200 {object} models.SuccessResponse
// @Success 404 {object} models.ErrorResponse
// @Router /simpanan/{id} [get]
// @Param Authorization header string true "Authorization token" example("Bearer <token>")
func GetSimpananById(c *gin.Context) {
	id := c.Param("id")
	var simpanan models.Simpanan
	if err := database.DB.Where("id = ?", id).First(&simpanan).Error; err != nil {
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
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Data Simpanan",
		Data:    simpanan,
	})
}

// @Summary Create Simpanan
// @Description Membuat data simpanan baru
// @Tags Simpanan
// @Accept json
// @Produce json
// @Param simpanan body models.Simpanan true "Data Simpanan"
// @Success 200 {object} models.SuccessResponse
// @Success 400 {object} models.ErrorResponse
// @Router /simpanan [post]
// @Param Authorization header string true "Authorization token" example("Bearer <token>")
func CreateSimpanan(c *gin.Context) {
	var simpanan models.Simpanan
	if err := c.ShouldBindJSON(&simpanan); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Validasi tidak boleh ada data yang kosong
	if simpanan.Nama == "" || simpanan.JenisSimpanan == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Nama dan Jenis Simpanan tidak boleh kosong",
		})
		return
	}

	database.DB.Create(&simpanan)
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Success",
		Data:    simpanan,
	})
}

// @Summary Update Simpanan
// @Description Memperbarui data simpanan
// @Tags Simpanan
// @Accept json
// @Produce json
// @Param id path string true "ID Simpanan"
// @Param simpanan body models.Simpanan true "Data Simpanan"
// @Success 200 {object} models.SuccessResponse
// @Success 400 {object} models.ErrorResponse
// @Router /simpanan/{id} [put]
// @Param Authorization header string true "Authorization token" example("Bearer <token>")
func UpdateSimpanan(c *gin.Context) {
	id := c.Param("id")
	var simpanan models.Simpanan
	if err := c.ShouldBindJSON(&simpanan); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	database.DB.Model(&simpanan).Where("id = ?", id).Updates(&simpanan)
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Data Simpanan",
		Data:    simpanan,
	})
}

// @Summary Delete Simpanan By ID
// @Description Menghapus data simpanan berdasarkan ID
// @Tags Simpanan
// @Accept json
// @Produce json
// @Param id path string true "ID Simpanan"
// @Success 200 {object} models.SuccessResponse
// @Success 404 {object} models.ErrorResponse
// @Router /simpanan/{id} [delete]
// @Param Authorization header string true "Authorization token" example("Bearer <token>")
func DeleteSimpanan(c *gin.Context) {
	id := c.Param("id")
	var simpanan models.Simpanan
	if err := database.DB.Where("id = ?", id).First(&simpanan).Error; err != nil {
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
	database.DB.Delete(&simpanan)
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Data Simpanan",
		Data:    simpanan,
	})
}
