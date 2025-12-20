package logger

import (
	"bytes"
	"encoding/json"
	"sync"
)

type Capture struct {
	mu        sync.Mutex
	entries   []map[string]any
	remainder []byte
}

func (c *Capture) Write(p []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	buf := append(c.remainder, p...)
	lines := bytes.Split(buf, []byte("\n"))

	if len(lines) > 0 {
		c.remainder = lines[len(lines)-1]
		lines = lines[:len(lines)-1]
	}

	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal(line, &m); err == nil {
			c.entries = append(c.entries, m)
		}
	}

	return len(p), nil
}

func (c *Capture) Sync() error { return nil }

func (c *Capture) All() []map[string]any {
	c.mu.Lock()
	defer c.mu.Unlock()

	out := make([]map[string]any, len(c.entries))
	copy(out, c.entries)
	return out
}

func (c *Capture) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = nil
	c.remainder = nil
}
