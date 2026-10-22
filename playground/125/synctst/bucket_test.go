package synctst_test

import (
	"testing"
	"testing/synctest"
	"time"

	"golang.org/x/time/rate"
)

func TestRateLimiterWithSyncTest(t *testing.T) {
	// 将测试逻辑封装在 synctest.Run 中
	// 区域内的所有 Goroutine 和 time 包调用都会绑定到当前的虚拟时钟环境（Bubble）中
	synctest.Test(t, func(t *testing.T) {
		// 创建一个限流器：每秒产生 10 个令牌 (r = 10)，桶容量为 10 (b = 10)
		limiter := rate.NewLimiter(rate.Limit(10), 10)

		// 1. 连续消耗 10 个令牌（瞬间完成）
		for i := range 10 {
			if !limiter.Allow() {
				t.Fatalf("预期第 %d 次请求应该成功，但被限流了", i+1)
			}
		}

		// 2. 第 11 次请求：此时令牌桶已空，预期被拒绝
		if limiter.Allow() {
			t.Fatal("预期第 11 次请求应该被限流，但成功获取到了令牌")
		}

		// 3. 睡眠 1 秒
		// 注意：在 synctest 作用域内，因为当前 Goroutine 进入了 sleep 阻塞状态，
		// 并且没有其他可运行的 Goroutine，虚拟时钟会【瞬间】推进 1 秒，真实测试消耗 0 毫秒！
		time.Sleep(1 * time.Second)

		// 4. 1 秒后令牌桶重新补满，再次校验第 12 次请求（即恢复后的第 1 个令牌）
		if !limiter.Allow() {
			t.Fatal("经过 1 秒的补充后，预期可以成功获取令牌，但依然被限流了")
		}

		t.Log("限流器边界与自动恢复逻辑测试通过！")
	})
}
