package init

import (
	"context"
	"github.com/middleware-labs/agent-apm-go/packages/metrics"
	"github.com/middleware-labs/agent-apm-go/packages/tracer"
	"log"
	"time"
)

func Execute() func(context.Context) error {
	cleanup := tracer.InitTracer()
	defer cleanup(context.Background())
	handler := metrics.MeltTracer{}
	err := handler.Init()
	if err != nil {
		log.Fatalf("failed to create the collector exporter: %v", err)
	}
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			handler.CollectMetrics(handler)
		}
	}
	return nil
}
