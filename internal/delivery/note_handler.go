package delivery

import (
	"net/http"
	"notes-project/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type NoteHandler struct {
	usecase usecase.NoteUsecase
}

func NewNoteHandler(u usecase.NoteUsecase) *NoteHandler {
	return &NoteHandler{usecase: u}
}

type NoteRequest struct {
	Title   string `json:"title" validate:"required,min=5,max=100" example:"Catatan Pertama"`
	Content string `json:"content" validate:"required,min=5" example:"Ini adalah isi catatan pertama saya"`
}

// NoteExample adalah contoh data catatan untuk Swagger
type NoteExample struct {
	ID        uint   `json:"id" example:"1"`
	Title     string `json:"title" example:"Catatan Pertama"`
	Content   string `json:"content" example:"Ini adalah isi catatan pertama saya"`
	CreatedAt string `json:"created_at" example:"2026-03-03T01:00:00+07:00"`
	UpdatedAt string `json:"updated_at" example:"2026-03-03T01:00:00+07:00"`
}

// PaginationMeta contoh meta data pagination
type PaginationMeta struct {
	CurrentPage int   `json:"current_page" example:"1"`
	Limit       int   `json:"limit" example:"10"`
	Total       int64 `json:"total" example:"25"`
}

// CreateNoteSuccessResponse contoh response sukses create note
type CreateNoteSuccessResponse struct {
	Status     string      `json:"status" example:"success"`
	Message    string      `json:"message" example:"Catatan berhasil dibuat"`
	StatusCode uint        `json:"status_code" example:"201"`
	Data       NoteExample `json:"data"`
}

// GetAllNotesSuccessResponse contoh response sukses get all notes
type GetAllNotesSuccessResponse struct {
	Status     string         `json:"status" example:"success"`
	Message    string         `json:"message" example:"Berhasil mendapatkan catatan"`
	StatusCode uint           `json:"status_code" example:"200"`
	Data       []NoteExample  `json:"data"`
	Meta       PaginationMeta `json:"meta"`
}

// DeleteNoteSuccessResponse contoh response sukses delete note
type DeleteNoteSuccessResponse struct {
	Status     string `json:"status" example:"success"`
	Message    string `json:"message" example:"Catatan Berhasil dihapus"`
	StatusCode uint   `json:"status_code" example:"200"`
}

// UpdateNoteSuccessResponse contoh response sukses update note
type UpdateNoteSuccessResponse struct {
	Status     string      `json:"status" example:"success"`
	Message    string      `json:"message" example:"Catatan Berhasil Diperbarui"`
	StatusCode uint        `json:"status_code" example:"200"`
	Data       NoteExample `json:"data"`
}

// ErrorCommonResponse contoh response error umum
type ErrorCommonResponse struct {
	Status     string `json:"status" example:"error"`
	Message    string `json:"message" example:"Terjadi kesalahan"`
	StatusCode uint   `json:"status_code" example:"500"`
}

// ValidationErrorDetail contoh detail error validasi
type ValidationErrorDetail struct {
	Field   string `json:"field" example:"Title"`
	Message string `json:"message" example:"Title minimal 5 karakter"`
}

// ErrorValidationResponse contoh response error validasi
type ErrorValidationResponse struct {
	Status     string                  `json:"status" example:"error"`
	Message    string                  `json:"message" example:"validasi gagal"`
	StatusCode uint                    `json:"status_code" example:"400"`
	Errors     []ValidationErrorDetail `json:"errors"`
}

