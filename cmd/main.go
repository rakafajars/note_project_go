package main

import (
	"notes-project/config"
	appconfig "notes-project/internal/config"
	"notes-project/internal/delivery"
	"notes-project/internal/models"
	"notes-project/internal/repository"
	"notes-project/internal/usecase"

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

	config.DB.AutoMigrate(&models.User{})

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
		v1.POST("/notes", noteHandler.CreateNote)
		v1.GET("/notes", noteHandler.GetAllNotes)
		v1.DELETE("/notes/:id", noteHandler.DeleteNote)
		v1.PUT("/notes/:id", noteHandler.UpdateNote)

		v1.POST("/register", userHandler.Register)

		v1.POST("/login", func(c *gin.Context) {
			userHandler.Login(c, cfg.JWTSecret)
		})
	}

	// 4. Jalankan Server
	r.Run(":8080")
}
