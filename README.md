# TR369 WebSocket Client

TR369 WebSocket 客户端，基于清洁架构设计。

## 项目结构

```
tr369-wss-client/
├── cmd/client/              # 应用程序入口点
├── internal/                # 私有应用代码
│   ├── domain/              # 领域层
│   ├── application/         # 应用层
│   └── infrastructure/      # 基础设施层
├── pkg/                     # 公共库代码
├── configs/                 # 配置文件
├── data/                    # 数据文件
├── docs/                    # 项目文档
├── examples/                # 示例代码
├── proto/                   # Protobuf 定义
├── scripts/                 # 构建脚本
└── test/                    # 测试代码
```

## 快速开始

### 构建
```powershell
.\scripts\build.ps1
```

### 运行
```powershell
.\bin\tr369-wss-client.exe
```

### 测试
```powershell
.\scripts\test.ps1 -Coverage
```

## 配置

配置文件位于 `configs/environments/` 目录。

## 许可证

MIT License