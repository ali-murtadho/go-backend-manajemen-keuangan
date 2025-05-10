package database

import (
	"backend-api/config"
	"backend-api/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres" // Menggunakan driver PostgreSQL
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	// Load konfigurasi database dari .env
	dbUser := config.GetEnv("DB_USER", "postgres") // Default user PostgreSQL biasanya 'postgres'
	dbPass := config.GetEnv("DB_PASS", "")
	dbHost := config.GetEnv("DB_HOST", "localhost")
	dbPort := config.GetEnv("DB_PORT", "5432")        // Port default PostgreSQL adalah 5432
	dbName := config.GetEnv("DB_NAME", "db_keuangan") // Nama database, misalnya 'mydb'

	// Format DSN untuk PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	// Koneksi ke database PostgreSQL
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connected successfully!")

	// **Auto Migrate Models**
	err = DB.AutoMigrate(&models.User{}) // Tambahkan model lain jika perlu
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database migrated successfully! good")
}
