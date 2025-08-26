package storage

import (
	"context"

	entityv1 "github.com/runeharvest/gserver/proto/rh/entity/v1"
)

type Storager interface {
	Shards(ctx context.Context) ([]*entityv1.Shard, error)
	ShardByShardID(ctx context.Context, shardID int32) (*entityv1.Shard, error)
	ShardByWSAddr(ctx context.Context, wsAddr string) (*entityv1.Shard, error)
	ShardCreate(ctx context.Context, shard *entityv1.Shard) (*entityv1.Shard, error)
	ShardUpdate(ctx context.Context, shard *entityv1.Shard) error
	ShardsByClientApplication(ctx context.Context, clientApp string) ([]*entityv1.Shard, error)

	Users(ctx context.Context) ([]*entityv1.User, error)
	UserByLogin(ctx context.Context, login string) (*entityv1.User, error)
	UserByUserID(ctx context.Context, userID int32) (*entityv1.User, error)
	UsersByState(ctx context.Context, state entityv1.UserState) ([]*entityv1.User, error)
	UserByShardID(ctx context.Context, shardID int32) ([]*entityv1.User, error)
	UserCreate(ctx context.Context, user *entityv1.User) (*entityv1.User, error)
	UserUpdate(ctx context.Context, user *entityv1.User) error
}
