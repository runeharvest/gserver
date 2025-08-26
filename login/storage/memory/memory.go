package memory

import (
	"sync"

	entityv1 "github.com/runeharvest/gserver/proto/rh/entity/v1"
)

type MemoryStorage struct {
	mux    sync.RWMutex
	shards map[int32]*entityv1.Shard
	users  map[int32]*entityv1.User
}

func NewMemoryStorage() (*MemoryStorage, error) {
	e := &MemoryStorage{
		shards: make(map[int32]*entityv1.Shard),
		users:  make(map[int32]*entityv1.User),
	}
	return e, nil
}
