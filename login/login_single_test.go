package login

import (
	"context"
	"net"
	"testing"

	"github.com/runeharvest/gserver/config"
	"github.com/runeharvest/gserver/login/storage/memory"
	netlisten "github.com/runeharvest/gserver/net"
	netdialgrpc "github.com/runeharvest/gserver/net/dial/grpc"
	netlistengrpc "github.com/runeharvest/gserver/net/listen/grpc"
	loginv1 "github.com/runeharvest/gserver/proto/rh/login/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestSingleDialLogin(t *testing.T) {
	err := config.SetConfig(defaultLoginConfig())
	if err != nil {
		t.Fatal("set config:", err)
	}
	memoryStorage, err := memory.NewMemoryStorage()
	if err != nil {
		t.Fatal("new memory storage:", err)
	}
	loginService, err := NewLoginService(memoryStorage)
	if err != nil {
		t.Fatal("new login service:", err)
	}
	err = config.SetConfig(defaultLoginConfig())
	if err != nil {
		t.Fatal("set config:", err)
	}
	gs := grpc.NewServer()

	grpcListen, err := netlistengrpc.NewGrpcNetwork(gs)
	if err != nil {
		t.Fatal("new grpc listen network:", err)
	}

	netListen, err := netlisten.NewNetListenService(grpcListen)
	if err != nil {
		t.Fatal("new net listen service:", err)
	}
	netListen.LoginRegister(loginService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		t.Fatal("net listen:", err)
	}

	go gs.Serve(lis)

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(nil))
	if err != nil {
		t.Fatal("new grpc client:", err)
	}
	defer conn.Close()

	grpcDial, err := netdialgrpc.NewGrpcNetwork(conn)
	if err != nil {
		t.Fatal("new grpc dial network:", err)
	}

	err = grpcDial.LoginRegister(loginv1.NewLoginServiceClient(conn))
	if err != nil {
		t.Fatal("login register:", err)
	}

	netDial, err := netlisten.NewNetDialService(grpcDial)
	if err != nil {
		t.Fatal("new net dial service:", err)
	}

	resp, err := netDial.LoginVerify(context.Background(), &loginv1.LoginVerifyRequest{})
	if err != nil {
		t.Fatal("login verify:", err)
	}
	if resp == nil {
		t.Fatal("login verify: response is nil")
	}

}
