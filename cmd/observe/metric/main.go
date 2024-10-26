package main

import (
	"context"
	"github.com/dangweiwu/microkit/observe/metricx"
	"go.opentelemetry.io/otel/metric"
	"log"
	"time"
)

// https://opentelemetry.io/docs/languages/go/instrumentation/#setup
func main() {
	cfg := metricx.Config{
		EndpointUrl: "http://localhost:8080/api/default/v1/metrics",
		Auth:        "Basic cm9vdEBxcS5jb206RDV0RTh2RTdjak9EYkdmYQ==",
		ServerName:  "metricServer",
		StreamName:  "metricStream",
		IsDebug:     false,
		Interval:    1,
	}

	m, err := metricx.NewMetricCli(cfg)
	if err != nil {
		panic(err)
	}
	defer m.Close()
	log.Println("counter指标")
	counter, err := m.Meter.Int64Counter("test_counter",
		metric.WithDescription("Number of API calls."),
		metric.WithUnit("{call}"),
	)
	//counter.Add(context.Background(), 1,
	//	metric.WithAttributes(semconv.HTTPResponseStatusCode(200)),
	//	metric.WithAttributes(attribute.String("test_key", "test_value")),
	//)
	//m.Provider.ForceFlush(context.Background())

	for i := 0; i < 5; i++ {
		counter.Add(context.Background(), 1)
		time.Sleep(time.Second * 2)
	}

	//log.Println("updown指标")
	//updown, err := m.Meter.Int64UpDownCounter("test_updown",
	//	metric.WithDescription("Number of API calls."),
	//	metric.WithUnit("{item}"),
	//)
	//
	//for i := 0; i < 5; i++ {
	//
	//	updown.Add(context.Background(), int64(rand.Intn(11)-5))
	//	time.Sleep(time.Second)
	//}
	//
	//log.Println("gauge指标")
	//
	//gauge, err := m.Meter.Int64Gauge("cpu.fam.speed",
	//	metric.WithDescription("Speed of CPU fan"),
	//	metric.WithUnit("RPM"),
	//)
	//if err != nil {
	//	panic(err)
	//}
	//getCPUFanSpeed := func() int64 {
	//	// Generates a random fan speed for demonstration purpose.
	//	// In real world applications, replace this to get the actual fan speed.
	//	return int64(1500 + rand.Intn(1000))
	//}
	//
	//for i := 0; i < 5; i++ {
	//	gauge.Record(context.Background(), getCPUFanSpeed())
	//	time.Sleep(time.Second)
	//}
	//
	//log.Println("histogram指标")
	//histogram, err := m.Meter.Float64Histogram("task.duration",
	//	metric.WithDescription("the duration of task execution."),
	//	metric.WithUnit("s"),
	//)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//for i := 0; i < 5; i++ {
	//	histogram.Record(context.Background(), float64(rand.Intn(10)))
	//	time.Sleep(time.Second)
	//}

}
