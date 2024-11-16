package tracex

/**
基于openObserver otel的链路追踪
文档地址 https://opentelemetry.io/zh/docs/languages/go/getting-started/
*/

type Config struct {
	EndpointUrl string      `yaml:"EndpointUrl" validate:"required"` // 链路追踪地址
	Auth        string      `yaml:"Auth" validate:"required"`        // 链路追踪认证
	ServerName  string      `yaml:"ServerName" validate:"required"`  // 服务名称
	StreamName  string      `yaml:"StreamName" default:"default"`
	SampleType  SamplerType `yaml:"SampleType" default:"0"` //0~3
	IsDebug     bool        `yaml:"IsDebug" default:"false"`
}

type SamplerType int

const (
	NeverSample SamplerType = iota
	AlwaysSample
	ParentBasedAlwaysSample
	ParentBasedNeverSample
	TraceIdRatioBased
)
