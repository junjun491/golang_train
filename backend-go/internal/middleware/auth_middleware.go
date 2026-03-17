package middleware

import (
	"net/http"
	"strconv"
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

		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": []string{"Invalid token"},
			})
			c.Abort()
			return
		}

		subValue, ok := claims["sub"].(string)
		if !ok || subValue == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": []string{"Invalid token subject"},
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(subValue, ":", 2)
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": []string{"Invalid token subject"},
			})
			c.Abort()
			return
		}

		role := parts[0]
		idStr := parts[1]

		if role != "teacher" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": []string{"Invalid token role"},
			})
			c.Abort()
			return
		}

		teacherID, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": []string{"Invalid token subject"},
			})
			c.Abort()
			return
		}

		c.Set("jwt_claims", claims)
		c.Set("current_teacher_id", teacherID)
		c.Next()
	}
}