package test

import (
	"backend-api/config"
	"backend-api/database"
	"backend-api/models"
	"backend-api/routes"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func initTestDB() {
	// pastikan env test dipakai
	os.Setenv("ENV", "test")

	// load env (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, dst.)
	config.LoadEnv()

	// konek ke PostgreSQL sesuai env
	database.InitDB()

	// migrate tabel yang dipakai
	database.DB.AutoMigrate(&models.User{}, &models.Waktu{})
}

func setupRouter() *gin.Engine {
	initTestDB()
	return routes.SetupRouter()
}

func registerAndLogin(router *gin.Engine, t *testing.T) string {
	// --- Register user ---
	registerBody := `{"name":"Tester","username":"tester","email":"tester@mail.com","password":"Password!1"}`
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBufferString(registerBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// --- Login user ---
	loginBody := `{"username":"tester","password":"Password!1"}`
	req2 := httptest.NewRequest("POST", "/api/login", bytes.NewBufferString(loginBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	var loginResp models.SuccessResponse
	json.Unmarshal(w2.Body.Bytes(), &loginResp)
	assert.True(t, loginResp.Success)

	userData := loginResp.Data.(map[string]interface{})

	// cek dulu isi field apa yang ada
	var token string
	if tkn, ok := userData["token"].(string); ok {
		token = tkn
	} else if tkn, ok := userData["Token"].(string); ok {
		token = tkn
	} else {
		t.Fatalf("token not found in login response: %+v", userData)
	}

	return token
}

func TestTimesCRUD(t *testing.T) {
	router := setupRouter()
	token := registerAndLogin(router, t)

	// ---- 1. Create ----
	createBody := `{"bulan_tahun": "09-2025"}`
	req := httptest.NewRequest("POST", "/api/times", bytes.NewBufferString(createBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var createResp models.SuccessResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)
	assert.True(t, createResp.Success)
	assert.NotNil(t, createResp.Data)

	created := createResp.Data.(map[string]interface{})

	// ambil id dengan aman
	var id int
	if v, ok := created["id"].(float64); ok {
		id = int(v)
	} else if v, ok := created["ID"].(float64); ok {
		id = int(v)
	} else {
		t.Fatalf("ID not found in create response: %+v", created)
	}

	// ---- 2. Get All ----
	req2 := httptest.NewRequest("GET", "/api/times", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// ---- 3. Get By BulanTahun ----
	req3 := httptest.NewRequest("GET", "/api/times/09-2025", nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)

	// ---- 4. Update ----
	updateBody := `{"bulan_tahun": "10-2025"}`
	req4 := httptest.NewRequest("PUT", fmt.Sprintf("/api/times/%d", id), bytes.NewBufferString(updateBody))
	req4.Header.Set("Content-Type", "application/json")
	req4.Header.Set("Authorization", "Bearer "+token)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)

	// ---- 5. Delete ----
	req5 := httptest.NewRequest("DELETE", fmt.Sprintf("/api/times/%d", id), nil)
	req5.Header.Set("Authorization", "Bearer "+token)
	w5 := httptest.NewRecorder()
	router.ServeHTTP(w5, req5)
	assert.Equal(t, http.StatusOK, w5.Code)
}
