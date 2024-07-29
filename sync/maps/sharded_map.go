package maps

import (
	"sync"
)

type ShardsCount uint

type shardedMapItem[K comparable, V any] struct {
	m map[K]V
	sync.RWMutex
}

type ShardedMap[K comparable, V any] struct {
	shards []*shardedMapItem[K, V]
}

func NewShardedMap[K comparable, V any](nshards ShardsCount) *ShardedMap[K, V] {
	shards := make([]*shardedMapItem[K, V], nshards)

	for i := 0; i < int(nshards); i++ {
		shardMap := make(map[K]V)
		shards[i] = &shardedMapItem[K, V]{
			m: shardMap,
		}
	}

	return &ShardedMap[K, V]{
		shards: shards,
	}
}

func (sm *ShardedMap[K, V]) Get(key K) (V, bool) {
	shard := sm.getShard(key)

	shard.RLock()
	defer shard.RUnlock()

	val, ok := shard.m[key]

	return val, ok
}

func (sm *ShardedMap[K, V]) Set(key K, value V) {
	shard := sm.getShard(key)

	shard.Lock()
	defer shard.Unlock()

	shard.m[key] = value
}

func (sm *ShardedMap[K, V]) Delete(key K) {
	shard := sm.getShard(key)

	shard.Lock()
	defer shard.Unlock()

	delete(shard.m, key)
}

func (sm *ShardedMap[K, V]) Len() int64 {
	var total int64 = 0

	for _, shard := range sm.shards {
		total += int64(len(shard.m))
	}

	return int64(total)
}

func (sm *ShardedMap[K, V]) getShardIndex(key K) ShardsCount {
	checkSum := hash(key)

	return ShardsCount(checkSum % uintptr(len(sm.shards)))
}

func (sm *ShardedMap[K, V]) getShard(key K) *shardedMapItem[K, V] {
	index := sm.getShardIndex(key)

	return sm.shards[index]
}
