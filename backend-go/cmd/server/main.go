package main

import (
	"log"

	"golang_train/backend-go/internal/auth"
	"golang_train/backend-go/internal/config"
	"golang_train/backend-go/internal/db"
	"golang_train/backend-go/internal/handler"
	"golang_train/backend-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	auth.InitJWT(cfg.JWTSecret)
	db.ConnectDB(cfg.DatabaseURL)

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/teachers", handler.GetTeachers)
	r.GET("/teachers/:id", handler.GetTeacher)
	r.POST("/teachers/register", handler.RegisterTeacher)
	r.POST("/teachers/login", handler.LoginTeacher)
	r.GET("/teachers/me", middleware.AuthenticateAPI(), handler.GetMe)

	r.Run(":8080")
}