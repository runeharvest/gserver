package grpc

import (
	"context"

	loginv1 "github.com/runeharvest/gserver/proto/rh/login/v1"
	"google.golang.org/grpc"
)

type GrpcNetwork struct {
	conn        *grpc.ClientConn
	loginClient loginv1.LoginServiceClient
}

func NewGrpcNetwork(conn *grpc.ClientConn) (*GrpcNetwork, error) {
	e := &GrpcNetwork{conn: conn}
	return e, nil
}

func (e *GrpcNetwork) LoginVerify(ctx context.Context, in *loginv1.LoginVerifyRequest, opts ...grpc.CallOption) (*loginv1.LoginVerifyResponse, error) {
	if e.loginClient == nil {
		e.loginClient = loginv1.NewLoginServiceClient(e.conn)
	}

	return e.loginClient.LoginVerify(ctx, in)
}

func (e *GrpcNetwork) LoginRegister(loginClient loginv1.LoginServiceClient) error {
	e.loginClient = loginClient
	return nil
}
