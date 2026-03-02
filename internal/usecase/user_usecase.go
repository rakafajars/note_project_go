package usecase

import (
	"notes-project/internal/models"
	"notes-project/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(email, password string) (*models.User, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) UserUsecase {
	return &userUsecase{
		repo: r,
	}
}

func (u *userUsecase) Register(email, password string) (*models.User, error) {
	// 1, hash passwor dmenggunakan bcrpy
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// buat objeck baru
	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	// simpan ke database
	err = u.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
