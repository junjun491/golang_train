package handler

import (
	"net/http"
	"strings"

	"golang_train/backend-go/internal/auth"
	"golang_train/backend-go/internal/repository"

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

type LoginTeacherRequest struct {
	Teacher LoginTeacherParams `json:"teacher"`
}

type LoginTeacherParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	existingTeacher, err := repository.FindTeacherByEmail(req.Teacher.Email)
	if err == nil && existingTeacher != nil {
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

	teacher, err := repository.CreateTeacher(
		req.Teacher.Email,
		hashedPassword,
		req.Teacher.Name,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Errors: []string{err.Error()},
		})
		return
	}

	token, err := auth.GenerateToken("teacher", int64(teacher.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Errors: []string{"failed to generate token"},
		})
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusCreated, TeacherResponse{
		Data: TeacherResponseData{
			ID:    teacher.ID,
			Name:  teacher.Name,
			Email: teacher.Email,
		},
	})
}

func LoginTeacher(c *gin.Context) {
	var req LoginTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors: []string{"invalid request"},
		})
		return
	}

	req.Teacher.Email = strings.TrimSpace(req.Teacher.Email)

	teacher, err := repository.FindTeacherByEmail(req.Teacher.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Errors: []string{"Invalid email or password"},
		})
		return
	}

	if !auth.CheckPassword(req.Teacher.Password, teacher.EncryptedPassword) {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Errors: []string{"Invalid email or password"},
		})
		return
	}

	token, err := auth.GenerateToken("teacher", int64(teacher.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Errors: []string{"failed to generate token"},
		})
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, TeacherResponse{
		Data: TeacherResponseData{
			ID:    teacher.ID,
			Name:  teacher.Name,
			Email: teacher.Email,
		},
	})
}

func GetMe(c *gin.Context) {
	currentTeacherIDValue, exists := c.Get("current_teacher_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Errors: []string{"Unauthorized"},
		})
		return
	}

	currentTeacherID := currentTeacherIDValue.(int)

	teacher, err := repository.FindTeacherByID(currentTeacherID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Errors: []string{"Unauthorized"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":    teacher.ID,
			"name":  teacher.Name,
			"email": teacher.Email,
		},
	})
}