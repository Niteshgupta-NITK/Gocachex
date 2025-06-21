package cache

import (
    "bufio"
    "fmt"
    "os"
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

    logLine := fmt.Sprintf("SET %s %s\n", key, value)
    _, err := w.writer.WriteString(logLine)
    if err != nil {
        return err
    }
    return w.writer.Flush()
}

func (w *WAL) Close() error {
    return w.file.Close()
}
