package grpc

import (
	"context"

	loginv1 "github.com/runeharvest/gserver/proto/rh/login/v1"
	"google.golang.org/grpc"
)

type GrpcNetwork struct {
	server *grpc.Server
}

func NewGrpcNetwork(server *grpc.Server) (*GrpcNetwork, error) {
	e := &GrpcNetwork{server: server}
	return e, nil
}

func (g *GrpcNetwork) LoginVerify(ctx context.Context, in *loginv1.LoginVerifyRequest) (*loginv1.LoginVerifyResponse, error) {

	return &loginv1.LoginVerifyResponse{}, nil
}

func (g *GrpcNetwork) LoginRegister(loginService loginv1.LoginServiceServer) error {
	loginv1.RegisterLoginServiceServer(g.server, loginService)
	return nil
}
