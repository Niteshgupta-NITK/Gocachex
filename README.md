# GoCacheX

**GoCacheX** is a multithreaded, in-memory key-value caching system built in Go, with support for:

- 🔁 **LRU (Least Recently Used) eviction**
- 🧵 **Thread-safe access using sync primitives**
- 💾 **Write-Ahead Logging (WAL)** for crash-safe persistence
- ⚙️ **Consistent Hashing-based sharding** for horizontal scalability

---

## 🛠️ Features

- 🚀 High performance in-memory storage
- 🧠 LRU eviction using Go’s `container/list`
- 📦 WAL for persistence and crash recovery
- 🧮 Consistent hashing to distribute load across shards
- ⚔️ Thread-safe via mutex locks

---
## Architecture
```mermaid
flowchart TD
    A[main.go] --> B[ShardedCache]
    B --> C1[LRUCache shard1]
    B --> C2[LRUCache shard2]
    B --> C3[LRUCache shard3]
    B --> D[ConsistentHashRing]
    B --> E[WAL Logger]

    subgraph cache/
        C1
        C2
        C3
        E
    end

    subgraph utils/
        D
    end
