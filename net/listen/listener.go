package listen

import (
	"context"

	loginv1 "github.com/runeharvest/gserver/proto/rh/login/v1"
)

// Listener defines the interface for network operations.
type Listener interface {
	LoginRegister(loginService loginv1.LoginServiceServer) error
	LoginVerify(ctx context.Context, in *loginv1.LoginVerifyRequest) (*loginv1.LoginVerifyResponse, error)
}
