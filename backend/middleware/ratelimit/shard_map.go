package ratelimit

import (
	"hash/fnv"
	"sync"
)

type shard struct {
	mu sync.RWMutex
	m  map[string]*Item
}

type ConcurrentMap struct {
	shards []shard
}

func NewConcurrentMap(shardCount int) *ConcurrentMap {
	shards := make([]shard, shardCount)
	for i := range shards {
		shards[i].m = make(map[string]*Item)
	}
	return &ConcurrentMap{shards: shards}
}

func (cm *ConcurrentMap) getShard(key string) *shard {
	h := fnv32(key)
	return &cm.shards[h%uint32(len(cm.shards))]
}

func (cm *ConcurrentMap) Get(key string) (*Item, bool) {
	shard := cm.getShard(key)
	shard.mu.RLock()
	item, ok := shard.m[key]
	shard.mu.RUnlock()
	return item, ok
}

func (cm *ConcurrentMap) Set(key string, value *Item) {
	shard := cm.getShard(key)
	shard.mu.Lock()
	shard.m[key] = value
	shard.mu.Unlock()
}

func fnv32(key string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	return hash.Sum32()
}
