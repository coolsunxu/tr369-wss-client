# 架构文档

## 概述

TR369 WebSocket 客户端采用清洁架构（Clean Architecture）设计，将代码组织为四个主要层次：

1. **领域层 (Domain Layer)** - 业务逻辑核心
2. **应用层 (Application Layer)** - 用例编排
3. **基础设施层 (Infrastructure Layer)** - 外部依赖实现
4. **接口层 (Interface Layer)** - 外部接口适配

## 目录结构

```
tr369-wss-client/
├── cmd/                    # 应用程序入口点
├── internal/               # 私有应用代码
│   ├── domain/             # 领域层
│   │   ├── entities/       # 实体
│   │   ├── repositories/   # 仓储接口
│   │   ├── services/       # 服务接口
│   │   └── valueobjects/   # 值对象
│   ├── application/        # 应用层
│   │   └── usecases/       # 用例实现
│   ├── infrastructure/     # 基础设施层
│   │   ├── config/         # 配置管理
│   │   ├── logging/        # 日志实现
│   │   ├── persistence/    # 数据持久化
│   │   ├── protobuf/       # Protocol Buffers
│   │   └── websocket/      # WebSocket 实现
│   └── interfaces/         # 接口层
├── pkg/                    # 公共库代码
│   ├── api/                # 生成的 Protobuf 代码
│   ├── errors/             # 错误定义
│   └── utils/              # 通用工具
├── configs/                # 配置文件
├── scripts/                # 构建和部署脚本
├── test/                   # 测试代码
└── docs/                   # 项目文档
```

## 依赖规则

- 内层不依赖外层
- 外层通过接口依赖内层
- 依赖注入用于解耦组件

## 更多信息

- [领域层设计](./domain.md)
- [应用层设计](./application.md)
- [基础设施层设计](./infrastructure.md)
