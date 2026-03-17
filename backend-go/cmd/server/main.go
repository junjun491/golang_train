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

	// 既存の healthz は残す
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 既存ルートも残す（curl確認や移行中の保険）
	r.GET("/teachers", handler.GetTeachers)
	r.GET("/teachers/:id", handler.GetTeacher)
	r.POST("/teachers/register", handler.RegisterTeacher)
	r.POST("/teachers/login", handler.LoginTeacher)
	r.GET("/teachers/me", middleware.AuthenticateAPI(), handler.GetMe)

	// React / 本番ALB 用に /api 配下を追加
	api := r.Group("/api")
	{
		api.GET("/healthz", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		api.GET("/teachers", handler.GetTeachers)
		api.GET("/teachers/:id", handler.GetTeacher)
		api.POST("/teachers/register", handler.RegisterTeacher)

		// React側が叩いている sign_in に合わせる
		api.POST("/teachers/sign_in", handler.LoginTeacher)

		// 既存の login も /api 配下で使えるようにしておく
		api.POST("/teachers/login", handler.LoginTeacher)

		api.GET("/teachers/profile", middleware.AuthenticateAPI(), handler.GetMe)
	}

	log.Printf("starting server on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}