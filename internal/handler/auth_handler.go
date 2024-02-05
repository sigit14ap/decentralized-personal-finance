package handler

import (
	"context"

	"github.com/sigit14ap/github.com/sigit14ap/decentralized-personal-finance-auth-service/internal/domain"
	"github.com/sigit14ap/github.com/sigit14ap/decentralized-personal-finance-auth-service/internal/helpers"
	pb "github.com/sigit14ap/github.com/sigit14ap/decentralized-personal-finance-auth-service/internal/proto"
	"github.com/sigit14ap/github.com/sigit14ap/decentralized-personal-finance-auth-service/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (handler *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var responseError = []*pb.ResponseError{}

	request := domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	validate := helpers.ValidateRequest(request)
	if validate != nil {
		return &pb.RegisterResponse{
			Status:  422,
			Message: "Validation failed",
			Error:   validate,
		}, nil
	}

	_, err := handler.authService.Register(req.Email, req.Password)
	if err != nil {
		return &pb.RegisterResponse{
			Status:  400,
			Message: err.Error(),
			Error:   responseError,
		}, nil
	}

	return &pb.RegisterResponse{
		Status:  200,
		Message: "Register success",
		Error:   responseError,
	}, nil
}

func (handler *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var responseError = []*pb.ResponseError{}

	request := domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	validate := helpers.ValidateRequest(request)
	if validate != nil {
		return &pb.LoginResponse{
			Status:  422,
			Message: "Validation failed",
			Error:   validate,
		}, nil
	}

	token, err := handler.authService.Login(req.Email, req.Password)
	if err != nil {
		return &pb.LoginResponse{
			Status:  400,
			Message: err.Error(),
			Error:   responseError,
		}, nil
	}

	return &pb.LoginResponse{
		Status:  200,
		Message: "Login success",
		Data: &pb.DataLoginResponse{
			AccessToken: token,
		},
	}, nil
}

func (handler *AuthHandler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	_, err := handler.authService.ValidateToken(req.Token)

	if err != nil {
		return &pb.ValidateResponse{
			Status:  401,
			Message: "Unauthorized",
		}, nil
	}

	return &pb.ValidateResponse{
		Status:  200,
		Message: "Authorized",
	}, nil
}
