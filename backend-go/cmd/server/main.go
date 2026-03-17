package main

import (
	"golang_train/backend-go/internal/db"
	"golang_train/backend-go/internal/handler"
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

	r.Run(":8080")
}