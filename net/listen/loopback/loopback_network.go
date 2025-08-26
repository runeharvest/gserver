package loopback

import (
	"context"
	"fmt"
	"sync"

	loginv1 "github.com/runeharvest/gserver/proto/rh/login/v1"
)

type callbackLookup int

const (
	callbackLoginVerify callbackLookup = iota
)

type LoopbackNetwork struct {
	mutex        sync.RWMutex
	loginService loginv1.LoginServiceServer
}

func NewLoopbackNetwork() (*LoopbackNetwork, error) {
	e := &LoopbackNetwork{}
	return e, nil
}

func (g *LoopbackNetwork) LoginVerify(ctx context.Context, in *loginv1.LoginVerifyRequest) (*loginv1.LoginVerifyResponse, error) {
	if g.loginService == nil {
		return nil, fmt.Errorf("login service not registered")
	}

	return g.loginService.LoginVerify(ctx, in)
}

func (g *LoopbackNetwork) LoginRegister(loginService loginv1.LoginServiceServer) error {

	if g.loginService != nil {
		return fmt.Errorf("login service already registered")
	}
	g.loginService = loginService
	return nil
}
