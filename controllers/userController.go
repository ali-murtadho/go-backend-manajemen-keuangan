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

// @Summary Register User
// @Description Mendaftarkan user baru
// @Tags User - Auth
// @Accept json
// @Produce json
// @Param request body models.UserCreateRequest true "Request body"
// @Success 200 {object} models.SuccessResponse
// @Router /register [post]
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

// @Summary Login User
// @Description Melakukan login user
// @Tags User - Auth
// @Accept json
// @Produce json
// @Param request body models.UserLoginRequest true "Request body"
// @Success 200 {object} models.SuccessResponse
// @Success 401 {object} models.ErrorResponse
// @Router /login [post]
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

// @Summary Get All Users
// @Description Mendapatkan semua data user
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Router /users [get]
// @Param Authorization header string true "Authorization token" example("Bearer <token>")
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

// @Summary Create User
// @Description Membuat data user baru
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.UserCreateRequest true "Request body"
// @Success 201 {object} models.SuccessResponse
// @Success 422 {object} models.ErrorResponse
// @Router /users [post]
// @Param Authorization header string true "Authorization token" example("Bearer <token>")
func CreateUser(c *gin.Context) {

	//struct user request
	var req = models.UserCreateRequest{}

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
	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(req.Password),
	}

	// Simpan user ke database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirimkan response sukses
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
		},
	})

}

// @Summary Get User By ID
// @Description Mendapatkan data user berdasarkan ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "ID User"
// @Success 200 {object} models.SuccessResponse
// @Success 404 {object} models.ErrorResponse
// @Router /users/{id} [get]
// @Param Authorization header string true "Authorization token" example("Bearer <token>")
func FindUserById(c *gin.Context) {

	// Ambil ID user dari parameter URL
	id := c.Param("id")

	// Inisialisasi user
	var user models.User

	// Cari user berdasarkan ID
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirimkan response sukses dengan data user
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "User Found",
		Data: models.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

// @Summary Update User By ID
// @Description Memperbarui data user berdasarkan ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "ID User"
// @Param request body models.UserUpdateRequest true "Request body"
// @Success 200 {object} models.SuccessResponse
// @Success 422 {object} models.ErrorResponse
// @Success 404 {object} models.ErrorResponse
// @Router /users/{id} [put]
// @Param Authorization header string true "Authorization token" example("Bearer <token>")
func UpdateUser(c *gin.Context) {

	// Ambil ID user dari parameter URL
	id := c.Param("id")

	// Inisialisasi user
	var user models.User

	// Cari user berdasarkan ID
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	//struct user request
	var req = models.UserUpdateRequest{}

	// Bind JSON request ke struct UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Update user dengan data baru
	user.Name = req.Name
	user.Username = req.Username
	user.Email = req.Email
	user.Password = helpers.HashPassword(req.Password)

	// Simpan perubahan ke database
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to update user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirimkan response sukses
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "User updated successfully",
		Data: models.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

// @Summary Delete User By ID
// @Description Menghapus data user berdasarkan ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "ID User"
// @Success 200 {object} models.SuccessResponse
// @Success 404 {object} models.ErrorResponse
// @Router /users/{id} [delete]
// @Param Authorization header string true "Authorization token" example("Bearer <token>")
func DeleteUser(c *gin.Context) {

	// Ambil ID user dari parameter URL
	id := c.Param("id")

	// Inisialisasi user
	var user models.User

	// Cari user berdasarkan ID
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Hapus user dari database
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to delete user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirimkan response sukses
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}
