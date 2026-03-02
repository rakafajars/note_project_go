package main

import (
	"notes-project/config"
	"notes-project/internal/delivery"
	"notes-project/internal/repository"
	"notes-project/internal/usecase"

	_ "notes-project/docs"

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
	// 1. inisialisasi koneksi database
	config.ConnectDatabase()

	// 2. setup layers (Depedency Injection)
	noteRepo := repository.NewNoteRepository(config.DB)
	noteUsecase := usecase.NewTodoUsecase(noteRepo)
	noteHandler := delivery.NewNoteHandler(noteUsecase)

	// 3. Setup Router (Gin)
	r := gin.Default()

	// Route untuk swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	v1 := r.Group("/api/v1")
	{
		v1.POST("/notes", noteHandler.CreateNote)
		v1.GET("/notes", noteHandler.GetAllNotes)
		v1.DELETE("/notes/:id", noteHandler.DeleteNote)
		v1.PUT("/notes/:id", noteHandler.UpdateNote)
	}

	// 4. Jalankan Server
	r.Run(":8080")
}
