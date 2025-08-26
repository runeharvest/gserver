package memory

import (
	"context"
	"sync"

	entityv1 "github.com/runeharvest/gserver/proto/rh/entity/v1"
)

type MemoryStorage struct {
	mux    sync.RWMutex
	shards map[int32]*entityv1.Shard
}

func NewMemoryStorage() (*MemoryStorage, error) {
	e := &MemoryStorage{
		shards: make(map[int32]*entityv1.Shard),
	}
	return e, nil
}

func (e *MemoryStorage) Shards(ctx context.Context) ([]*entityv1.Shard, error) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	var shards []*entityv1.Shard
	for _, shard := range e.shards {
		shards = append(shards, shard)
	}
	return shards, nil
}

func (e *MemoryStorage) ShardByShardID(ctx context.Context, shardID int32) (*entityv1.Shard, error) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	if shard, ok := e.shards[shardID]; ok {
		return shard, nil
	}
	return nil, nil
}

func (e *MemoryStorage) ShardByWSAddr(ctx context.Context, wsAddr string) (*entityv1.Shard, error) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	for _, shard := range e.shards {
		if shard.WsAddr == wsAddr {
			return shard, nil
		}
	}
	return nil, nil
}

func (e *MemoryStorage) ShardCreate(ctx context.Context, shard *entityv1.Shard) (*entityv1.Shard, error) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.shards[shard.ShardId] = shard
	return shard, nil
}

func (e *MemoryStorage) ShardUpdate(ctx context.Context, shard *entityv1.Shard) error {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.shards[shard.ShardId] = shard
	return nil
}

func (e *MemoryStorage) ShardsByClientApplication(ctx context.Context, clientApp string) ([]*entityv1.Shard, error) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	var shards []*entityv1.Shard
	for _, shard := range e.shards {
		if shard.ClientApp == clientApp {
			shards = append(shards, shard)
		}
	}
	return shards, nil
}
