package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sigit14ap/github.com/sigit14ap/decentralized-personal-finance-auth-service/internal/domain"
	"github.com/sigit14ap/github.com/sigit14ap/decentralized-personal-finance-auth-service/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(email, password string) (domain.User, error)
	Login(email, password string) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
	FindUserById(id string) (domain.User, error)
}

type authService struct {
	userRepository *repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{
		userRepository: &userRepository,
	}
}

func (service *authService) Register(email, password string) (domain.User, error) {
	currentTime := time.Now()
	_, err := service.userRepository.FindUserByEmail(email)

	if err == nil {
		return domain.User{}, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return domain.User{}, err
	}

	userData := domain.User{
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: currentTime.Format("2006-01-02 15:04:05"),
		UpdatedAt: currentTime.Format("2006-01-02 15:04:05"),
	}

	user, err := service.userRepository.CreateUser(userData)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (service *authService) Login(email string, password string) (string, error) {
	user, err := service.userRepository.FindUserByEmail(email)

	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := createToken(user.ID.String())

	if err != nil {
		return "", errors.New("failed generate token")
	}

	return token, nil
}

func (service *authService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token or claims are invalid")
}

func (service *authService) FindUserById(id string) (domain.User, error) {
	user, err := service.userRepository.FindUserById(id)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func createToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userId,
		"expired_at": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
