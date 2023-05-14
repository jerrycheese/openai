package limiter

import (
	"sync"
	"sync/atomic"
	"time"
)

type CallLimit struct {
	counts  sync.Map  // 记录每个UID的API调用次数
	limit   int32     // 每个UID每天的API调用次数限制，0表示无限次
	resetAt time.Time // 每天0点重置计数器
}

func NewCallLimit(limit int32) *CallLimit {
	cl := &CallLimit{
		counts:  sync.Map{},
		limit:   limit,
		resetAt: nextDayZeroTime(time.Now()),
	}
	// 启动计时器，以每天0点为触发点，重置所有UID的调用次数记录
	go cl.resetCounts()
	return cl
}

func (cl *CallLimit) Add(uid string) bool {
	if cl.limit <= 0 {
		return true
	}

	var initValue int32
	count, ok := cl.counts.LoadOrStore(uid, &initValue)
	if !ok {
		count = &initValue
	}
	countV := count.(*int32)
	if atomic.AddInt32(countV, 1) > cl.limit {
		return false // 超过每天限制调用次数，返回false
	}
	return true
}

func (cl *CallLimit) resetCounts() {
	for {
		durToWait := cl.resetAt.Sub(time.Now())
		if durToWait > 0 {
			tm := time.NewTimer(durToWait)
			<-tm.C
		}
		cl.counts = sync.Map{} // 重置所有UID的调用次数记录
		cl.resetAt = nextDayZeroTime(cl.resetAt)
	}
}

func nextDayZeroTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, 1)
}
