# 需求文档

## 简介

本文档定义了对 `client/repository/client.go` 中 `ClientRepository` 接口进行优化重构的需求。当前接口混合了数据访问职责和事件监听职责，违反了单一职责原则（SRP）。优化目标是将接口按职责分离，提高代码的可维护性、可测试性和架构清晰度。

## 术语表

- **ClientRepository**: 客户端数据仓库接口，负责 TR181 数据模型的 CRUD 操作
- **TR181DataModel**: TR181 数据模型，存储设备参数的树形结构
- **Listener**: 监听器，用于监听参数变化事件的回调函数
- **SRP**: Single Responsibility Principle，单一职责原则
- **ISP**: Interface Segregation Principle，接口隔离原则

## 需求

### 需求 1

**用户故事:** 作为开发者，我希望 Repository 接口只包含数据访问相关的方法，以便接口职责清晰且易于理解。

#### 验收标准

1. WHEN 定义 ClientRepository 接口 THEN ClientRepository SHALL 仅包含数据读写操作方法
2. WHEN 需要监听器功能 THEN 系统 SHALL 通过独立的 ListenerManager 接口提供
3. WHEN 查看接口定义 THEN 开发者 SHALL 能够清晰识别每个接口的单一职责

### 需求 2

**用户故事:** 作为开发者，我希望监听器管理功能独立于数据访问，以便可以单独测试和替换监听器实现。

#### 验收标准

1. WHEN 创建 ListenerManager 接口 THEN ListenerManager SHALL 包含 AddListener、RemoveListener、ResetListener 和 NotifyListeners 方法
2. WHEN UseCase 需要监听器功能 THEN UseCase SHALL 通过 ListenerManager 接口访问
3. WHEN 测试监听器功能 THEN 测试代码 SHALL 能够独立模拟 ListenerManager 而不影响 Repository

### 需求 3

**用户故事:** 作为开发者，我希望重构后的代码保持向后兼容，以便现有功能不受影响。

#### 验收标准

1. WHEN 重构完成后 THEN 所有现有的 UseCase 调用 SHALL 继续正常工作
2. WHEN 启动客户端 THEN 数据同步和监听器功能 SHALL 保持原有行为
3. WHEN 处理 USP 消息 THEN 消息处理流程 SHALL 产生与重构前相同的结果

### 需求 4

**用户故事:** 作为开发者，我希望接口命名清晰且符合 Go 语言惯例，以便代码易于阅读和维护。

#### 验收标准

1. WHEN 命名接口 THEN 接口名称 SHALL 使用描述性名词（如 DataRepository、ListenerManager）
2. WHEN 命名方法 THEN 方法名称 SHALL 使用动词开头并清晰表达操作意图
3. WHEN 定义接口 THEN 接口 SHALL 遵循 Go 语言的小接口设计原则

### 需求 5

**用户故事:** 作为开发者，我希望依赖注入更加灵活，以便可以根据需要组合不同的实现。

#### 验收标准

1. WHEN 创建 ClientUseCase 实例 THEN ClientUseCase SHALL 接受独立的 Repository 和 ListenerManager 依赖
2. WHEN 需要不同的监听器实现 THEN 系统 SHALL 支持注入自定义的 ListenerManager 实现
3. WHEN 进行单元测试 THEN 测试代码 SHALL 能够分别模拟 Repository 和 ListenerManager

