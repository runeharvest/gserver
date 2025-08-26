package dial

import (
	loginv1 "github.com/runeharvest/gserver/proto/rh/login/v1"
)

// Dialer defines the interface for network operations.
type Dialer interface {
	loginv1.LoginServiceClient
}
