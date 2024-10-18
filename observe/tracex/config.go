package tracex

/**
基于openObserver otel的链路追踪
文档地址 https://opentelemetry.io/zh/docs/languages/go/getting-started/
*/

type Config struct {
	EndpointUrl string      `validate:"empty=false"` // 链路追踪地址
	Auth        string      `validate:"empty=false"` // 链路追踪认证
	ServerName  string      `validate:"empty=false"` // 服务名称
	StreamName  string      `default:"default"`
	SampleType  SamplerType `default:"0"` //0~3
	IsDebug     bool        `default:"false"`
}

type SamplerType int

const (
	NeverSample SamplerType = iota
	AlwaysSample
	ParentBasedAlwaysSample
	ParentBasedNeverSample
	TraceIdRatioBased
)
