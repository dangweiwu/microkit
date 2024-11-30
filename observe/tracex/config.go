package tracex

/**
基于openObserver otel的链路追踪
文档地址 https://opentelemetry.io/zh/docs/languages/go/getting-started/
*/

type Config struct {
	EndpointUrl string      `yaml:"endpointUrl" validate:"required"` // 链路追踪地址
	Auth        string      `yaml:"auth" validate:"required"`        // 链路追踪认证
	ServerName  string      `yaml:"serverName" validate:"required"`  // 服务名称
	StreamName  string      `yaml:"streamName" default:"default"`
	SampleType  SamplerType `yaml:"sampleType" default:"0"` //0~3
	IsDebug     bool        `yaml:"isDebug" default:"false"`
}

type SamplerType int

const (
	NeverSample SamplerType = iota
	AlwaysSample
	ParentBasedAlwaysSample
	ParentBasedNeverSample
	TraceIdRatioBased
)
