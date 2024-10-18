package metricx

type Config struct {
	EndpointUrl string `validate:"empty=false"` // 链路追踪地址
	Auth        string `validate:"empty=false"` // 链路追踪认证
	ServerName  string `validate:"empty=false"` // 服务名称
	StreamName  string `default:"default"`
	Interval    int    `default:"60"` //导出时间间隔 单位秒
}
