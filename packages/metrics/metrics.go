package metrics

import (
	"context"
	"log"
	"os"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"

	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"

	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

var (
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
)

type MeltClientInterface interface {
	Init() error
	CollectMetrics(m MeltClientInterface)
	createMetric(name string, description string, value float64)
}

type MeltTracer struct{}

func (t *MeltTracer) Init() error {
	client := otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(collectorURL),
	)
	ctx := context.Background()
	exp, err := otlpmetric.New(ctx, client)
	if err != nil {
		return err
	}
	defer func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := exp.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	}()
	pusher := controller.New(
		processor.NewFactory(
			simple.NewWithHistogramDistribution(),
			exp,
		),
		controller.WithExporter(exp),
		controller.WithCollectPeriod(2*time.Second),
	)

	global.SetMeterProvider(pusher)

	if err := pusher.Start(ctx); err != nil {
		return err
		// log.Fatalf("could not start metric controller: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := pusher.Stop(ctx); err != nil {
			otel.Handle(err)
		}
	}()
	return nil
}

func (t *MeltTracer) CollectMetrics(m MeltTracer) {
	var ms runtime.MemStats

	m.createMetric("num_cpu", "", float64(runtime.NumCPU()))
	m.createMetric("num_goroutine", "", float64(runtime.NumGoroutine()))
	m.createMetric("num_cgo_call", "", float64(runtime.NumCgoCall()))

	// Memory stats
	m.createMetric("mem_stats.alloc", "", float64(ms.Alloc))
	m.createMetric("mem_stats.total_alloc", "", float64(ms.TotalAlloc))
	m.createMetric("mem_stats.sys", "", float64(ms.Sys))
	m.createMetric("mem_stats.lookups", "", float64(ms.Lookups))
	m.createMetric("mem_stats.mallocs", "", float64(ms.Mallocs))
	m.createMetric("mem_stats.frees", "", float64(ms.Frees))
}

func (t *MeltTracer) createMetric(name string, description string, value float64) {
	ctx := context.Background()
	meter := global.Meter("golang-agent")
	counter, err := meter.SyncFloat64().Counter(
		name,
		instrument.WithDescription(description))
	if err != nil {
		log.Fatalf("Failed to create the instrument: %v", err)
	}
	counter.Add(ctx, value)
}
