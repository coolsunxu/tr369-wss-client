# 设计文档

## 概述

本设计文档描述了将 `client/usecase/client.go` 文件按职责拆分为多个小文件的方案。当前文件约 500 行代码，包含消息处理、订阅管理、通知发送等多个职责。拆分后将提高代码的可读性、可维护性和团队协作效率。

## 架构

拆分后的文件结构：

```
client/usecase/
├── client.go              # 核心结构定义、构造函数、消息分发入口
├── message_handler.go     # USP 消息处理（GET/SET/ADD/DELETE/OPERATE）
├── subscription.go        # 订阅管理（添加/删除订阅）
├── notification.go        # 通知发送（ValueChange/ObjectCreation/ObjectDeletion）
└── helpers.go             # 辅助函数（参数提取等）
```

架构图：

```
┌─────────────────────────────────────────────────────────────────┐
│                         client.go                                │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  ClientUseCase 结构体定义                                │    │
│  │  NewClientUseCase() 构造函数                             │    │
│  │  HandleMessage() 消息分发入口                            │    │
│  │  HandleMTPMsgTransmit() 消息传输                         │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                              │
          ┌───────────────────┼───────────────────┐
          ▼                   ▼                   ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│message_handler.go│  │ subscription.go │  │ notification.go │
├─────────────────┤  ├─────────────────┤  ├─────────────────┤
│HandleGetRequest │  │HandleAddLocal.. │  │HandleValueChange│
│HandleSetRequest │  │HandleDeleteLo.. │  │HandleObjCreation│
│HandleAddRequest │  │HandleSubscript..│  │HandleObjDeletion│
│HandleDeleteReq..│  │isSubscription.. │  │sendNotification │
│HandleOperateReq │  │extractSubscri.. │  │notifyValueChange│
│HandleNotifyResp │  │deleteAllSubs.. │  │notifyObjCreation│
└─────────────────┘  └─────────────────┘  └─────────────────┘
          │                   │                   │
          └───────────────────┼───────────────────┘
                              ▼
                    ┌─────────────────┐
                    │   helpers.go    │
                    ├─────────────────┤
                    │extractParamSet..│
                    │constructGetResp │
                    │isExistPath      │
                    │getNewInstance   │
                    └─────────────────┘
```

## 组件和接口

### 1. client.go - 核心模块

保留核心结构定义和入口方法：

```go
// 结构体定义
type ClientUseCase struct {
    Config         *config.Config
    DataRepo       model.DataRepository
    ListenerMgr    model.ListenerManager
    ctx            context.Context
    messageChannel chan []byte
}

// 构造函数
func NewClientUseCase(...) *ClientUseCase

// 消息分发入口
func (uc *ClientUseCase) HandleMessage(msg *api.Msg)

// 消息传输
func (uc *ClientUseCase) HandleMTPMsgTransmit(msg *api.Msg) error
```

### 2. message_handler.go - 消息处理模块

处理各类 USP 消息请求：

```go
// GET 请求处理
func (uc *ClientUseCase) HandleGetRequest(msg *api.Msg)

// SET 请求处理
func (uc *ClientUseCase) HandleSetRequest(msg *api.Msg)

// ADD 请求处理
func (uc *ClientUseCase) HandleAddRequest(msg *api.Msg)

// DELETE 请求处理
func (uc *ClientUseCase) HandleDeleteRequest(msg *api.Msg)

// OPERATE 请求处理
func (uc *ClientUseCase) HandleOperateRequest(msg *api.Msg)

// NOTIFY_RESP 处理
func (uc *ClientUseCase) HandleNotifyResp(msg *api.Msg)

// OPERATE_COMPLETE 通知发送
func (uc *ClientUseCase) SendOperateCompleteNotify(...)

// 副作用处理
func (uc *ClientUseCase) handleObjectCreationSideEffects(...)
func (uc *ClientUseCase) handleObjectDeletionPreEffects(...)
```

### 3. subscription.go - 订阅管理模块

集中管理订阅相关逻辑：

