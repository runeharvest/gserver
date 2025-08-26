package login

import (
	"context"
	"testing"

	"github.com/runeharvest/gserver/login/storage/memory"
	net "github.com/runeharvest/gserver/net"
	netdialloopback "github.com/runeharvest/gserver/net/dial/loopback"
	netlistenloopback "github.com/runeharvest/gserver/net/listen/loopback"
	loginv1 "github.com/runeharvest/gserver/proto/rh/login/v1"
)

func TestDialLogin(t *testing.T) {

	memoryStorage, err := memory.NewMemoryStorage()
	if err != nil {
		t.Fatal("new memory storage:", err)
	}
	loginService, err := NewLoginService(memoryStorage)
	if err != nil {
		t.Fatal("new login service:", err)
	}

	loopbackListen, err := netlistenloopback.NewLoopbackNetwork()
	if err != nil {
		t.Fatal("new loopback listen network:", err)
	}

	netListen, err := net.NewNetListenService(loopbackListen)
	if err != nil {
		t.Fatal("new net listen service:", err)
	}

	netListen.LoginRegister(loginService)

	loopbackDial, err := netdialloopback.NewLoopbackNetwork()
	if err != nil {
		t.Fatal("new loopback dial network:", err)
	}

	loopbackDial.LoginRegister(loginService)

	netDial, err := net.NewNetDialService(loopbackDial)
	if err != nil {
		t.Fatal("new net dial service:", err)
	}

	resp, err := netDial.LoginVerify(context.Background(), &loginv1.LoginVerifyRequest{})
	if err != nil {
		t.Fatal("login verify:", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}

}
