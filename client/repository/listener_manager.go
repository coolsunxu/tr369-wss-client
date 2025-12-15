package repository

import (
	tr181Model "tr369-wss-client/tr181/model"
)

// ListenerManager 实现 model.ListenerManager 接口
type ListenerManager struct {
	*BaseRepository
}

// NewListenerManager 创建 ListenerManager 实例
func NewListenerManager(base *BaseRepository) *ListenerManager {
	return &ListenerManager{BaseRepository: base}
}

// AddListener 添加参数变化监听器
func (lm *ListenerManager) AddListener(paramName string, listener tr181Model.Listener) error {
	lm.TR181DataModel.Listeners[paramName] = append(lm.TR181DataModel.Listeners[paramName], listener)
	return nil
}

// RemoveListener 移除指定参数的监听器
func (lm *ListenerManager) RemoveListener(paramName string) error {
	delete(lm.TR181DataModel.Listeners, paramName)
	return nil
}

// ResetListener 重置所有监听器
func (lm *ListenerManager) ResetListener() error {
	lm.TR181DataModel.Listeners = make(map[string][]tr181Model.Listener)
	return nil
}

// NotifyListeners 通知指定参数的所有监听器
func (lm *ListenerManager) NotifyListeners(paramName string, value interface{}) {
	listeners, exists := lm.TR181DataModel.Listeners[paramName]
	if exists {
		for _, listener := range listeners {
			// 在goroutine中执行监听器，避免阻塞
			go listener.Listener(listener.SubscriptionId, value)
		}
	}
}
