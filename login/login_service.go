package login

import (
	"context"

	"github.com/runeharvest/gserver/login/storage"
	loginv1 "github.com/runeharvest/gserver/proto/rh/login/v1"
)

// LoginService is the use case handler for login-related operations.
type LoginService struct {
	loginv1.UnimplementedLoginServiceServer
	storager storage.Storager
}

func NewLoginService(storage storage.Storager) (*LoginService, error) {
	e := &LoginService{storager: storage}
	return e, nil
}

func (e *LoginService) LoginVerify(ctx context.Context, in *loginv1.LoginVerifyRequest) (*loginv1.LoginVerifyResponse, error) {

	return &loginv1.LoginVerifyResponse{}, nil
}
