package delivery

import (
	"net/http"
	"notes-project/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(u usecase.UserUsecase) *UserHandler {
	return &UserHandler{u}
}

func (h *UserHandler) Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, "Input tidak valid", http.StatusBadRequest, "error", nil)
		return
	}

	// jalankan validasi
	user, err := h.usecase.Register(input.Email, input.Password)
	if err != nil {
		ErrorResponse(c, "Email sudah terdaftar atau terjadi kesalahan", http.StatusBadRequest, "error", nil)
		return
	}

	SuccessResponse(c, "Registrasi Berhasil", http.StatusCreated, "success", user, nil)
}

func (h *UserHandler) Login(c *gin.Context, secretKey string) {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, "Input tidak valid", http.StatusBadRequest, "error", nil)
		return
	}

	// panggil usecase login
	token, err := h.usecase.Login(input.Email, input.Password, secretKey)
	if err != nil {
		ErrorResponse(c, err.Error(), http.StatusUnauthorized, "error", nil)
		return
	}

	// Kembalikan token ke client
	SuccessResponse(c, "Login Berhasil", http.StatusOK, "success", gin.H{"token": token}, nil)
}
