package cache

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "sync"
)

type WAL struct {
    mu     sync.Mutex
    file   *os.File
    writer *bufio.Writer
}

func NewWAL(path string) (*WAL, error) {
    f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    return &WAL{file: f, writer: bufio.NewWriter(f)}, nil
}

func (w *WAL) LogSet(key, value string) error {
    w.mu.Lock()
    defer w.mu.Unlock()
    _, err := w.writer.WriteString(fmt.Sprintf("SET %s %s\n", key, value))
    if err != nil {
        return err
    }
    return w.writer.Flush()
}

func (w *WAL) Replay(sc *ShardedCache) error {
    f, err := os.Open(w.file.Name())
    if err != nil {
        return err
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "SET") {
            parts := strings.SplitN(line, " ", 3)
            if len(parts) == 3 {
                sc.Set(parts[1], parts[2])
            }
        }
    }
    return scanner.Err()
}

func (w *WAL) Close() error {
    return w.file.Close()
}
