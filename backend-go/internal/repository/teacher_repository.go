package repository

import (
	"context"
	"time"

	"golang_train/backend-go/internal/db"
	"golang_train/backend-go/internal/model"
)

func FindTeacherByEmail(email string) (*model.Teacher, error) {
	var teacher model.Teacher

	err := db.DB.QueryRow(
		context.Background(),
		`
		SELECT id, email, encrypted_password,
		       reset_password_token, reset_password_sent_at,
		       remember_created_at, created_at, updated_at, name
		FROM teachers
		WHERE email = $1
		`,
		email,
	).Scan(
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
		return nil, err
	}

	return &teacher, nil
}

func FindTeacherByID(id int) (*model.Teacher, error) {
	var teacher model.Teacher

	err := db.DB.QueryRow(
		context.Background(),
		`
		SELECT id, email, encrypted_password,
		       reset_password_token, reset_password_sent_at,
		       remember_created_at, created_at, updated_at, name
		FROM teachers
		WHERE id = $1
		`,
		id,
	).Scan(
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
		return nil, err
	}

	return &teacher, nil
}

func CreateTeacher(email string, hashedPassword string, name string) (*model.Teacher, error) {
	now := time.Now()

	var teacher model.Teacher

	err := db.DB.QueryRow(
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
		RETURNING id, email, encrypted_password,
		          reset_password_token, reset_password_sent_at,
		          remember_created_at, created_at, updated_at, name
		`,
		email,
		hashedPassword,
		now,
		now,
		name,
	).Scan(
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
		return nil, err
	}

	return &teacher, nil
}