package tracex

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	trac "go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"sync"
	"time"
)

// 创建Provicder
func newTraceProvider(exp trace.SpanExporter, res *resource.Resource, sampler SamplerType) *trace.TracerProvider {
	var samplerFunc trace.Sampler
	switch sampler {
	case AlwaysSample:
		samplerFunc = trace.AlwaysSample()
	case NeverSample:
		samplerFunc = trace.NeverSample()
	case ParentBasedAlwaysSample:
		samplerFunc = trace.ParentBased(trace.AlwaysSample())
	case ParentBasedNeverSample:
		samplerFunc = trace.ParentBased(trace.NeverSample())
	case TraceIdRatioBased:
		samplerFunc = trace.TraceIDRatioBased(0.1) // 10% 采样率
	default:
		samplerFunc = trace.AlwaysSample()
	}

	return trace.NewTracerProvider(
		trace.WithSampler(samplerFunc),
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)
}

func newExporter() (trace.SpanExporter, error) {
	return stdouttrace.New()
}

type TraceCli struct {
	lock     sync.Mutex
	export   trace.SpanExporter
	resource *resource.Resource
	Provider *trace.TracerProvider
	Tracer   trac.Tracer
	config   Config
}

func NewTraceCli(config Config) (*TraceCli, error) {
	a := &TraceCli{}
	a.config = config
	if err := a.Start(); err != nil {
		return a, err
	} else {
		return a, nil
	}
}

func (a *TraceCli) Start() error {
	a.lock.Lock()
	defer a.lock.Unlock()
	var export trace.SpanExporter
	var err error
	if a.config.IsDebug {
		export, err = stdouttrace.New()
	} else {
		export, err = otlptracehttp.New(context.Background(),
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpointURL(a.config.EndpointUrl),
			otlptracehttp.WithHeaders(map[string]string{
				"Authorization": a.config.Auth,
				"stream-name":   a.config.StreamName,
			}))
	}
	if err != nil {
		return err
	}
	a.export = export

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(a.config.ServerName),
		),
	)
	if err != nil {
		return err
	}
	a.resource = res
	a.Provider = newTraceProvider(a.export, a.resource, a.config.SampleType)
	otel.SetTracerProvider(a.Provider)
	a.Tracer = a.NewTracer("")
	return nil
}

func (t *TraceCli) Close() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.Provider != nil {
		ctx, c := context.WithTimeout(context.Background(), time.Second)
		defer c()

		t.Provider.ForceFlush(ctx)
		ctx, c = context.WithTimeout(context.Background(), time.Second)
		defer c()

		t.Provider.Shutdown(ctx)
		t.Provider = nil
	}
}

func (t *TraceCli) NewTracer(name string, opt ...trac.TracerOption) trac.Tracer {
	if t.Provider != nil {
		return t.Provider.Tracer(name, opt...)
	}
	return noop.NewTracerProvider().Tracer(name, opt...)
}
