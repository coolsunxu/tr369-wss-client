// Package tr181 定义 TR181 数据模型相关的领域实体
package tr181

// DataModel 表示 TR181 数据模型
// 包含参数映射和监听器管理
type DataModel struct {
	Parameters map[string]interface{}
	Listeners  map[string][]Listener
}

// NewDataModel 创建新的数据模型实例
func NewDataModel() *DataModel {
	return &DataModel{
		Parameters: make(map[string]interface{}),
		Listeners:  make(map[string][]Listener),
	}
}

// Parameter 表示 TR181 参数
type Parameter struct {
	Path     string      // 参数路径
	Value    interface{} // 参数值
	Type     string      // 参数类型
	Writable bool        // 是否可写
}

// NewParameter 创建新的参数实例
func NewParameter(path string, value interface{}, paramType string, writable bool) *Parameter {
	return &Parameter{
		Path:     path,
		Value:    value,
		Type:     paramType,
		Writable: writable,
	}
}
