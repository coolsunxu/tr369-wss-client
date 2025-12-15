# 设计文档

## 概述

本设计文档描述了对 `client/usecase/client.go` 文件的优化方案。优化目标是提高代码的可读性、可维护性、健壮性和性能，同时保持现有功能的完整性。

## 架构

当前架构采用分层设计：
- **UseCase 层**: 处理业务逻辑（本次优化目标）
- **Repository 层**: 数据访问和持久化
- **Model 层**: 数据模型定义

优化后的架构保持不变，但在 UseCase 层内部进行重构：

```
┌─────────────────────────────────────────────────────────┐
│                    ClientUseCase                         │
├─────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐              │
│  │ Message Handler │  │ Response Helper │              │
│  │   (消息处理)     │  │   (响应辅助)     │              │
│  └────────┬────────┘  └────────┬────────┘              │
│           │                    │                        │
│  ┌────────▼────────────────────▼────────┐              │
│  │         Subscription Manager          │              │
│  │           (订阅管理器)                 │              │
│  └──────────────────────────────────────┘              │
│                                                         │
│  ┌──────────────────────────────────────┐              │
│  │         Notification Sender           │              │
│  │           (通知发送器)                 │              │
│  └──────────────────────────────────────┘              │
└─────────────────────────────────────────────────────────┘
```

## 组件和接口

### 1. 消息处理组件

负责处理不同类型的 USP 消息：

```go
// 消息处理入口
func (uc *ClientUseCase) HandleMessage(msg *api.Msg)

// 各类型消息处理器
func (uc *ClientUseCase) HandleGetRequest(msg *api.Msg)
func (uc *ClientUseCase) HandleSetRequest(msg *api.Msg)
func (uc *ClientUseCase) HandleAddRequest(msg *api.Msg)
func (uc *ClientUseCase) HandleDeleteRequest(msg *api.Msg)
func (uc *ClientUseCase) HandleOperateRequest(msg *api.Msg)
```

### 2. 响应辅助组件（新增）

统一响应发送逻辑：

```go
// 统一的响应发送函数
func (uc *ClientUseCase) sendResponse(msg *api.Msg, operation string) error

// 参数设置提取辅助函数
func extractParamSettings(settings []*api.Set_UpdateParamSetting) map[string]string
```

### 3. 订阅管理组件

集中管理订阅相关逻辑：

```go
// 订阅处理
func (uc *ClientUseCase) HandleSubscription(path, subscriptionId, subscriptionType string) error
func (uc *ClientUseCase) HandleAddLocalAgentSubscription(requestPath string, paramSettings map[string]string) error
func (uc *ClientUseCase) HandleDeleteLocalAgentSubscription(requestPath string) error
```

### 4. 通知发送组件

统一通知发送逻辑：

```go
// 通用通知发送函数
func (uc *ClientUseCase) sendNotification(subscriptionId string, notification interface{}, notifyType string)

// 具体通知处理器
func (uc *ClientUseCase) HandleValueChange(subscriptionId string, change interface{})
func (uc *ClientUseCase) HandleObjectCreation(subscriptionId string, change interface{})
func (uc *ClientUseCase) HandleObjectDeletion(subscriptionId string, change interface{})
```

## 数据模型

现有数据模型保持不变：

```go
type ClientUseCase struct {
    Config           *config.Config
    ClientRepository model.ClientRepository
    ctx              context.Context
    messageChannel   chan []byte
}
```

## 正确性属性

*属性是一种特征或行为，应该在系统的所有有效执行中保持为真——本质上是关于系统应该做什么的形式化陈述。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### Property 1: 错误处理完整性

*对于任何* 无效输入（nil 消息、空路径、无效参数），错误处理函数应返回包含上下文信息的错误，而不是 panic 或返回空错误。

**验证: 需求 1.2, 1.3, 1.4**

### Property 2: 日志格式一致性

