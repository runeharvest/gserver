package login

import (
	"context"
	"testing"

	"github.com/runeharvest/gserver/config"
	"github.com/runeharvest/gserver/login/storage/memory"
	net "github.com/runeharvest/gserver/net"
	netdialloopback "github.com/runeharvest/gserver/net/dial/loopback"
	netlistenloopback "github.com/runeharvest/gserver/net/listen/loopback"
	loginv1 "github.com/runeharvest/gserver/proto/rh/login/v1"
)

func TestDialLogin(t *testing.T) {

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

	resp, err := netDial.LoginVerify(context.Background(), &loginv1.LoginVerifyRequest{
		Username: "testuser",
		Password: "testpassword",
	})
	if err != nil {
		t.Fatal("login verify:", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
	if resp.Error != "" {
		t.Fatal("response error:", resp.Error)
	}

}

func defaultLoginConfig() map[string]any {
	return map[string]any{
		"login": map[string]any{
			"displayed_variables":         []string{},
			"ws_port":                     8080,
			"web_port":                    8081,
			"client_port":                 8082,
			"is_external_shard_allowed":   true,
			"is_unknown_user_allowed":     true,
			"is_user_creation_allowed":    true,
			"beep":                        true,
			"database_host":               "localhost",
			"database_name":               "login",
			"database_username":           "user",
			"database_password":           "password",
			"force_database_reconnection": "5s",
			"is_naming_service_used":      false,
			"is_aes_used":                 false,
			"shard_id":                    1,
			"is_login_verbose_to_client":  true,
		},
	}
}
