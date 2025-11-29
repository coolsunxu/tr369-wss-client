package datamodel

import (
	"sync"
	"time"
	"tr369-wss-client/common"

	"tr369-wss-client/protocol"
)

// TR181DataModel represents the TR181 data model
type TR181DataModel struct {
	parameters map[string]interface{}
	lock       sync.RWMutex
	listeners  map[string][]ParameterChangeListener
	startTime  time.Time
}

// ParameterChangeListener defines a callback for parameter changes
type ParameterChangeListener func(name string, oldValue, newValue interface{})

// NewTR181DataModel creates a new TR181 data model instance with default values
func NewTR181DataModel() *TR181DataModel {
	model := &TR181DataModel{
		parameters: make(map[string]interface{}),
		listeners:  make(map[string][]ParameterChangeListener),
		startTime:  time.Now(),
	}

	// 初始化默认参数值
	model.initializeDefaultValues()

	return model
}

// initializeDefaultValues initializes the data model with default TR181 values
func (m *TR181DataModel) initializeDefaultValues() {
	// 加载TR181默认节点值
	m.loadDefaultTR181Nodes()
}

// GetParameters retrieves the value of a parameter
func (m *TR181DataModel) GetParameters() map[string]interface{} {
	return m.parameters
}

// AddParameterChangeListener adds a listener for parameter changes
func (m *TR181DataModel) AddParameterChangeListener(paramName string, listener ParameterChangeListener) error {
	if err := protocol.ValidateParameterName(paramName); err != nil {
		return err
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	m.listeners[paramName] = append(m.listeners[paramName], listener)
	return nil
}

// RemoveParameterChangeListener removes a listener for parameter changes
func (m *TR181DataModel) RemoveParameterChangeListener(paramName string, listener ParameterChangeListener) error {
	if err := protocol.ValidateParameterName(paramName); err != nil {
		return err
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	// 在Go中无法直接比较函数，这里我们提供一个简单的移除机制
	// 注意：这个实现会移除所有监听器，因为无法精确比较函数
	delete(m.listeners, paramName)

	return nil
}

// notifyListeners notifies all listeners of a parameter change
func (m *TR181DataModel) notifyListeners(name string, oldValue, newValue interface{}) {
	m.lock.RLock()
	listeners, exists := m.listeners[name]
	m.lock.RUnlock()

	if exists {
		for _, listener := range listeners {
			// 在goroutine中执行监听器，避免阻塞
			go listener(name, oldValue, newValue)
		}
	}

	// 通知通配符监听器（监听所有参数变化）
	m.lock.RLock()
	wildcardListeners, exists := m.listeners["*"]
	m.lock.RUnlock()

	if exists {
		for _, listener := range wildcardListeners {
			go listener(name, oldValue, newValue)
		}
	}
}

// loadDefaultTR181Nodes loads default TR181 node values from JSON file
func (m *TR181DataModel) loadDefaultTR181Nodes() {
	// 读取默认TR181节点配置文件
	filePath := "datamodel/default_tr181_nodes.json"

	m.parameters = common.LoadJsonFile(filePath)
}
