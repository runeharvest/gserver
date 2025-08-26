package loopback

import (
	"context"
	"fmt"
	"sync"

	loginv1 "github.com/runeharvest/gserver/proto/rh/login/v1"
	"google.golang.org/grpc"
)

type LoopbackNetwork struct {
	mutex       sync.RWMutex
	loginServer loginv1.LoginServiceServer
}

func NewLoopbackNetwork() (*LoopbackNetwork, error) {
	e := &LoopbackNetwork{}
	return e, nil
}

func (e *LoopbackNetwork) LoginVerify(ctx context.Context, in *loginv1.LoginVerifyRequest, opts ...grpc.CallOption) (*loginv1.LoginVerifyResponse, error) {
	if e.loginServer == nil {
		return nil, fmt.Errorf("login service not registered")
	}

	return e.loginServer.LoginVerify(ctx, in)
}

func (e *LoopbackNetwork) LoginRegister(loginServer loginv1.LoginServiceServer) error {
	e.loginServer = loginServer
	return nil
}
