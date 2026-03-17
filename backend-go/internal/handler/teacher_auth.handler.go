package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"golang_train/backend-go/internal/auth"
	"golang_train/backend-go/internal/db"

	"github.com/gin-gonic/gin"
)

type RegisterTeacherRequest struct {
	Teacher RegisterTeacherParams `json:"teacher"`
}

type RegisterTeacherParams struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type TeacherResponseData struct {
	ID    int     `json:"id"`
	Name  *string `json:"name"`
	Email string  `json:"email"`
}

type TeacherResponse struct {
	Data TeacherResponseData `json:"data"`
}

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func RegisterTeacher(c *gin.Context) {
	var req RegisterTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors: []string{"invalid request"},
		})
		return
	}

	req.Teacher.Email = strings.TrimSpace(req.Teacher.Email)
	req.Teacher.Name = strings.TrimSpace(req.Teacher.Name)

	if req.Teacher.Email == "" {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Errors: []string{"email can't be blank"},
		})
		return
	}

	if req.Teacher.Password == "" {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Errors: []string{"password can't be blank"},
		})
		return
	}

	if req.Teacher.Password != req.Teacher.PasswordConfirmation {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Errors: []string{"password confirmation doesn't match"},
		})
		return
	}

	var existingID int
	err := db.DB.QueryRow(
		context.Background(),
		`SELECT id FROM teachers WHERE email = $1`,
		req.Teacher.Email,
	).Scan(&existingID)

	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Errors: []string{"email has already been taken"},
		})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Teacher.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Errors: []string{"failed to hash password"},
		})
		return
	}

	now := time.Now()

	var id int
	var name *string
	var email string

	err = db.DB.QueryRow(
		context.Background(),
		`
		INSERT INTO teachers (
			email,
			encrypted_password,
			created_at,
			updated_at,
			name
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, email
		`,
		req.Teacher.Email,
		hashedPassword,
		now,
		now,
		req.Teacher.Name,
	).Scan(&id, &name, &email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Errors: []string{err.Error()},
		})
		return
	}

	token, err := auth.GenerateTeacherToken(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Errors: []string{"failed to generate token"},
		})
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusCreated, TeacherResponse{
		Data: TeacherResponseData{
			ID:    id,
			Name:  name,
			Email: email,
		},
	})
}