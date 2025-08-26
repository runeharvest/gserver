package main

import (
	"fmt"
	"net"
	"os"

	"github.com/runeharvest/gserver/login"
	"github.com/runeharvest/gserver/login/storage/memory"
	netlisten "github.com/runeharvest/gserver/net"
	netlistengrpc "github.com/runeharvest/gserver/net/listen/grpc"
	"google.golang.org/grpc"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}
}

func run() error {
	memoryStorage, err := memory.NewMemoryStorage()
	if err != nil {
		return fmt.Errorf("new memory storage: %w", err)
	}
	loginService, err := login.NewLoginService(memoryStorage)
	if err != nil {
		return fmt.Errorf("new login service: %w", err)
	}

	gs := grpc.NewServer()
	grpcNetwork, err := netlistengrpc.NewGrpcNetwork(gs)
	if err != nil {
		return fmt.Errorf("new grpc network: %w", err)
	}

	netListen, err := netlisten.NewNetListenService(grpcNetwork)
	if err != nil {
		return fmt.Errorf("new network service: %w", err)
	}

	netListen.LoginRegister(loginService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	fmt.Println("Login Server listening on port 50051:")
	err = gs.Serve(lis)
	if err != nil {
		return fmt.Errorf("grpc serve: %w", err)
	}
	return nil
}
