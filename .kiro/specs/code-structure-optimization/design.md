# 代码目录结构优化设计文档

## 概述

本设计文档基于清洁架构（Clean Architecture）原则，为 TR369 WebSocket 客户端项目提供系统性的代码目录结构重构方案。当前项目虽然采用了分层架构，但存在职责边界不清晰、依赖关系混乱、测试组织不规范等问题。

通过实施本设计，项目将获得更好的可维护性、可测试性和可扩展性，同时保持向后兼容性，确保现有功能不受影响。

## 架构设计

### 清洁架构四层模型

```
┌─────────────────────────────────────────────────────────────┐
│                    Interface Layer                          │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              Application Layer                      │    │
│  │  ┌─────────────────────────────────────────────┐    │    │
│  │  │              Domain Layer                   │    │    │
│  │  │  ┌─────────────────────────────────────┐    │    │    │
│  │  │  │        Infrastructure Layer         │    │    │    │
│  │  │  └─────────────────────────────────────┘    │    │    │
│  │  └─────────────────────────────────────────────┘    │    │
│  └─────────────────────────────────────────────────────┐    │
└─────────────────────────────────────────────────────────────┘
```

**依赖方向**: 外层依赖内层，内层不依赖外层

### 目录结构设计

```
tr369-wss-client/
├── cmd/                           # 应用程序入口点
│   └── client/
│       └── main.go
├── internal/                      # 私有应用代码
│   ├── domain/                    # 领域层 - 业务逻辑核心
│   │   ├── entities/              # 实体
│   │   │   ├── tr181/
│   │   │   │   ├── datamodel.go
│   │   │   │   ├── subscription.go
│   │   │   │   └── listener.go
│   │   │   └── usp/
│   │   │       ├── message.go
│   │   │       └── record.go
│   │   ├── repositories/          # 仓储接口
│   │   │   ├── client.go
│   │   │   └── tr181.go
│   │   ├── services/              # 领域服务接口
│   │   │   ├── message_handler.go
│   │   │   └── subscription_manager.go
│   │   └── valueobjects/          # 值对象
│   │       ├── endpoint_id.go
│   │       └── message_type.go
│   ├── application/               # 应用层 - 用例编排
│   │   ├── usecases/              # 用例实现
│   │   │   ├── client/
│   │   │   │   ├── connect.go
│   │   │   │   ├── disconnect.go
│   │   │   │   └── message_handling.go
│   │   │   └── tr181/
│   │   │       ├── get_parameter.go
│   │   │       ├── set_parameter.go
│   │   │       ├── add_object.go
│   │   │       ├── delete_object.go
│   │   │       └── subscription_management.go
│   │   ├── services/              # 应用服务
│   │   │   ├── message_processor.go
│   │   │   └── notification_service.go
│   │   └── dto/                   # 数据传输对象
│   │       ├── request/
│   │       └── response/
│   ├── infrastructure/            # 基础设施层 - 外部依赖
│   │   ├── websocket/             # WebSocket 实现
│   │   │   ├── client.go
│   │   │   ├── connection.go
│   │   │   └── message_handler.go
│   │   ├── persistence/           # 数据持久化
│   │   │   ├── json/
│   │   │   │   ├── tr181_repository.go
│   │   │   │   └── file_manager.go
│   │   │   └── memory/
│   │   │       └── cache_repository.go
│   │   ├── config/                # 配置管理
│   │   │   ├── loader.go
│   │   │   ├── validator.go
│   │   │   └── models.go
│   │   ├── logging/               # 日志实现
│   │   │   └── zap_logger.go
│   │   └── protobuf/              # Protocol Buffers 处理
│   │       ├── encoder.go
│   │       ├── decoder.go
│   │       └── message_factory.go
│   └── interfaces/                # 接口层 - 外部接口
│       ├── cli/                   # 命令行接口
│       │   ├── commands/
│       │   └── flags.go
│       ├── websocket/             # WebSocket 接口适配器
│       │   ├── handlers/
│       │   └── middleware/
│       └── api/                   # API 接口（如果需要）
├── pkg/                           # 公共库代码
│   ├── api/                       # 生成的 Protocol Buffers 代码
│   ├── utils/                     # 通用工具
│   │   ├── json/
│   │   ├── string/
│   │   └── crypto/
│   └── errors/                    # 错误定义
├── configs/                       # 配置文件
│   ├── environments/
│   │   ├── development.json
│   │   ├── production.json
│   │   └── test.json
│   └── schemas/
├── scripts/                       # 构建和部署脚本
│   ├── build.sh
│   ├── test.sh
│   └── generate.sh
├── test/                          # 测试代码
│   ├── integration/               # 集成测试
│   ├── e2e/                       # 端到端测试
│   ├── fixtures/                  # 测试数据
│   └── mocks/                     # 模拟对象
├── docs/                          # 项目文档
│   ├── architecture/
│   ├── api/
│   └── deployment/
├── deployments/                   # 部署配置
│   ├── docker/
│   └── kubernetes/
└── examples/                      # 示例代码
```

## 组件和接口设计

### 领域层接口定义

