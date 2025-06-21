# GoCacheX

# GoCacheX

**GoCacheX** is a high-performance, concurrent, in-memory cache written in Go with:

- 🔁 **LRU eviction**
- ⏳ **TTL-based key expiry**
- 🧵 **Thread-safe access (sync.Mutex, goroutines)**
- 💾 **Write-Ahead Logging (WAL)** for persistence
- ⚙️ **Consistent Hashing for sharding**
- 🌐 **REST APIs** with Gin

---

## ✨ Features

- ⚡ Multi-shard LRU cache with automatic key expiration
- 📦 WAL ensures crash recovery
- 🛡️ Safe for concurrent access via Go’s sync primitives
- 🧠 Sharded via consistent hashing (`crc32`)
- 🧪 Gin APIs: `/set`, `/get/:key`, `/clear`
- 🧱 Extensible architecture for adding metrics, auth, etc.

---

## Architecture
```mermaid
flowchart TD
    A[main.go] --> B[ShardedCache]
    B --> C1[LRUCache: shard1]
    B --> C2[LRUCache: shard2]
    B --> C3[LRUCache: shard3]
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
