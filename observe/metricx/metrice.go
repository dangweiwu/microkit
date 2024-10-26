package metricx

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"log"
	"sync"
	"time"
)

type MetricCli struct {
	lock     sync.Mutex
	Provider *sdkmetric.MeterProvider
	export   sdkmetric.Exporter
	resource *resource.Resource
	Meter    metric.Meter
	config   Config
}

func NewMetricCli(config Config) (*MetricCli, error) {
	a := &MetricCli{config: config}
	if err := a.Start(); err != nil {
		return nil, err
	}
	return a, nil
}

func (this *MetricCli) Start() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(this.config.ServerName),
		),
	)
	if err != nil {
		return err
	}
	this.resource = res

	if this.config.IsDebug {
		if export, err := stdoutmetric.New(); err != nil {
			return err
		} else {
			this.export = export
		}
	} else {
		export, err := otlpmetrichttp.New(context.Background(), otlpmetrichttp.WithInsecure(),
			otlpmetrichttp.WithEndpointURL(this.config.EndpointUrl),
			otlpmetrichttp.WithHeaders(map[string]string{
				"Authorization": this.config.Auth,
				"stream-name":   this.config.StreamName,
			}))
		if err != nil {
			log.Printf("Failed to initialize trace exporter: %v", err)
			return err
		} else {
			this.export = export
		}
	}

	this.Provider = sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(this.export,
			// Default is 1m. Set to 3s for demonstrative purposes.
			sdkmetric.WithInterval(time.Duration(this.config.Interval)*time.Second)),
		),
	)

	otel.SetMeterProvider(this.Provider)

	this.Meter = otel.Meter("")
	return nil
}

func (this *MetricCli) Close() {
	this.lock.Lock()
	defer this.lock.Unlock()
	ctx, f := context.WithTimeout(context.Background(), time.Second)
	defer f()
	this.Provider.ForceFlush(ctx)
	ctx, f = context.WithTimeout(context.Background(), time.Second)
	defer f()
	this.Provider.Shutdown(ctx)
}
