package login

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/runeharvest/gserver/config"
	"github.com/runeharvest/gserver/login/storage"
	entityv1 "github.com/runeharvest/gserver/proto/rh/entity/v1"
	loginv1 "github.com/runeharvest/gserver/proto/rh/login/v1"
	"github.com/runeharvest/gserver/stringfmt"
)

// LoginService is the use case handler for login-related operations.
type LoginService struct {
	loginv1.UnimplementedLoginServiceServer
	storager storage.Storager
}

func NewLoginService(storage storage.Storager) (*LoginService, error) {

	err := configValidate()
	if err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	e := &LoginService{storager: storage}

	return e, nil
}

func (e *LoginService) LoginVerify(ctx context.Context, req *loginv1.LoginVerifyRequest) (*loginv1.LoginVerifyResponse, error) {
	resp := &loginv1.LoginVerifyResponse{}

	isLoginVerboseToClient := config.ValueBool("login", "is_login_verbose_to_client")

	if req.Username == "" {
		resp.Error = "Username is empty"
		return resp, nil
	}
	if len(req.Username) > 16 {
		resp.Error = "Username is too long"
		return resp, nil
	}
	if req.Password == "" {
		resp.Error = "Password is empty"
		return resp, nil
	}
	if len(req.Password) > 16 {
		resp.Error = "Password is too long"
		return resp, nil
	}

	user, err := e.storager.UserByLogin(ctx, req.Username)
	if err != nil {
		resp.Error = "Failed to login for an unknown reason"
		if isLoginVerboseToClient {
			resp.Error = "Failed to create user: " + err.Error()
		}
		return resp, nil
	}

	if user == nil {
		isUnknownUserAllowed := config.ValueBool("login", "is_unknown_user_allowed")
		if !isUnknownUserAllowed {
			resp.Error = "Invalid username or password"
			if isLoginVerboseToClient {
				resp.Error = "User not found and unknown users are not allowed"
			}
			return resp, nil
		}

		isUserCreationAllowed := config.ValueBool("login", "is_user_creation_allowed")
		if !isUserCreationAllowed {
			resp.Error = "Invalid username or password"
			if isLoginVerboseToClient {
				resp.Error = "User not found and user creation is not allowed"
			}
			return resp, nil
		}

		err := stringfmt.UsernameValidate(req.Username)
		if err != nil {
			resp.Error = "Username is invalid: " + err.Error()
			return resp, nil
		}

		err = stringfmt.PasswordValidate(req.Password)
		if err != nil {
			resp.Error = "Password is invalid: " + err.Error()
			return resp, nil
		}

		newUser := &entityv1.User{
			Username: req.Username,
			Password: req.Password,
			State:    entityv1.UserState_OFFLINE,
		}
		user, err = e.storager.UserCreate(ctx, newUser)
		if err != nil {
			resp.Error = "User creation failed for an unknown reason"
			if isLoginVerboseToClient {
				resp.Error = "Failed to create user: " + err.Error()
			}
			return resp, nil
		}
		slog.Info("User created", "username", req.Username, "application", req.Application)
	}

	if user.Password != req.Password {
		resp.Error = "Invalid username or password"
		if isLoginVerboseToClient {
			resp.Error = "Password is incorrect"
		}
		return resp, nil
	}

	if user.State != entityv1.UserState_OFFLINE {
		// TODO: unified service DC broadcast
		//  CMessage msgout("DC");
		// msgout.serial(uid);
		// CUnifiedNetwork::getInstance()->send("WS", msgout);

		// reason = toString("User '%s' is already connected", req.toUtf8().c_str());
		// CMessage vplMsgout("VLP");
		// vplMsgout.serial(reason);
		// netbase.send(vplMsgout, from);
		// return
		resp.Error = fmt.Sprintf("User '%s' is already connected", req.Username)
		return resp, nil
	}

	// TODO: jwt token

	user.State = entityv1.UserState_ONLINE
	err = e.storager.UserUpdate(ctx, user)
	if err != nil {
		resp.Error = "Failed to update user cookie"
		return resp, nil
	}

	shards, err := e.storager.ShardsByClientApplication(ctx, req.Application)
	if err != nil {
		resp.Error = "Failed to get shards"
		return resp, nil
	}

	for _, shard := range shards {
		resp.Shards = append(resp.Shards, &loginv1.LoginVerifyShardResponse{
			Name:        shard.Name,
			PlayerCount: shard.PlayerCount,
			ShardId:     shard.ShardId,
		})
	}

	return resp, nil
}

func configValidate() error {

	requiredKeys := []struct {
		section string
		key     string
		typ     string // "string", "int", "bool"
	}{
		{"login", "displayed_variables", "[]string"},
		{"login", "ws_port", "int"},
		{"login", "web_port", "int"},
		{"login", "client_port", "int"},
		{"login", "is_external_shard_allowed", "bool"},
		{"login", "is_unknown_user_allowed", "bool"},
		{"login", "is_user_creation_allowed", "bool"},
		{"login", "beep", "bool"},
		{"login", "database_host", "string"},
		{"login", "database_name", "string"},
		{"login", "database_username", "string"},
		{"login", "database_password", "string"},
		{"login", "force_database_reconnection", "string"},
		{"login", "is_naming_service_used", "bool"},
		{"login", "is_aes_used", "bool"},
		{"login", "is_login_verbose_to_client", "bool"},
		{"login", "shard_id", "int"},
	}

	for _, k := range requiredKeys {
		var err error
		switch k.typ {
		case "[]string":
			_, err = config.ValueSliceStrE(k.section, k.key)
		case "string":
			_, err = config.ValueStrE(k.section, k.key)
		case "int":
			_, err = config.ValueIntE(k.section, k.key)
		case "bool":
			_, err = config.ValueBoolE(k.section, k.key)
		}
		if err != nil {
			return fmt.Errorf("%s.%s: %w", k.section, k.key, err)
		}
	}
	return nil
}