// CreateNote godoc
// @Summary      Membuat catatan baru
// @Description  Menyimpan judul dan konten catatan ke database
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        note  body      NoteRequest  true  "Note Data"
// @Success      201  {object}  CreateNoteSuccessResponse
// @Failure      400  {object}  ErrorValidationResponse
// @Failure      500  {object}  ErrorCommonResponse
// @Router       /notes [post]
func (h *NoteHandler) CreateNote(c *gin.Context) {
	var input NoteRequest

	val, _ := c.Get("user_id")
	userID := val.(uint)

	// Bind JSON dari body request ke struct input
	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, "Validasi gagal", http.StatusBadRequest, "error", gin.H{"details": err.Error()})
		return
	}

	// 3. jalankan validasi
	if err := validate.Struct(input); err != nil {

		report := []ValidationError{}
		for _, err := range err.(validator.ValidationErrors) {
			// ambil nama tag (misal: min)
			tag := err.Tag()
			// ambil parameter tag (misal: 5)
			param := err.Param()

			message := "Format " + err.Field() + " tidak sesuai"

			// logika untuk mempercantik pesan
			switch tag {
			case "min":
				message = err.Field() + " minimal " + param + " karakter"
			case "required":
				message = err.Field() + " wajib diisi"
			case "max":
				message = err.Field() + " maksimal " + param + " karakter"
			}

			report = append(report, ValidationError{
				Field:   err.Field(),
				Message: message,
			})
		}

		ErrorResponse(c, "validasi gagal", http.StatusBadRequest, "error", report)
		return
	}

	note, err := h.usecase.CreateNote(input.Title, input.Content, userID)
	if err != nil {
		ErrorResponse(c, "Gagal membuat catatan", http.StatusInternalServerError, "error", gin.H{"details": err.Error()})
		return
	}

	SuccessResponse(c, "Catatan berhasil dibuat", http.StatusCreated, "success", note, nil)

}

// GetAllNotes godoc
// @Summary      Mendapatkan semua catatan
// @Description  Mengambil semua data catatan dari database dengan fitur pencarian dan pagination
// @Tags         notes
// @Produce      json
// @Param        q      query     string  false  "Kata kunci pencarian berdasarkan judul"
// @Param        page   query     int     false  "Nomor halaman (default: 1)"
// @Param        limit  query     int     false  "Jumlah data per halaman (default: 10)"
// @Success      200  {object}  GetAllNotesSuccessResponse
// @Failure      500  {object}  ErrorCommonResponse
// @Router       /notes [get]
func (h *NoteHandler) GetAllNotes(c *gin.Context) {
	// ambil query params
	val, _ := c.Get("user_id")
	userID := val.(uint)

	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// panggil use case
	notes, total, err := h.usecase.GetAllNotes(userID, query, page, limit)
	if err != nil {
		ErrorResponse(c, "Gagal memuat catatan", http.StatusInternalServerError, "error", gin.H{"details": err.Error()})
		return
	}

	// 3. susu meta data pagination
	meta := gin.H{
		"current_page": page,
		"limit":        limit,
		"total":        total,
	}

	SuccessResponse(c, "Berhasil mendapatkan catatan", http.StatusOK, "success", notes, meta)

}

// DeleteNote godoc
// @Summary      Menghapus catatan
// @Description  Menghapus data catatan berdasarkan ID
// @Tags         notes
// @Produce      json
// @Param        id   path      int  true  "Note ID"
// @Success      200  {object}  DeleteNoteSuccessResponse
// @Failure      404  {object}  ErrorCommonResponse
// @Failure      500  {object}  ErrorCommonResponse
// @Router       /notes/{id} [delete]
func (h *NoteHandler) DeleteNote(c *gin.Context) {

	val, _ := c.Get("user_id")
	userID := val.(uint)
	// Mengambil ID dari URL parameter /notes/:id
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam) // Convert string ke int

	err := h.usecase.DeleteNote(uint(id), userID)
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
// @Success      200  {object}  UpdateNoteSuccessResponse
// @Failure      400  {object}  ErrorValidationResponse
// @Failure      404  {object}  ErrorCommonResponse
// @Failure      500  {object}  ErrorCommonResponse
// @Router       /notes/{id} [put]
func (h *NoteHandler) UpdateNote(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID := val.(uint)

	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var input NoteRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, "input tidak valid", http.StatusBadRequest, "error", nil)
		return
	}

	note, err := h.usecase.UpdateNote(uint(id), input.Title, input.Content, userID)
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
