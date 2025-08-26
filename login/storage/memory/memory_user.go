package memory

import (
	"context"

	entityv1 "github.com/runeharvest/gserver/proto/rh/entity/v1"
)

func (e *MemoryStorage) Users(ctx context.Context) ([]*entityv1.User, error) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	var users []*entityv1.User
	for _, user := range e.users {
		users = append(users, user)
	}
	return users, nil
}

func (e *MemoryStorage) UserByLogin(ctx context.Context, login string) (*entityv1.User, error) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	for _, user := range e.users {
		if user.Username == login {
			return user, nil
		}
	}
	return nil, nil
}

func (e *MemoryStorage) UserByUserID(ctx context.Context, userID int32) (*entityv1.User, error) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	if user, ok := e.users[userID]; ok {
		return user, nil
	}
	return nil, nil
}

func (e *MemoryStorage) UsersByState(ctx context.Context, state entityv1.UserState) ([]*entityv1.User, error) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	var users []*entityv1.User
	for _, user := range e.users {
		if user.State == state {
			users = append(users, user)
		}
	}
	return users, nil
}

func (e *MemoryStorage) UserByShardID(ctx context.Context, shardID int32) ([]*entityv1.User, error) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	var users []*entityv1.User
	for _, user := range e.users {
		if user.ShardId == shardID {
			users = append(users, user)
		}
	}
	return users, nil
}

func (e *MemoryStorage) UserCreate(ctx context.Context, user *entityv1.User) (*entityv1.User, error) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.users[user.UserId] = user
	return user, nil
}

func (e *MemoryStorage) UserUpdate(ctx context.Context, user *entityv1.User) error {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.users[user.UserId] = user
	return nil
}