#### 仓储接口
```go
// internal/domain/repositories/tr181.go
type TR181Repository interface {
    GetParameter(path string) (*entities.Parameter, error)
    SetParameter(path, key, value string) error
    AddObject(path string, params map[string]string) (string, error)
    DeleteObject(path string) error
    SaveData() error
    LoadData() error
}

// internal/domain/repositories/client.go
type ClientRepository interface {
    AddListener(path string, listener *entities.Listener) error
    RemoveListener(path string) error
    NotifyListeners(path string, event interface{}) error
}
```

#### 服务接口
```go
// internal/domain/services/message_handler.go
type MessageHandler interface {
    HandleGet(msg *entities.USPMessage) (*entities.USPMessage, error)
    HandleSet(msg *entities.USPMessage) (*entities.USPMessage, error)
    HandleAdd(msg *entities.USPMessage) (*entities.USPMessage, error)
    HandleDelete(msg *entities.USPMessage) (*entities.USPMessage, error)
    HandleOperate(msg *entities.USPMessage) (*entities.USPMessage, error)
}
```

### 应用层用例

#### 连接管理用例
```go
// internal/application/usecases/client/connect.go
type ConnectUseCase struct {
    wsClient     domain.WebSocketClient
    config       *config.Config
    logger       domain.Logger
}

func (uc *ConnectUseCase) Execute(ctx context.Context) error {
    // 连接逻辑实现
}
```

### 基础设施层实现

#### WebSocket 客户端实现
```go
// internal/infrastructure/websocket/client.go
type WSClient struct {
    conn           *websocket.Conn
    messageHandler domain.MessageHandler
    config         *config.Config
}

func (c *WSClient) Connect(ctx context.Context, url string) error {
    // WebSocket 连接实现
}
```

## 数据模型设计

### 领域实体

#### TR181 数据模型
```go
// internal/domain/entities/tr181/datamodel.go
type DataModel struct {
    parameters map[string]*Parameter
    listeners  map[string][]*Listener
}

type Parameter struct {
    Path     string
    Value    interface{}
    Type     ParameterType
    Writable bool
}
```

#### USP 消息实体
```go
// internal/domain/entities/usp/message.go
type Message struct {
    Header *Header
    Body   *Body
}

type Header struct {
    MessageType MessageType
    MessageID   string
}
```

### 值对象

#### 端点标识符
```go
// internal/domain/valueobjects/endpoint_id.go
type EndpointID struct {
    value string
}

func NewEndpointID(id string) (*EndpointID, error) {
    if id == "" {
        return nil, errors.New("endpoint ID cannot be empty")
    }
    return &EndpointID{value: id}, nil
}
```

## 错误处理设计

### 错误类型定义
```go
// pkg/errors/domain.go
type DomainError struct {
    Code    string
    Message string
    Cause   error
}

// pkg/errors/application.go
type ApplicationError struct {
    Type    ErrorType
    Message string
    Details map[string]interface{}
}
```

### 错误处理策略
- 领域层：定义业务错误类型
- 应用层：处理用例级别错误
- 基础设施层：处理技术错误
- 接口层：错误转换和响应

## 测试策略

### 单元测试组织
```
internal/
├── domain/
│   ├── entities/
│   │   └── tr181/
│   │       ├── datamodel.go
│   │       └── datamodel_test.go
│   └── services/
│       ├── message_handler.go
│       └── message_handler_test.go
├── application/
│   └── usecases/
│       └── client/
│           ├── connect.go
│           └── connect_test.go
└── infrastructure/
    └── websocket/
        ├── client.go
        └── client_test.go
```

### 集成测试
```
test/
├── integration/
│   ├── websocket_integration_test.go
│   ├── tr181_repository_test.go
│   └── message_flow_test.go
└── fixtures/
    ├── test_config.json
    └── sample_messages.json
```

### 属性基础测试库选择
- **Go 语言推荐**: `gopter` 库
- **配置要求**: 每个属性测试运行最少 100 次迭代
- **测试标记**: 使用注释标记属性测试与设计文档的对应关系

## 正确性属性

*属性是指在系统的所有有效执行中都应该成立的特征或行为——本质上是关于系统应该做什么的正式声明。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### 架构依赖属性

**属性 1: 依赖倒置原则遵循**
*对于任何* 内层模块到外层模块的调用，都应该通过接口进行而不是直接依赖具体实现
**验证需求: 需求 1.2**

**属性 2: 基础设施层隔离性**
*对于任何* 基础设施层的修改，领域层和应用层的测试应该继续通过而不受影响
**验证需求: 需求 1.3**

**属性 3: 外部依赖隔离**
*对于任何* 新添加的外部依赖，都应该只出现在基础设施层中
**验证需求: 需求 1.4**

**属性 4: 跨层调用接口化**
*对于任何* 跨层的方法调用，都应该通过定义的接口进行
**验证需求: 需求 1.5**

### 代码组织属性

**属性 5: 功能模块内聚性**
*对于任何* 特定功能，其相关的所有文件都应该组织在对应的功能模块目录下
**验证需求: 需求 2.2**

**属性 6: 目录结构一致性**
*对于任何* 功能模块，都应该遵循统一的目录结构模式
**验证需求: 需求 2.3**

