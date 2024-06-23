package auth

import (
	"context"
	"strconv"

	ssov1 "github.com/atlasir0/proto/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

const (
	emptyValue = int64(0) 
)

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}
	appID, err := strconv.Atoi(req.GetAppId()) 
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid App ID")
	}
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), appID)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &ssov1.RegisterResponse{
		UsedId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "Email cannot be empty")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "Password cannot be empty")
	}
	if req.GetAppId() == "" {
		return status.Error(codes.InvalidArgument, "APPID is required")
	}

	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "Email cannot be empty")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "Password cannot be empty")
	}

	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue { // Изменение условия на сравнение с int64
		return status.Error(codes.InvalidArgument, "User ID is required")
	}

	return nil
}
