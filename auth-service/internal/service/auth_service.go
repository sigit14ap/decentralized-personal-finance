package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sigit14ap/personal-finance/auth-service/internal/domain"
	"github.com/sigit14ap/personal-finance/auth-service/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
}

type authService struct {
	userRepository *repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{
		userRepository: &userRepository,
	}
}

func (service *authService) Register(username, password string) error {
	_, err := service.userRepository.FindUserByUsername(username)

	if err == nil {
		return errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user := domain.User{
		Username: username,
		Password: string(hashedPassword),
	}

	err = service.userRepository.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (service *authService) Login(username string, password string) (string, error) {
	user, err := service.userRepository.FindUserByUsername(username)

	if err != nil {
		return "", errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", errors.New("invalid username or password")
	}

	token, err := createToken(username)

	if err != nil {
		return "", errors.New("failed generate token")
	}

	return token, nil
}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
