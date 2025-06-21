package cache

import "gocachex/utils"

type ShardedCache struct {
    shards map[string]*LRUCache
    ring   *utils.HashRing
}

func NewShardedCache(nodes []string, cap int) *ShardedCache {
    shards := make(map[string]*LRUCache)
    for _, node := range nodes {
        shards[node] = NewLRU(cap)
    }
    return &ShardedCache{
        shards: shards,
        ring:   utils.NewHashRing(nodes, 3),
    }
}

func (s *ShardedCache) Set(key, value string) {
    node := s.ring.GetNode(key)
    s.shards[node].Set(key, value)
}

func (s *ShardedCache) Get(key string) (string, bool) {
    node := s.ring.GetNode(key)
    return s.shards[node].Get(key)
}

func (s *ShardedCache) Clear() {
    for _, shard := range s.shards {
        shard.Clear()
    }
}
