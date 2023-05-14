package limiter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallLimit_Add(t *testing.T) {
	// 创建一个每个UID每天最多调用10次API的限流器，并重置调用次数的时间间隔为每天0点
	cl := NewCallLimit(1)

	// 模拟用户调用API的过程
	for i := 0; i < 15; i++ {
		uid := fmt.Sprintf("user%d", i%5)
		suc := cl.Add(uid)
		if suc {
			t.Logf("%s API called successfully.\n", uid)
		} else {
			t.Logf("%s API called failed due to rate limit.\n", uid)
		}
		if i < 5 {
			assert.Equal(t, suc, true)
		} else {
			assert.Equal(t, suc, false)
		}
	}
}
