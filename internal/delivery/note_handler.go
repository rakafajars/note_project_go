package delivery

import (
	"net/http"
	"notes-project/internal/usecase"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	usecase usecase.NoteUsecase
}

func NewNoteHandler(u usecase.NoteUsecase) *NoteHandler {
	return &NoteHandler{usecase: u}
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	// Bind JSON dari body request ke struct input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.usecase.CreateNote(input.Title, input.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, note)

}

func (h *NoteHandler) GetAllNotes(c *gin.Context) {
	notes, err := h.usecase.GetAllNotes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notes)
}
