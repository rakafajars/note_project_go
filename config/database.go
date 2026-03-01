package config

import (
	"fmt"
	"notes-project/internal/models"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// load file .env
	err := godotenv.Load()
	if err != nil {
		panic("Gagal memuat file .env")
	}

	// ambil data dari .env
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disabled",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal koneksi ke datbase!")
	}

	// Auto Migrate: Membuat tabel otomatis berdasarkan struct Model
	database.AutoMigrate(&models.Note{})

	DB = database
	fmt.Println("Koneksi database berhasil")
}
