package login

import (
	loginv1 "github.com/runeharvest/gserver/proto/rh/login/v1"
)

// LoginClient is the use case handler for login-related operations.
type LoginClient struct {
	loginv1.LoginServiceClient
}

func NewLoginClient() (*LoginClient, error) {
	e := &LoginClient{}
	return e, nil
}
