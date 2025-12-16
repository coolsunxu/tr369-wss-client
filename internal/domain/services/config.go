// Package services 定义领域层的服务接口
package services

// ConfigProvider 定义配置提供者接口
// 该接口定义了获取配置信息的方法，具体实现由基础设施层提供
type ConfigProvider interface {
	// GetServerURL 获取服务器 URL
	GetServerURL() string

	// GetEndpointID 获取端点 ID
	GetEndpointID() string

	// GetControllerEndpointID 获取控制器端点 ID
	GetControllerEndpointID() string

	// GetMessageChannelSize 获取消息通道大小
	GetMessageChannelSize() int

	// GetDataRefreshInterval 获取数据刷新间隔（秒）
	GetDataRefreshInterval() int

	// GetWriteCountThreshold 获取写入计数阈值
	GetWriteCountThreshold() int

	// GetTR181DataModelPath 获取 TR181 数据模型路径
	GetTR181DataModelPath() string
}
