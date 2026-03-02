package delivery

import (
	"net/http"
	"notes-project/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	usecase usecase.NoteUsecase
}

func NewNoteHandler(u usecase.NoteUsecase) *NoteHandler {
	return &NoteHandler{usecase: u}
}

type NoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// CreateNote godoc
// @Summary      Membuat catatan baru
// @Description  Menyimpan judul catatan ke database
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        note  body      NoteRequest  true  "Note Data"
// @Success      201  {object}  Response
// @Failure      400  {object}  Response
// @Router       /notes [post]
func (h *NoteHandler) CreateNote(c *gin.Context) {
	var input NoteRequest

	// Bind JSON dari body request ke struct input
	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, "Validasi gagal", http.StatusBadRequest, "error", gin.H{"details": err.Error()})
		return
	}

	note, err := h.usecase.CreateNote(input.Title, input.Content)
	if err != nil {
		ErrorResponse(c, "Gagal membuat catatan", http.StatusInternalServerError, "error", gin.H{"details": err.Error()})
		return
	}

	SuccessResponse(c, "Catatan berhasil dibuat", http.StatusCreated, "success", note, nil)

}

// GetAllNotes godoc
// @Summary      Mendapatkan semua catatan
// @Description  Mengambil semua data catatan dari database
// @Tags         notes
// @Produce      json
// @Success      200  {object}  Response
// @Failure      500  {object}  Response
// @Router       /notes [get]
func (h *NoteHandler) GetAllNotes(c *gin.Context) {
	notes, err := h.usecase.GetAllNotes()
	if err != nil {
		ErrorResponse(c, "Gagal memuat catatan", http.StatusInternalServerError, "error", gin.H{"details": err.Error()})
		return
	}

	SuccessResponse(c, "Berhasil mendapatkan catatan", http.StatusOK, "success", notes, nil)

}

// DeleteNote godoc
// @Summary      Menghapus catatan
// @Description  Menghapus data berdasarkan ID
// @Tags         notes
// @Param        id   path      int  true  "Note ID"
// @Success      200  {object}  Response
// @Failure      404  {object}  Response
// @Router       /notes/{id} [delete]
func (h *NoteHandler) DeleteNote(c *gin.Context) {
	// Mengambil ID dari URL parameter /notes/:id
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam) // Convert string ke int

	err := h.usecase.DeleteNote(uint(id))
	if err != nil {
		if err.Error() == "catatan tidak ditemukan" {
			ErrorResponse(c, err.Error(), http.StatusNotFound, "error", nil)
			return
		}
		ErrorResponse(c, "Terjadi kesalahan server", http.StatusInternalServerError, "error", nil)
		return
	}

	SuccessResponse(c, "Catatan Berhasil dihapus", http.StatusOK, "success", nil, nil)
}

// UpdateNote godoc
// @Summary      Memperbarui catatan
// @Description  Memperbarui judul dan isi catatan berdasarkan ID
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Note ID"
// @Param        note body      NoteRequest  true  "Note Data"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      404  {object}  Response
// @Failure      500  {object}  Response
// @Router       /notes/{id} [put]
func (h *NoteHandler) UpdateNote(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var input NoteRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, "input tidak valid", http.StatusBadRequest, "error", nil)
		return
	}

	note, err := h.usecase.UpdateNote(uint(id), input.Title, input.Content)
	if err != nil {
		if err.Error() == "catatan tidak ditemukan" {
			ErrorResponse(c, err.Error(), http.StatusNotFound, "error", nil)
			return
		}

		ErrorResponse(c, "Gagal memperbarui catatan", http.StatusInternalServerError, "error", nil)
		return
	}

	SuccessResponse(c, "Catatan Berhasil Diperbarui", http.StatusOK, "success", note, nil)
}
