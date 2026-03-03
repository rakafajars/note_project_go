package delivery

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. ambil header auth
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			ErrorResponse(c, "Otorisasi diperlukan", http.StatusUnauthorized, "error", nil)
			c.Abort() //hentikan request di sini
			return
		}

		// 2. format bearer "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ErrorResponse(c, "Format token tidak valid", http.StatusUnauthorized, "error", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. parse dan validasi token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			ErrorResponse(c, "Token tidak valid atau kadaluarsa", http.StatusUnauthorized, "error", nil)
			c.Abort()
			return
		}

		// 4. Ambil user_id dari claims dan simpan ke context gin
		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["user_id"].(float64))

		// simpen userID agar bisa dipakai oleh handler (misal: saat simpan note)
		c.Set("user_id", userID)

		c.Next()

	}
}
