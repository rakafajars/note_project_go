package delivery

import "github.com/gin-gonic/gin"

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Response struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode uint   `json:"status_code"`
	Data       any    `json:"data,omitempty"`
	Meta       any    `json:"meta,omitempty"`   // Tambahkan Meta di struct
	Errors     any    `json:"errors,omitempty"` // Gunakan any agar lebih fleksibel
}

func SuccessResponse(c *gin.Context, message string, statusCode int, status string, data any, meta any) {
	// Gunakan struct Response agar konsisten
	jsonResponse := Response{
		Status:     status,
		StatusCode: uint(statusCode),
		Message:    message,
		Data:       data,
		Meta:       meta,
	}

	c.JSON(statusCode, jsonResponse) // Gunakan parameter statusCode untuk Header
}

func ErrorResponse(c *gin.Context, message string, statusCode int, status string, errs any) {
	jsonResponse := Response{
		Status:     status,
		StatusCode: uint(statusCode),
		Message:    message,
		Errors:     errs,
	}

	c.JSON(statusCode, jsonResponse)
}
