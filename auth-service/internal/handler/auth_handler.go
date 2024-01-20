package handler

import (
	"context"

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
	err := handler.authService.Register(req.Username, req.Password)
	if err != nil {
		return nil, errToRPCError(err)
	}

	return &pb.RegisterResponse{}, nil
}

func (handler *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := handler.authService.Login(req.Username, req.Password)
	if err != nil {
		return nil, errToRPCError(err)
	}

	return &pb.LoginResponse{Token: token}, nil
}

func errToRPCError(err error) error {
	st := status.New(codes.Internal, err.Error())
	return st.Err()
}
