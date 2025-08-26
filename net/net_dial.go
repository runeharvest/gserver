package network

import (
	"context"
	"fmt"

	"github.com/runeharvest/gserver/net/dial"
	loginv1 "github.com/runeharvest/gserver/proto/rh/login/v1"
)

type NetDialService struct {
	dialer dial.Dialer
}

func NewNetDialService(dialer dial.Dialer) (*NetDialService, error) {
	e := &NetDialService{dialer: dialer}
	return e, nil
}

func (e *NetDialService) LoginVerify(ctx context.Context, in *loginv1.LoginVerifyRequest) (*loginv1.LoginVerifyResponse, error) {
	if e.dialer == nil {
		return nil, fmt.Errorf("dialer is nil")
	}
	return e.dialer.LoginVerify(ctx, in)
}
