package deterministic

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spojchil/torchbearing/internal/core"
)

// IDGenerator 按固定前缀生成进程内递增请求 ID。
type IDGenerator struct {
	mu     sync.Mutex
	prefix string
	next   uint64
}

// NewIDGenerator 创建从 1 开始计数的确定性 ID 生成器。
func NewIDGenerator(prefix string) *IDGenerator {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		prefix = "mock-analysis"
	}
	return &IDGenerator{prefix: prefix, next: 1}
}

var _ core.IDGenerator = (*IDGenerator)(nil)

// Next 返回下一个固定格式 ID；互斥锁只保证并发调用不会重复。
func (g *IDGenerator) Next() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	id := fmt.Sprintf("%s-%03d", g.prefix, g.next)
	g.next++
	return id
}
