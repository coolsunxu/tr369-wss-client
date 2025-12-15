# 设计文档

## 概述

本设计文档描述了对 `client/repository/client.go` 中 `ClientRepository` 接口的重构方案。核心思想是将当前混合职责的接口拆分为两个独立接口：`DataRepository`（数据访问）和 `ListenerManager`（监听器管理），遵循单一职责原则（SRP）和接口隔离原则（ISP）。

## 架构

### 当前架构

```
┌─────────────────────────────────────────────────────────┐
│                   ClientRepository                       │
├─────────────────────────────────────────────────────────┤
│  数据访问方法:                                            │
│  - GetValue()                                            │
│  - SetValue()                                            │
│  - ConstructGetResp()                                    │
│  - IsExistPath()                                         │
│  - GetNewInstance()                                      │
│  - HandleDeleteRequest()                                 │
│  - StartClientRepository()                               │
├─────────────────────────────────────────────────────────┤
│  监听器管理方法:                                          │
│  - AddListener()                                         │
│  - RemoveListener()                                      │
│  - ResetListener()                                       │
│  - NotifyListeners()                                     │
└─────────────────────────────────────────────────────────┘
```

### 重构后架构

```
┌─────────────────────────────┐  ┌─────────────────────────────┐
│       DataRepository        │  │      ListenerManager        │
├─────────────────────────────┤  ├─────────────────────────────┤
│  - GetValue()               │  │  - AddListener()            │
│  - SetValue()               │  │  - RemoveListener()         │
│  - ConstructGetResp()       │  │  - ResetListener()          │
│  - IsExistPath()            │  │  - NotifyListeners()        │
│  - GetNewInstance()         │  │                             │
│  - HandleDeleteRequest()    │  │                             │
│  - Start()                  │  │                             │
└─────────────────────────────┘  └─────────────────────────────┘
              │                              │
              └──────────────┬───────────────┘
                             │
                             ▼
              ┌─────────────────────────────┐
              │      clientRepository       │
              │    (实现两个接口)            │
              └─────────────────────────────┘
```

## 组件和接口

### 1. DataRepository 接口

负责 TR181 数据模型的纯数据 CRUD 操作，不包含业务逻辑：

```go
// DataRepository 定义数据访问接口
type DataRepository interface {
    // GetValue 获取指定路径的值
    GetValue(path string) (interface{}, error)

    // GetParameters 获取底层参数数据（供 UseCase 构建响应使用）
    GetParameters() map[string]interface{}

    // SetValue 设置指定路径的值
    // 返回: changed (是否发生变化), oldValue (旧值)
    SetValue(path string, key string, value string) (changed bool, oldValue string)

    // DeleteNode 删除指定路径的节点
    DeleteNode(path string) (nodePath string, isFound bool)

    // Start 启动数据仓库（初始化和数据同步）
    Start()
}
```

**已移除的方法:**
- `ConstructGetResp`: 构造 API 响应是业务逻辑，移至 UseCase 层
- `IsExistPath`: 包含路径表达式解析（`[条件]` 语法），移至 UseCase 层
- `GetNewInstance`: 包含实例编号生成策略，移至 UseCase 层

Repository 层只提供 `GetParameters()` 方法暴露底层数据，业务逻辑由 UseCase 层处理。

### 2. ListenerManager 接口

负责事件监听器的管理：

```go
// ListenerManager 定义监听器管理接口
type ListenerManager interface {
    // AddListener 添加参数变化监听器
    AddListener(paramName string, listener tr181Model.Listener) error

    // RemoveListener 移除指定参数的监听器
    RemoveListener(paramName string) error

    // ResetListener 重置所有监听器
    ResetListener() error

    // NotifyListeners 通知指定参数的所有监听器
    NotifyListeners(paramName string, value interface{})
}
```

### 3. ClientRepository 组合接口（向后兼容）

为保持向后兼容，保留组合接口：

```go
// ClientRepository 组合接口，包含数据访问和监听器管理
// 保留此接口以保持向后兼容性
type ClientRepository interface {
    DataRepository
    ListenerManager
}
```

### 4. clientRepository 实现

实现类同时实现两个接口：

```go
type clientRepository struct {
    Config         *config.Config
    TR181DataModel *tr181Model.TR181DataModel
    writeCount     int
    lastWriteTime  int64
    pingTicker     *time.Ticker
    ctx            context.Context
    cancel         context.CancelFunc
}

// 实现 DataRepository 接口
func (repo *clientRepository) GetValue(path string) (interface{}, error) { ... }
func (repo *clientRepository) GetParameters() map[string]interface{} {
    return repo.TR181DataModel.Parameters
}
func (repo *clientRepository) SetValue(path, key, value string) (bool, string) { ... }
func (repo *clientRepository) DeleteNode(path string) (string, bool) { ... }
// ... 其他数据访问方法

// 实现 ListenerManager 接口
func (repo *clientRepository) AddListener(paramName string, listener tr181Model.Listener) error { ... }
func (repo *clientRepository) RemoveListener(paramName string) error { ... }
// ... 其他监听器方法
```

### 5. UseCase 层业务逻辑

将业务逻辑方法移至 UseCase 层：

