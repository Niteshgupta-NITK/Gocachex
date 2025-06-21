package main

import (
    "github.com/gin-gonic/gin"
    "gocachex/cache"
    "net/http"
)

var cacheInstance *cache.ShardedCache
var wal *cache.WAL

func main() {
    nodes := []string{"shard1", "shard2", "shard3"}
    cacheInstance = cache.NewShardedCache(nodes, 100)

    var err error
    wal, err = cache.NewWAL("wal.log")
    if err != nil {
        panic(err)
    }
    defer wal.Close()

    // Replay WAL to restore state
    wal.Replay(cacheInstance)

    r := gin.Default()

    r.POST("/set", func(c *gin.Context) {
        var json struct {
            Key   string `json:"key"`
            Value string `json:"value"`
        }
        if err := c.ShouldBindJSON(&json); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        cacheInstance.Set(json.Key, json.Value)
        wal.LogSet(json.Key, json.Value)
        c.JSON(http.StatusOK, gin.H{"status": "success"})
    })

    r.GET("/get/:key", func(c *gin.Context) {
        key := c.Param("key")
        if val, ok := cacheInstance.Get(key); ok {
            c.JSON(http.StatusOK, gin.H{"key": key, "value": val})
        } else {
            c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
        }
    })

    r.POST("/clear", func(c *gin.Context) {
        cacheInstance.Clear()
        c.JSON(http.StatusOK, gin.H{"status": "cache cleared"})
    })

    r.Run(":8080")
}
