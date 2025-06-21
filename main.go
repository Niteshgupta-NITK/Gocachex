package main

import (
    "fmt"
    "gocachex/cache"
    "math/rand"
    "sync"
    "time"
)

func main() {
    nodes := []string{"shard1", "shard2", "shard3"}
    shardedCache := cache.NewShardedCache(nodes, 100)

    wal, err := cache.NewWAL("wal.log")
    if err != nil {
        panic(err)
    }
    defer wal.Close()

    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := 0; j < 100; j++ {
                key := fmt.Sprintf("key-%d-%d", id, j)
                value := fmt.Sprintf("value-%d-%d", id, j)
                shardedCache.Set(key, value)

                // Log this write to WAL
                if err := wal.LogSet(key, value); err != nil {
                    fmt.Println("WAL write error:", err)
                }

                time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
            }
        }(i)
    }

    wg.Wait()
    fmt.Println("All goroutines finished. Check wal.log for output.")
}