**属性 7: 配置代码集中性**
*对于任何* 配置相关的代码，都应该位于配置管理模块中
**验证需求: 需求 2.4**

**属性 8: 共享代码位置规范**
*对于任何* 被多个模块使用的代码，都应该放置在明确的共享目录中
**验证需求: 需求 2.5**

### 接口设计属性

**属性 9: 业务接口领域层定位**
*对于任何* 业务接口定义，都应该位于领域层中
**验证需求: 需求 3.1**

**属性 10: 外部服务实现基础设施层定位**
*对于任何* 外部服务接口的实现，都应该位于基础设施层中
**验证需求: 需求 3.2**

**属性 11: 模块间接口通信**
*对于任何* 模块间的通信，都应该通过定义良好的接口进行
**验证需求: 需求 3.3**

**属性 12: 接口模拟测试支持**
*对于任何* 定义的接口，都应该能够创建模拟对象进行单元测试
**验证需求: 需求 3.4**

**属性 13: 接口向后兼容性**
*对于任何* 接口的变更，都应该保持向后兼容性不破坏现有功能
**验证需求: 需求 3.5**

### 测试组织属性

**属性 14: 测试文件同目录组织**
*对于任何* 源代码文件，其对应的单元测试文件应该位于相同目录中
**验证需求: 需求 4.1**

**属性 15: 集成测试专门目录**
*对于任何* 集成测试，都应该位于专门的集成测试目录中
**验证需求: 需求 4.2**

**属性 16: 测试数据统一管理**
*对于任何* 测试数据文件，都应该位于统一的测试数据目录中
**验证需求: 需求 4.3**

**属性 17: 测试工具独立组织**
*对于任何* 测试工具代码，都应该位于独立的测试工具目录中
**验证需求: 需求 4.5**

### 配置和文档属性

**属性 18: 环境配置分离**
*对于任何* 环境配置，不同环境的配置文件应该分别存储
**验证需求: 需求 5.1**

**属性 19: 文档类型化组织**
*对于任何* 项目文档，都应该按类型和用途进行组织
**验证需求: 需求 5.2**

**属性 20: 生成代码分离**
*对于任何* 代码生成操作，生成的代码应该与手写代码分离存储
**验证需求: 需求 5.4**

**属性 21: 示例代码独立管理**
*对于任何* 示例代码，都应该在独立的示例目录中管理
**验证需求: 需求 5.5**

### 模块化属性

**属性 22: 模块独立开发支持**
*对于任何* 功能模块，都应该支持独立的编译和测试
**验证需求: 需求 6.1**

**属性 23: 模块依赖接口化**
*对于任何* 模块间的依赖关系，都应该通过明确的接口定义
**验证需求: 需求 6.2**

**属性 24: 模块边界封装性**
*对于任何* 模块，其内部实现不应该被外部模块直接访问
**验证需求: 需求 6.3**

**属性 25: 模块独立测试能力**
*对于任何* 模块，都应该能够独立运行其测试套件
**验证需求: 需求 6.4**

**属性 26: 模块配置独立性**
*对于任何* 需要配置的模块，都应该有独立的配置管理机制
**验证需求: 需求 6.5**

## 测试策略

### 双重测试方法

本项目将采用单元测试和属性基础测试相结合的方法：

- **单元测试**: 验证具体示例、边界情况和错误条件
- **属性基础测试**: 验证应该在所有输入中保持的通用属性
- **集成测试**: 验证组件间的交互和端到端流程

### 属性基础测试要求

- **测试库**: 使用 `gopter` 库进行属性基础测试
- **迭代次数**: 每个属性测试最少运行 100 次迭代
- **测试标记**: 每个属性基础测试必须使用注释明确标记对应的设计文档属性
- **标记格式**: `**Feature: code-structure-optimization, Property {number}: {property_text}**`

### 测试覆盖策略

1. **架构测试**: 使用静态分析验证架构约束
2. **结构测试**: 验证目录结构和文件组织
3. **依赖测试**: 检查模块间依赖关系
4. **接口测试**: 验证接口定义和实现的一致性
5. **集成测试**: 验证重构后的系统功能完整性

### 测试工具和框架

- **单元测试**: Go 标准 `testing` 包
- **属性测试**: `gopter` 库
- **静态分析**: `go vet`, `golangci-lint`
- **架构测试**: 自定义架构验证工具
- **覆盖率**: `go test -cover`

## 迁移策略

### 渐进式重构方法

1. **阶段一**: 创建新的目录结构
2. **阶段二**: 定义领域层接口
3. **阶段三**: 迁移实体和值对象
4. **阶段四**: 重构应用层用例
5. **阶段五**: 迁移基础设施层实现
6. **阶段六**: 更新接口层适配器
7. **阶段七**: 完善测试覆盖

### 向后兼容性保证

- 保持现有 API 接口不变
- 逐步废弃旧的内部接口
- 提供迁移指南和工具
- 维护功能等价性测试

### 风险缓解

- 每个阶段都有回滚计划
- 保持完整的测试覆盖
- 渐进式部署和验证
- 详细的变更日志记录