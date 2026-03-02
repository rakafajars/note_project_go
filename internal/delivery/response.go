package delivery

import "github.com/gin-gonic/gin"

type Response struct {
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	StatusCode uint        `json:"status_code"`
	Data       interface{} `json:"data,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
}

func APIResponse(c *gin.Context, message string, statusCode int, status string, data interface{}) {
	jsonResponse := Response{
		Status:     status,
		StatusCode: uint(statusCode),
		Message:    message,
		Data:       data,
	}

	c.JSON(statusCode, jsonResponse)
}