*对于任何* 日志记录操作，日志消息应包含操作类型标识，并且相同类型的操作应使用相同的日志格式模板。

**验证: 需求 4.1, 4.2, 4.3**

## 错误处理

### 错误处理策略

1. **防御性检查**: 所有公共方法入口处检查 nil 参数
2. **错误包装**: 使用 `fmt.Errorf` 添加上下文信息
3. **优雅降级**: 非致命错误记录日志后继续执行
4. **错误传播**: 致命错误向上层传播

### 错误类型

```go
// 常见错误场景
- 空消息错误: "received nil message"
- 路径不存在: "path not found: %s"
- 订阅类型未知: "unknown subscription type: %s"
- 消息发送失败: "failed to send %s response: %v"
```

## 测试策略

### 单元测试

- 测试各消息处理函数的正常流程
- 测试边界条件（空输入、无效参数）
- 测试错误处理路径

### 属性测试

使用 Go 的 `testing/quick` 包进行属性测试：

- 配置每个属性测试运行至少 100 次迭代
- 每个属性测试需标注对应的正确性属性：`**Feature: client-usecase-optimization, Property {number}: {property_text}**`

### 测试覆盖范围

1. **HandleMessage**: 测试所有消息类型分发
2. **错误处理**: 测试 nil 输入和无效参数
3. **订阅管理**: 测试添加、删除订阅的各种场景
4. **通知发送**: 测试各类型通知的发送逻辑

## 优化详细设计

### 1. 错误处理增强

**当前问题:**
- 部分函数缺少 nil 检查
- 错误信息缺少上下文

**优化方案:**
```go
func (uc *ClientUseCase) HandleMessage(msg *api.Msg) {
    if msg == nil {
        logger.Warnf("received nil message, ignoring")
        return
    }
    if msg.Header == nil {
        logger.Warnf("received message with nil header, ignoring")
        return
    }
    // ... 原有逻辑
}
```

### 2. 代码重复消除

**当前问题:**
- HandleValueChange、HandleObjectCreation、HandleObjectDeletion 有大量重复代码

**优化方案:**
```go
// 通用通知发送函数
func (uc *ClientUseCase) sendNotification(subscriptionId string, notification isNotify_Notification, notifyType string) {
    notify := &api.Notify{
        SubscriptionId: subscriptionId,
        SendResp:       true,
        Notification:   notification,
    }
    msg := utils.CreateNotifyMessage(notify)
    
    if err := uc.HandleMTPMsgTransmit(msg); err != nil {
        logger.Warnf("send %s notify error: %v", notifyType, err)
        return
    }
    logger.Infof("client sent %s notify msg %s", notifyType, msg)
}
```

### 3. 参数提取辅助函数

**当前问题:**
- HandleSetRequest 和 HandleAddRequest 中参数提取逻辑重复

**优化方案:**
```go
// extractParamSettings 从参数设置列表中提取键值对
func extractParamSettings(settings interface{}) map[string]string {
    result := make(map[string]string)
    // 使用类型断言处理不同类型的设置
    // ...
    return result
}
```

### 4. 日志格式统一

**当前问题:**
- 日志格式不一致，有的用 "client receive"，有的用 "client sent"

**优化方案:**
```go
const (
    logPrefixReceive = "[USP] receive %s request: %s"
    logPrefixSend    = "[USP] send %s response: %s"
    logPrefixError   = "[USP] %s error: %v"
)
```

### 5. 常量集中定义

**优化方案:**
```go
// 在文件顶部定义所有常量
const (
    // 父路径常量
    pathDevice              = "Device."
    pathLocalAgent          = "Device.LocalAgent."
    pathSubscription        = "Device.LocalAgent.Subscription."
)

// 预编译正则表达式（已存在，保持不变）
var subscriptionPathRegex = regexp.MustCompile(`^Device\.LocalAgent\.Subscription\.([1-9]\d*)\.$`)
```
