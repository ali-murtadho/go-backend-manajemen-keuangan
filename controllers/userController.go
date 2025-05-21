package controllers

import (
	"backend-api/database"
	"backend-api/helpers"
	"backend-api/models"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register menangani proses registrasi user baru
func Register(c *gin.Context) {

	// Inisialisasi struct untuk menangkap data request
	var req = models.UserCreateRequest{}

	// Validasi request JSON menggunakan binding dari Gin
	if err := c.ShouldBindJSON(&req); err != nil {
		// Jika validasi gagal, kirimkan response error
		c.JSON(http.StatusUnprocessableEntity, models.ErrorResponse{
			Success: false,
			Message: "Validasi Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	pwd := req.Password
	// quick manual checks
	if len(pwd) < 8 ||
		!regexp.MustCompile(`[A-Z]`).MatchString(pwd) ||
		!regexp.MustCompile(`[a-z]`).MatchString(pwd) ||
		!regexp.MustCompile(`\d`).MatchString(pwd) ||
		!regexp.MustCompile(`[!@#~\$%\^&\*\(\)]`).MatchString(pwd) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Password must be at least 8 characters, include a symbol, a number, uppercase and lowercase letters",
		})
		return
	}

	// Buat data user baru dengan password yang sudah di-hash
	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(pwd),
	}

	// Simpan data user ke database
	if err := database.DB.Create(&user).Error; err != nil {
		// Cek apakah error karena data duplikat (misalnya username/email sudah terdaftar)
		if helpers.IsDuplicateEntryError(err) {
			// Jika duplikat, kirimkan response 409 Conflict
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Success: false,
				Message: "Duplicate entry error",
				Errors:  helpers.TranslateErrorMessage(err),
			})
		} else {
			// Jika error lain, kirimkan response 500 Internal Server Error
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Message: "Failed to create user",
				Errors:  helpers.TranslateErrorMessage(err),
			})
		}
		return
	}

	// Jika berhasil, kirimkan response sukses
	c.JSON(http.StatusCreated, models.SuccessResponse{
		Success: true,
		Message: "User created successfully",
		Data: models.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			// CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			// UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func Login(c *gin.Context) {

	// Inisialisasi struct untuk menampung data dari request
	var req = models.UserLoginRequest{}
	var user = models.User{}

	// Validasi input dari request body menggunakan ShouldBindJSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Cari user berdasarkan username yang diberikan di database
	// Jika tidak ditemukan, kirimkan respons error Unauthorized
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "User Not Found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Bandingkan password yang dimasukkan dengan password yang sudah di-hash di database
	// Jika tidak cocok, kirimkan respons error Unauthorized
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Message: "Invalid Password",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Jika login berhasil, generate token untuk user
	token := helpers.GenerateToken(user.Username)

	// Kirimkan response sukses dengan status OK dan data user serta token
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Login Success",
		Data: models.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Token:     &token,
		},
	})
}

func FindUsers(c *gin.Context) {

	// Inisialisasi slice untuk menampung data user
	var users []models.User

	// Ambil data user dari database
	database.DB.Find(&users)

	// Kirimkan response sukses dengan data user
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Lists Data Users",
		Data:    users,
	})
}
