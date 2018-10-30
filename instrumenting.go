package main

import (
	"fmt"
	"httpproject/service"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	next           service.PlayerStore
}

func (mw instrumentingMiddleware) GetPlayerScore(s string) (output int) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetPlayerScore", "error", "false"}
		fmt.Println("INSTRUMENTING")
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		fmt.Printf("Latency? maybe: %v\n", time.Since(begin).Seconds())
		mw.countResult.Observe(float64(output))
	}(time.Now())
	output = mw.next.GetPlayerScore(s)
	return output
}

func (mw instrumentingMiddleware) RecordWin(s string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "RecordWin", "error", "false"}
		fmt.Println("INSTRUMENTING")
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		fmt.Printf("Latency? maybe: %v\n", time.Since(begin).Seconds())
	}(time.Now())

	mw.next.RecordWin(s)
	return
}
