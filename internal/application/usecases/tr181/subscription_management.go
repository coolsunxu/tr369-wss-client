// Package tr181 提供 TR181 相关的用例实现
package tr181

import (
	"fmt"

	"tr369-wss-client/internal/domain/entities/tr181"
	"tr369-wss-client/internal/domain/services"
)

// SubscriptionRepository 订阅仓储接口
type SubscriptionRepository interface {
	AddListener(paramName string, listener tr181.Listener) error
	RemoveListener(paramName string) error
	ResetListeners() error
	GetValueByPath(path string) (interface{}, error)
}

// SubscriptionManagementUseCase 订阅管理用例
type SubscriptionManagementUseCase struct {
	repository SubscriptionRepository
	logger     services.Logger
}

// NewSubscriptionManagementUseCase 创建新的订阅管理用例
func NewSubscriptionManagementUseCase(repo SubscriptionRepository, logger services.Logger) *SubscriptionManagementUseCase {
	return &SubscriptionManagementUseCase{
		repository: repo,
		logger:     logger,
	}
}

// HandleSubscription 处理订阅
func (uc *SubscriptionManagementUseCase) HandleSubscription(path, subscriptionId, subscriptionType string, handler tr181.Handler) error {
	uc.logger.Debug("处理订阅, 路径: %s, ID: %s, 类型: %s", path, subscriptionId, subscriptionType)

	listener := tr181.Listener{
		SubscriptionId: subscriptionId,
		Listener:       handler,
	}

	switch subscriptionType {
	case tr181.ValueChange:
		return uc.repository.AddListener(path, listener)
	case tr181.ObjectCreation:
		return uc.repository.AddListener(path, listener)
	case tr181.ObjectDeletion:
		return uc.repository.AddListener(path, listener)
	case tr181.OperationComplete:
		return uc.repository.AddListener(path, listener)
	case tr181.Event:
		return uc.repository.AddListener(path, listener)
	default:
		return fmt.Errorf("未知的订阅类型: %s", subscriptionType)
	}
}

// RemoveSubscription 移除订阅
func (uc *SubscriptionManagementUseCase) RemoveSubscription(path string) error {
	uc.logger.Debug("移除订阅, 路径: %s", path)
	return uc.repository.RemoveListener(path)
}

// ResetAllSubscriptions 重置所有订阅
func (uc *SubscriptionManagementUseCase) ResetAllSubscriptions() error {
	uc.logger.Debug("重置所有订阅")
	return uc.repository.ResetListeners()
}

// HandleAddSubscription 处理添加订阅请求
func (uc *SubscriptionManagementUseCase) HandleAddSubscription(requestPath string, paramSettings map[string]string, handler tr181.Handler) error {
	if requestPath == "" {
		return nil
	}

	if requestPath != tr181.DeviceLocalAgentSubscription {
		return nil
	}

	referenceList := ""
	subscriptionId := ""
	subscriptionType := ""

	for key, value := range paramSettings {
		switch key {
		case "ID":
			subscriptionId = value
		case "ReferenceList":
			referenceList = value
		case "NotifType":
			subscriptionType = value
		}
	}

	return uc.HandleSubscription(referenceList, subscriptionId, subscriptionType, handler)
}
