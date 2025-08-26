package network

import (
	"fmt"

	loginv1 "github.com/runeharvest/gserver/proto/rh/login/v1"

	"github.com/runeharvest/gserver/net/listen"
)

type NetListenService struct {
	listener listen.Listener
}

func NewNetListenService(listener listen.Listener) (*NetListenService, error) {
	e := &NetListenService{listener: listener}
	return e, nil
}

func (e *NetListenService) LoginRegister(loginServer loginv1.LoginServiceServer) error {
	if e.listener == nil {
		return fmt.Errorf("listener is nil")
	}
	return e.listener.LoginRegister(loginServer)
}
