package main

import (
	"context"
	"log"
	"net/http"
	"notes-project/config"
	appconfig "notes-project/internal/config"
	"notes-project/internal/delivery"
	"notes-project/internal/models"
	"notes-project/internal/repository"
	"notes-project/internal/usecase"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "notes-project/docs" // Import library CORS

	"github.com/gin-contrib/cors" // Import library CORS
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Notes API
// @version         1.0
// @description     Ini adalah server API untuk aplikasi catatan (Notes).
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
	cfg := appconfig.LoadConfig()
	// 1. inisialisasi koneksi database
	config.ConnectDatabase(cfg)

	// 2. setup layers (Depedency Injection)
	noteRepo := repository.NewNoteRepository(config.DB)
	noteUsecase := usecase.NewTodoUsecase(noteRepo)
	noteHandler := delivery.NewNoteHandler(noteUsecase)

	config.DB.AutoMigrate(&models.User{}, &models.Note{})

	userRepo := repository.NewUserRepository(config.DB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := delivery.NewUserHandler(userUsecase)

	// 3. Setup Router (Gin)
	r := gin.Default()

	// 1. Logger & Recovery (Bawaan Gin)
	// Logger: Mencatat setiap request di terminal
	// Recovery: Mencegah server mati jika terjadi panic/error fatal
	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	// 2. custom cors middleware
	// ini akan mengizinakn aplikasi
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	// Route untuk swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	v1 := r.Group("/api/v1")
	{

		v1.POST("/register", userHandler.Register)

		v1.POST("/login", func(c *gin.Context) {
			userHandler.Login(c, cfg.JWTSecret)
		})

		notes := v1.Group("/notes")
		notes.Use(delivery.AuthMiddleware(cfg.JWTSecret))
		{
			notes.POST("", noteHandler.CreateNote)
			notes.GET("", noteHandler.GetAllNotes)
			notes.DELETE("/:id", noteHandler.DeleteNote)
			notes.PUT("/:id", noteHandler.UpdateNote)
		}

	}

	// 1. Definisikan konfigurasi HTTP Server secara manual
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// 2. Jalankan server di dalam goroutine (jalur terpisah)
	// Ini agar program tidak berhenti di sini dan bisa lanjut ke baris berikutnya
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Gagal menjalankan server: %s\n", err)
		}
	}()

	// 3. Tunggu sinyal interupsi untuk mematikan server secara halus
	// quit adalah channel yang akan menerima sinyal dari Sistem Operasi (seperti Ctrl+C)
	quit := make(chan os.Signal, 1)

	// SIGINT: sinyal dari Ctrl+C
	// SIGTERM: sinyal stop standar dari sistem (misal saat deploy)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Mematikan server secara halus...")

	// 4. Berikan batas waktu (timeout) 5 detik untuk menyelesaikan request yang tersisa
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server dimatikan paksa:", err)
	}

	log.Println("Server berhasil berhenti dengan aman")

}
