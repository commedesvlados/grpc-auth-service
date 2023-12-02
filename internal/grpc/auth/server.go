package auth

import (
	"context"
	"errors"

	gspv1 "github.com/commedesvlados/grpc-service-protos/gen/go/grpc_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/commedesvlados/grpc_service/internal/servises/auth"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int32) (token string, err error)
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
}

type serverAPI struct {
	gspv1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	gspv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

const (
	emptyAppID  = 0
	emptyUserID = 0
)

func (s *serverAPI) Register(ctx context.Context, req *gspv1.RegisterRequest) (*gspv1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &gspv1.RegisterResponse{UserId: userID}, nil
}

func validateRegister(req *gspv1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func (s *serverAPI) Login(ctx context.Context, req *gspv1.LoginRequest) (*gspv1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &gspv1.LoginResponse{Token: token}, nil
}

func validateLogin(req *gspv1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyAppID {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *gspv1.IsAdminRequest) (*gspv1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &gspv1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func validateIsAdmin(req *gspv1.IsAdminRequest) error {
	if req.GetUserId() == emptyUserID {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	return nil
}