```go
// 构建 GET 响应
func (uc *ClientUseCase) constructGetResp(paths []string) api.Response_GetResp {
    params := uc.DataRepo.GetParameters()
    return trtree.ConstructGetResp(params, paths)
}

// 检查路径是否存在（包含路径表达式解析）
func (uc *ClientUseCase) isExistPath(path string) (isSuccess bool, nodePath string) {
    params := uc.DataRepo.GetParameters()
    return trtree.IsExistPath(params, path)
}

// 获取新实例路径（包含实例编号生成策略）
func (uc *ClientUseCase) getNewInstance(path string) string {
    params := uc.DataRepo.GetParameters()
    return trtree.GetNewInstance(params, path)
}
```

**设计原则:**
- Repository 层只负责数据存储和基本 CRUD
- UseCase 层负责业务逻辑（路径解析、响应构建、实例编号生成）
- trtree 包提供底层树操作工具函数

## 数据模型

数据模型保持不变：

```go
// TR181DataModel 存储 TR181 参数和监听器
type TR181DataModel struct {
    Parameters map[string]interface{}
    Listeners  map[string][]Listener
}

// Listener 监听器定义
type Listener struct {
    SubscriptionId string
    Listener       func(subscriptionId string, value interface{})
}
```



## 正确性属性

*属性是一种特征或行为，应该在系统的所有有效执行中保持为真——本质上是关于系统应该做什么的形式化陈述。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### Property 1: 行为等价性

*对于任何* 有效的数据操作序列（GetValue、SetValue、HandleDeleteRequest 等），通过新的 DataRepository 接口执行应产生与通过原 ClientRepository 接口执行相同的结果。

**验证: 需求 3.1, 3.3**

### Property 2: 监听器管理等价性

*对于任何* 有效的监听器操作序列（AddListener、RemoveListener、NotifyListeners），通过新的 ListenerManager 接口执行应产生与通过原 ClientRepository 接口执行相同的行为。

**验证: 需求 3.1**

## 错误处理

### 错误处理策略

重构不改变现有的错误处理逻辑：

1. **GetValue**: 路径不存在时返回 error
2. **SetValue**: 始终成功，返回变化状态和旧值
3. **AddListener/RemoveListener**: 返回 error 表示操作失败
4. **ResetListener**: 返回 error 表示重置失败

### 错误类型

```go
// 常见错误场景
- 路径不存在: "path not found: %s"
- 监听器操作失败: 返回具体的 error
```

## 测试策略

### 单元测试

- 测试 DataRepository 接口的各方法独立工作
- 测试 ListenerManager 接口的各方法独立工作
- 测试接口可以被独立模拟

### 属性测试

使用 Go 的 `testing/quick` 包进行属性测试：

- 配置每个属性测试运行至少 100 次迭代
- 每个属性测试需标注对应的正确性属性：`**Feature: client-repository-optimization, Property {number}: {property_text}**`

### 测试覆盖范围

1. **DataRepository**: 测试所有数据访问方法
2. **ListenerManager**: 测试监听器的添加、移除、通知
3. **接口隔离**: 验证两个接口可以独立模拟和测试
4. **向后兼容**: 验证组合接口 ClientRepository 保持原有行为

## 重构详细设计

### 1. 接口分离

**步骤:**

1. 在 `client/model/client.go` 中定义 `DataRepository` 接口
2. 在 `client/model/client.go` 中定义 `ListenerManager` 接口
3. 修改 `ClientRepository` 为组合接口（嵌入上述两个接口）
4. 确保 `clientRepository` 实现保持不变

### 2. 方法重命名（可选优化）

考虑将 `StartClientRepository()` 重命名为 `Start()`，更简洁：

```go
// 旧方法
func (repo *clientRepository) StartClientRepository()

// 新方法
func (repo *clientRepository) Start()
```

### 3. UseCase 层适配

UseCase 层采用分别注入 `DataRepository` 和 `ListenerManager` 的方式，提供更好的灵活性和可测试性：

```go
// 新的依赖注入方式
type ClientUseCase struct {
    Config          *config.Config
    DataRepo        model.DataRepository     // 数据访问接口
    ListenerMgr     model.ListenerManager    // 监听器管理接口
    ctx             context.Context
    messageChannel  chan []byte
}

// 构造函数更新
func NewClientUseCase(
    ctx context.Context,
    cfg *config.Config,
    dataRepo model.DataRepository,
    listenerMgr model.ListenerManager,
    messageChannel chan []byte,
) *ClientUseCase {
    return &ClientUseCase{
        ctx:            ctx,
        Config:         cfg,
        DataRepo:       dataRepo,
        ListenerMgr:    listenerMgr,
        messageChannel: messageChannel,
    }
}
```

**优势:**
- 可以独立模拟 DataRepository 和 ListenerManager 进行单元测试
- 支持注入不同的实现（如内存实现、持久化实现）
- 更清晰地表达 UseCase 的依赖关系

### 4. 文件组织

接口定义集中在 `client/model/client.go`：

```go
package model

// DataRepository 数据访问接口
type DataRepository interface { ... }

// ListenerManager 监听器管理接口
type ListenerManager interface { ... }

// ClientRepository 组合接口（向后兼容）
type ClientRepository interface {
    DataRepository
    ListenerManager
}
```

## 迁移计划

1. **阶段一**: 定义新接口，保持实现不变
2. **阶段二**: 验证所有现有功能正常工作
3. **阶段三**: （可选）逐步将 UseCase 迁移到使用细粒度接口
