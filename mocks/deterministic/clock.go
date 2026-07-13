package deterministic

import (
	"time"

	"github.com/spojchil/torchbearing/internal/core"
)

// Clock 始终返回构造时注入的固定时间。
type Clock struct {
	now time.Time
}

// NewClock 创建不读取系统时钟的确定性时钟。
func NewClock(now time.Time) *Clock {
	return &Clock{now: now}
}

var _ core.Clock = (*Clock)(nil)

// Now 返回固定时间，不产生随机漂移。
func (c *Clock) Now() time.Time {
	return c.now
}
