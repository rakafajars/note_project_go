package main

import (
	"notes-project/config"
	"notes-project/internal/delivery"
	"notes-project/internal/repository"
	"notes-project/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. inisialisasi koneksi database
	config.ConnectDatabase()

	// 2. setup layers (Depedency Injection)
	noteRepo := repository.NewNoteRepository(config.DB)
	noteUsecase := usecase.NewTodoUsecase(noteRepo)
	noteHandler := delivery.NewNoteHandler(noteUsecase)

	// 3. Setup Router (Gin)
	r := gin.Default()

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
