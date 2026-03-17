package main

import (
	"golang_train/backend-go/internal/db"
	"golang_train/backend-go/internal/handler"
	"golang_train/backend-go/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDB()
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	r.GET("/teachers", handler.GetTeachers)
	r.GET("/teachers/:id", handler.GetTeacher)
	r.POST("/teachers/register", handler.RegisterTeacher)
	r.POST("/teachers/login", handler.LoginTeacher)
	r.GET("/teachers/me", middleware.AuthenticateAPI(), handler.GetMe)

	r.Run(":8080")
}