```go
// 常量定义（从 model 包引用）
var subscriptionInstanceRegex = regexp.MustCompile(...)

// 添加订阅
func (uc *ClientUseCase) HandleAddLocalAgentSubscription(requestPath string, paramSettings map[string]string) error

// 删除订阅
func (uc *ClientUseCase) HandleDeleteLocalAgentSubscription(requestPath string) error

// 订阅处理
func (uc *ClientUseCase) HandleSubscription(path, subscriptionId, subscriptionType string) error

// 辅助函数
func (uc *ClientUseCase) isSubscriptionPath(path string) bool
func (uc *ClientUseCase) isSubscriptionParentPath(path string) bool
func (uc *ClientUseCase) isSubscriptionInstancePath(path string) bool
func (uc *ClientUseCase) extractSubscriptionParams(paramSettings map[string]string) (*model.SubscriptionParams, error)
func (uc *ClientUseCase) deleteAllSubscriptions(parentPath string) error
func (uc *ClientUseCase) deleteSingleSubscription(instancePath string) error
func (uc *ClientUseCase) getSubscriptionReferenceList(instancePath string) (string, error)
```

### 4. notification.go - 通知发送模块

统一管理通知发送逻辑：

```go
// 通用通知发送
func (uc *ClientUseCase) sendNotification(subscriptionId string, notify *api.Notify, notifyType string)

// 事件处理器
func (uc *ClientUseCase) HandleValueChange(subscriptionId string, change interface{})
func (uc *ClientUseCase) HandleObjectCreation(subscriptionId string, change interface{})
func (uc *ClientUseCase) HandleObjectDeletion(subscriptionId string, change interface{})

// 通知触发函数
func (uc *ClientUseCase) notifyValueChange(paramPath string, newValue string)
func (uc *ClientUseCase) notifyObjectCreation(path string, uniqueKeys map[string]string)
func (uc *ClientUseCase) notifyObjectDeletion(objPath string)
```

### 5. helpers.go - 辅助函数模块

通用辅助函数：

```go
// 参数提取
func extractParamSettings[T model.ParamSetting](settings []T) map[string]string

// 数据访问辅助
func (uc *ClientUseCase) constructGetResp(paths []string) api.Response_GetResp
func (uc *ClientUseCase) isExistPath(path string) (isSuccess bool, nodePath string)
func (uc *ClientUseCase) getNewInstance(path string) string
```

## 数据模型

数据模型保持不变，继续使用现有的 `ClientUseCase` 结构体：

```go
type ClientUseCase struct {
    Config         *config.Config
    DataRepo       model.DataRepository
    ListenerMgr    model.ListenerManager
    ctx            context.Context
    messageChannel chan []byte
}
```

## 正确性属性

*属性是一种特征或行为，应该在系统的所有有效执行中保持为真——本质上是关于系统应该做什么的形式化陈述。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### Property 1: 拆分前后行为等价性

*对于任何* 有效的 USP 消息输入，拆分后的代码应产生与拆分前完全相同的输出和副作用。

**验证: 需求 4.1, 4.2**

## 错误处理

错误处理策略保持不变：

1. **防御性检查**: 所有公共方法入口处检查 nil 参数
2. **错误包装**: 使用 `fmt.Errorf` 添加上下文信息
3. **优雅降级**: 非致命错误记录日志后继续执行
4. **错误传播**: 致命错误向上层传播

## 测试策略

### 单元测试

- 确保拆分后所有现有测试继续通过
- 验证各模块的独立功能

### 属性测试

使用 Go 的 `testing/quick` 包进行属性测试：

- 配置每个属性测试运行至少 100 次迭代
- 每个属性测试需标注对应的正确性属性：`**Feature: client-usecase-split, Property {number}: {property_text}**`

### 编译验证

- 拆分后代码必须编译通过
- 无新增编译警告
- 所有导入正确解析

## 拆分原则

1. **单一职责**: 每个文件只负责一类功能
2. **最小改动**: 只移动代码，不修改逻辑
3. **保持兼容**: 所有公共方法签名不变
4. **清晰命名**: 文件名直接反映其职责

