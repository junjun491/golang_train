package middleware

import (
	"net/http"
	"strings"

	"golang_train/backend-go/internal/auth"

	"github.com/gin-gonic/gin"
)

func AuthenticateAPI() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": []string{"Authorization header missing"},
			})
			c.Abort()
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": []string{"Invalid authorization header"},
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, prefix)

		teacherID, err := auth.ParseTeacherToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": []string{"Invalid token"},
			})
			c.Abort()
			return
		}

		c.Set("current_teacher_id", teacherID)
		c.Next()
	}
}