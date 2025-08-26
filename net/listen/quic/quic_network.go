package quic

import (
	"context"

	loginv1 "github.com/runeharvest/gserver/proto/rh/login/v1"
)

type QuicNetwork struct {
}

func (q *QuicNetwork) LoginVerify(ctx context.Context, in *loginv1.LoginVerifyRequest) (*loginv1.LoginVerifyResponse, error) {
	// Implement the QUIC call to the login service here.
	return &loginv1.LoginVerifyResponse{}, nil
}
