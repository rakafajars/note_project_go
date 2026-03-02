package config

import (
	"fmt"
	appconfig "notes-project/internal/config"
	"notes-project/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg *appconfig.Config) {
	// load file .env
	err := godotenv.Load()
	if err != nil {
		panic("Gagal memuat file .env")
	}

	// ambil data dari .env
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Gagal koneksi: %v", err))
	}

	// Auto Migrate: Membuat tabel otomatis berdasarkan struct Model
	database.AutoMigrate(&models.Note{})

	DB = database
	fmt.Println("Koneksi database berhasil")
}
