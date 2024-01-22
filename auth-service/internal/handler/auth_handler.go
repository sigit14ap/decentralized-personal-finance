package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/sigit14ap/personal-finance/auth-service/internal/domain"
	pb "github.com/sigit14ap/personal-finance/auth-service/internal/proto"
	"github.com/sigit14ap/personal-finance/auth-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	err := validateAuthRequest(req.Username, req.Password)
	if err != nil {
		return nil, errToRPCError(err)
	}

	err = handler.authService.Register(req.Username, req.Password)
	if err != nil {
		return nil, errToRPCError(err)
	}

	return &pb.RegisterResponse{}, nil
}

func (handler *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	err := validateAuthRequest(req.Username, req.Password)
	if err != nil {
		return nil, errToRPCError(err)
	}

	token, err := handler.authService.Login(req.Username, req.Password)
	if err != nil {
		return nil, errToRPCError(err)
	}

	return &pb.LoginResponse{Token: token}, nil
}

func validateAuthRequest(username, password string) error {
	validate := validator.New()

	request := domain.User{
		Username: username,
		Password: password,
	}

	err := validate.Struct(request)
	if err != nil {
		return err
	}

	return nil
}

func errToRPCError(err error) error {
	st := status.New(codes.Internal, err.Error())
	return st.Err()
}
