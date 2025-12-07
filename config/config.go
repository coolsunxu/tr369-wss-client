package config

// DataRefreshConfig 定义数据刷库的参数
type DataRefreshConfig struct {
	// 刷新间隔时间，单位为秒
	IntervalSeconds int `mapstructure:"interval_seconds"`

	// 写入次数
	WriteCountThreshold int `mapstructure:"write_count_threshold"`

	// tr181节点
	TR181DataModelPath string `mapstructure:"tr181_data_model_path"`
}

type WebsocketConfig struct {
	// 服务器地址
	ServerURL string `mapstructure:"server_url"`

	// controller Id
	ControllerId string `mapstructure:"controller_id"`

	// ping 间隔
	PingInterval int `mapstructure:"ping_interval"`

	// 读取消息大小
	MaxMessageSize int64 `mapstructure:"max_message_size"`

	// 设备信息
	EndpointId string `mapstructure:"endpoint_id"`
	
	// 消息发送通道容量
	MessageChannelSize int `mapstructure:"message_channel_size"`
}

// Config represents the client configuration
type Config struct {
	DataRefreshConfig *DataRefreshConfig `mapstructure:"data_refresh_config"`
	WebsocketConfig   *WebsocketConfig   `mapstructure:"websocket_config"`
}

// GlobalConfig returns the default configuration
var GlobalConfig = Config{
	DataRefreshConfig: &DataRefreshConfig{},
	WebsocketConfig:   &WebsocketConfig{},
}
