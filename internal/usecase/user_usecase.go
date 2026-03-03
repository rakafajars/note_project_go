package usecase

import (
	"errors"
	"notes-project/internal/models"
	"notes-project/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Login(email, password, secretKey string) (string, error)
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

func (u *userUsecase) Login(email, password, secretKey string) (string, error) {
	// cari user di database
	user, err := u.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("Email atau Password salah")
	}

	// 2. Bandingkan password bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", errors.New("Email atau Password salah")
	}

	// 3. Buat JWT Token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	return tokenString, err
}
