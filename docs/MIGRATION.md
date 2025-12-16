# 代码结构迁移指南

## 概述

本文档描述了 TR369 WebSocket 客户端项目从旧结构迁移到新的清洁架构结构的过程。迁移已完成，旧代码已删除。

## 新结构

```
tr369-wss-client/
├── cmd/client/                    # 应用入口点
├── internal/                      # 私有代码
│   ├── domain/                    # 领域层
│   │   ├── entities/              # 领域实体
│   │   │   ├── tr181/             # TR181 数据模型
│   │   │   └── usp/               # USP 消息实体
│   │   ├── repositories/          # 仓储接口
│   │   ├── services/              # 服务接口
│   │   ├── valueobjects/          # 值对象
│   │   └── errors.go              # 领域错误
│   ├── application/               # 应用层
│   │   ├── usecases/              # 用例实现
│   │   │   ├── client/            # 客户端用例
│   │   │   └── tr181/             # TR181 用例
│   │   └── services/              # 应用服务
│   └── infrastructure/            # 基础设施层
│       ├── config/                # 配置管理
│       ├── di/                    # 依赖注入
│       ├── logging/               # 日志实现
│       ├── persistence/           # 持久化
│       │   ├── json/              # JSON 文件管理
│       │   ├── repository/        # 仓储实现
│       │   └── trtree/            # 树结构操作
│       ├── protobuf/              # Protobuf 编解码
│       └── websocket/             # WebSocket 客户端
├── pkg/                           # 公共库
│   ├── api/                       # API 定义
│   ├── errors/                    # 错误定义
│   └── utils/                     # 工具函数
├── configs/environments/          # 环境配置文件
├── test/                          # 测试代码
│   ├── fixtures/                  # 测试数据
│   ├── integration/               # 集成测试
│   └── mocks/                     # 模拟对象
├── docs/                          # 文档
├── examples/                      # 示例代码
├── proto/                         # Protobuf 定义
└── scripts/                       # 构建脚本
```

## 迁移步骤

### 1. 领域层迁移

将业务实体和接口迁移到 `internal/domain/`:

- `tr181/model/` → `internal/domain/entities/tr181/`
- 接口定义 → `internal/domain/repositories/` 和 `internal/domain/services/`

### 2. 应用层迁移

将用例逻辑迁移到 `internal/application/`:

- `client/usecase/` → `internal/application/usecases/`

### 3. 基础设施层迁移

将外部依赖实现迁移到 `internal/infrastructure/`:

- `client/client.go` → `internal/infrastructure/websocket/`
- `config/` → `internal/infrastructure/config/`
- `log/` → `internal/infrastructure/logging/`

### 4. 公共库迁移

将共享代码迁移到 `pkg/`:

- `utils/` → `pkg/utils/`
- `common/` → `pkg/utils/`

## 架构原则

### 依赖规则
- 领域层（domain）不依赖任何其他层
- 应用层（application）只依赖领域层
- 基础设施层（infrastructure）实现领域层定义的接口

### 接口定义
所有业务接口定义在领域层：
- `services.Logger` - 日志服务接口
- `services.ConfigProvider` - 配置提供者接口
- `services.WebSocketClient` - WebSocket 客户端接口
- `repositories.ClientRepository` - 客户端仓储接口
- `repositories.TR181Repository` - TR181 仓储接口

### 依赖注入
使用 `internal/infrastructure/di/container.go` 管理依赖注入。

## 最佳实践

1. 领域层不应依赖基础设施层
2. 通过接口实现依赖倒置
3. 使用依赖注入管理组件
4. 测试文件与源文件放在同一目录
5. 集成测试放在 `test/integration/` 目录
6. 外部依赖只在基础设施层使用

## 常见问题

### Q: 如何处理循环依赖？
A: 通过在领域层定义接口，在基础设施层实现来打破循环依赖。

### Q: 如何添加新功能？
A: 
1. 在领域层定义接口和实体
2. 在应用层实现用例
3. 在基础设施层实现具体依赖

### Q: 如何运行测试？
A: 使用 `.\scripts\test.ps1 -Verbose` 运行所有测试。
