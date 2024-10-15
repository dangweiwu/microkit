package tracex

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
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

type TraceCli struct {
	lock     sync.Mutex
	export   *otlptrace.Exporter
	resource *resource.Resource
	Provider *trace.TracerProvider
	Tracer   trac.Tracer
}

func NewTraceCli(config Config) (*TraceCli, error) {
	a := &TraceCli{}
	export, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpointURL(config.EndpointUrl),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization": config.Auth,
			"stream-name":   config.StreamName,
		}))
	if err != nil {
		return nil, err
	}
	a.export = export

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(config.ServerName),
		),
	)
	if err != nil {
		return nil, err
	}
	a.resource = res
	a.Provider = newTraceProvider(a.export, a.resource, config.SampleType)
	otel.SetTracerProvider(a.Provider)
	a.Tracer = a.NewTracer("")
	return a, nil
}

func (t *TraceCli) ChangeSample(samplerType SamplerType) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	//强制刷新
	if t.Provider != nil {
		t.Provider.ForceFlush(context.Background())
		t.Provider = newTraceProvider(t.export, t.resource, samplerType)
		otel.SetTracerProvider(t.Provider)
		t.Tracer = t.NewTracer("")
	}
	return nil
}

func (t *TraceCli) Close() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.Provider != nil {
		ctx, c := context.WithTimeout(context.Background(), time.Second)
		defer c()

		t.Provider.ForceFlush(ctx)
		t.Provider.Shutdown(ctx)
	}
}

func (t *TraceCli) NewTracer(name string, opt ...trac.TracerOption) trac.Tracer {
	if t.Provider != nil {
		return t.Provider.Tracer(name, opt...)
	}
	return noop.NewTracerProvider().Tracer(name, opt...)
}
