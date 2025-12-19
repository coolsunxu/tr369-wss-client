package repository

import (
	"sort"

	tr181Model "tr369-wss-client/tr181/model"
)

// matchedListener 匹配的监听器信息
type matchedListener struct {
	listener  tr181Model.Listener
	matchType MatchType
}

// ListenerManager 实现 model.ListenerManager 接口
type ListenerManager struct {
	*BaseRepository
	validator *PathValidator // 路径校验器
	matcher   *PathMatcher   // 路径匹配器
}

// NewListenerManager 创建 ListenerManager 实例
func NewListenerManager(base *BaseRepository) *ListenerManager {
	return &ListenerManager{
		BaseRepository: base,
		validator:      NewPathValidator(),
		matcher:        NewPathMatcher(),
	}
}

// AddListener 添加参数变化监听器
// 在添加前校验路径是否符合 TR181 Path Name 规范
func (lm *ListenerManager) AddListener(paramName string, listener tr181Model.Listener) error {
	// 校验路径
	if err := lm.validator.ValidatePath(paramName); err != nil {
		return err
	}

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

// NotifyListeners 通知匹配参数路径的所有监听器
// 支持层级前缀匹配和通配符匹配
// 匹配优先级：精确匹配 > 前缀匹配 > 通配符匹配
func (lm *ListenerManager) NotifyListeners(paramName string, value interface{}) {
	var matchedListeners []matchedListener

	// 遍历所有订阅，找出匹配的监听器
	for subPath, listeners := range lm.TR181DataModel.Listeners {
		result := lm.matcher.Match(subPath, paramName)
		if result.Matched {
			for _, listener := range listeners {
				matchedListeners = append(matchedListeners, matchedListener{
					listener:  listener,
					matchType: result.MatchType,
				})
			}
		}
	}

	// 如果没有匹配的监听器，直接返回
	if len(matchedListeners) == 0 {
		return
	}

	// 按匹配类型排序：精确匹配(1) > 前缀匹配(2) > 通配符匹配(3)
	sort.Slice(matchedListeners, func(i, j int) bool {
		return matchedListeners[i].matchType < matchedListeners[j].matchType
	})

	// 触发所有匹配的监听器
	for _, ml := range matchedListeners {
		// 在 goroutine 中执行监听器，避免阻塞
		go ml.listener.Listener(ml.listener.SubscriptionId, value)
	}
}
