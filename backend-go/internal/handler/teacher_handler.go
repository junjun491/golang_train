package handler

import (
	"context"
	"net/http"
	"strconv"

	"golang_train/backend-go/internal/db"
	"golang_train/backend-go/internal/model"

	"github.com/gin-gonic/gin"
)

func GetTeachers(c *gin.Context) {
	rows, err := db.DB.Query(context.Background(), `
		SELECT id, email, encrypted_password, reset_password_token,
		       reset_password_sent_at, remember_created_at,
		       created_at, updated_at, name
		FROM teachers
		ORDER BY id ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	teachers := []model.Teacher{}

	for rows.Next() {
		var teacher model.Teacher
		if err := rows.Scan(
			&teacher.ID,
			&teacher.Email,
			&teacher.EncryptedPassword,
			&teacher.ResetPasswordToken,
			&teacher.ResetPasswordSentAt,
			&teacher.RememberCreatedAt,
			&teacher.CreatedAt,
			&teacher.UpdatedAt,
			&teacher.Name,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		teachers = append(teachers, teacher)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teachers)
}

func GetTeacher(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var teacher model.Teacher

	err = db.DB.QueryRow(context.Background(), `
		SELECT id, email, encrypted_password,
		       reset_password_token, reset_password_sent_at,
		       remember_created_at, created_at, updated_at, name
		FROM teachers
		WHERE id = $1
	`, id).Scan(
		&teacher.ID,
		&teacher.Email,
		&teacher.EncryptedPassword,
		&teacher.ResetPasswordToken,
		&teacher.ResetPasswordSentAt,
		&teacher.RememberCreatedAt,
		&teacher.CreatedAt,
		&teacher.UpdatedAt,
		&teacher.Name,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "teacher not found"})
		return
	}

	c.JSON(http.StatusOK, teacher)
}