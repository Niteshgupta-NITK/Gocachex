package utils

import (
    "fmt"
    "hash/crc32"
    "sort"
)

type HashRing struct {
    nodes      []string
    ring       map[uint32]string
    sortedKeys []uint32
    replicas   int
}

func NewHashRing(nodes []string, replicas int) *HashRing {
    h := &HashRing{ring: make(map[uint32]string), replicas: replicas}
    for _, node := range nodes {
        for i := 0; i < replicas; i++ {
            hash := crc32.ChecksumIEEE([]byte(fmt.Sprintf("%s%d", node, i)))
            h.ring[hash] = node
            h.sortedKeys = append(h.sortedKeys, hash)
        }
    }
    sort.Slice(h.sortedKeys, func(i, j int) bool { return h.sortedKeys[i] < h.sortedKeys[j] })
    return h
}

func (h *HashRing) GetNode(key string) string {
    hash := crc32.ChecksumIEEE([]byte(key))
    idx := sort.Search(len(h.sortedKeys), func(i int) bool {
        return h.sortedKeys[i] >= hash
    })
    if idx == len(h.sortedKeys) {
        idx = 0
    }
    return h.ring[h.sortedKeys[idx]]
}
