# 需求文档

## 简介

本文档定义了对订阅路径匹配和参数校验功能的优化需求。当前 `NotifyListeners` 使用精确匹配逻辑，无法支持 TR181 数据模型中的层级路径匹配。同时，订阅参数（如 `ReferenceList`）缺乏符合 USP/TR181 Path Name 规范的校验。优化目标是实现符合 TR-369/TR-181 规范的路径匹配和参数校验机制。

## 术语表

- **Path Name**: TR181 数据模型中的路径名称，如 `Device.WiFi.Radio.1.Enabled`
- **ReferenceList**: 订阅中的引用路径列表，包含一个或多个 Path Name
- **完全限定路径**: 从 `Device.` 开始的绝对路径，不允许相对路径
- **多实例对象**: 包含实例标识符的对象路径，如 `Device.Ethernet.Interface.{i}` 或 `Device.Ethernet.Interface.2`
- **对象路径**: 指向对象的路径，通常以 `.` 结尾表示对象本身
- **参数路径**: 指向具体参数的路径，如 `Device.WiFi.Radio.1.Enabled`
- **NotifyListeners**: 监听器通知方法，当参数变化时通知相关订阅
- **ValueChange**: 值变化订阅类型
- **ObjectCreation**: 对象创建订阅类型
- **ObjectDeletion**: 对象删除订阅类型

## 需求

### 需求 1

**用户故事:** 作为 USP Agent 开发者，我希望 NotifyListeners 支持层级路径匹配，以便订阅父路径时能够接收到子路径的变化通知。

#### 验收标准

1. WHEN 订阅路径为 `Device.WiFi.` 且参数 `Device.WiFi.Radio.1.Enabled` 发生变化 THEN NotifyListeners SHALL 触发该订阅的监听器
2. WHEN 订阅路径为 `Device.WiFi.Radio.1.Enabled` 且该参数发生变化 THEN NotifyListeners SHALL 触发该订阅的监听器（精确匹配）
3. WHEN 订阅路径为 `Device.WiFi.Radio.1.` 且参数 `Device.WiFi.Radio.2.Enabled` 发生变化 THEN NotifyListeners SHALL 不触发该订阅的监听器
4. WHEN 订阅路径为 `Device.` 且任意 Device 下的参数发生变化 THEN NotifyListeners SHALL 触发该订阅的监听器

### 需求 2

**用户故事:** 作为 USP Agent 开发者，我希望订阅的 ReferenceList 参数在添加时被校验，以便拒绝不符合 Path Name 规范的路径。

#### 验收标准

1. WHEN 添加订阅且 ReferenceList 为完全限定路径（以 `Device.` 开头）THEN 系统 SHALL 接受该订阅
2. WHEN 添加订阅且 ReferenceList 为相对路径（不以 `Device.` 开头）THEN 系统 SHALL 拒绝该订阅并返回错误
3. WHEN 添加订阅且 ReferenceList 包含非法字符（如空字符、控制字符）THEN 系统 SHALL 拒绝该订阅并返回错误
4. WHEN 添加订阅且 ReferenceList 路径长度超过 256 字符 THEN 系统 SHALL 拒绝该订阅并返回错误
5. WHEN 添加订阅且 ReferenceList 路径格式正确 THEN 系统 SHALL 使用 UTF-8 编码存储路径

### 需求 3

**用户故事:** 作为 USP Agent 开发者，我希望路径校验能够识别不同类型的路径格式，以便正确处理对象路径和参数路径。

#### 验收标准

1. WHEN 校验对象路径（如 `Device.WiFi.Radio.1.`）THEN 校验器 SHALL 识别为有效的对象路径
2. WHEN 校验参数路径（如 `Device.WiFi.Radio.1.Enabled`）THEN 校验器 SHALL 识别为有效的参数路径
3. WHEN 校验多实例对象路径（如 `Device.Ethernet.Interface.2`）THEN 校验器 SHALL 识别为有效的实例路径
4. WHEN 校验路径末尾有多余的点（如 `Device.WiFi..`）THEN 校验器 SHALL 识别为无效路径
5. WHEN 校验路径包含连续的点（如 `Device..WiFi`）THEN 校验器 SHALL 识别为无效路径

### 需求 4

**用户故事:** 作为 USP Agent 开发者，我希望路径匹配逻辑能够正确处理通配符和实例占位符，以便支持灵活的订阅模式。

#### 验收标准

1. WHEN 订阅路径包含实例占位符 `{i}`（如 `Device.Ethernet.Interface.{i}.`）且任意实例的参数变化 THEN NotifyListeners SHALL 触发该订阅的监听器
2. WHEN 订阅路径为具体实例（如 `Device.Ethernet.Interface.1.`）且其他实例参数变化 THEN NotifyListeners SHALL 不触发该订阅的监听器
3. WHEN 变化路径与订阅路径完全匹配 THEN NotifyListeners SHALL 优先使用精确匹配的监听器

### 需求 5

**用户故事:** 作为 USP Agent 开发者，我希望路径校验提供清晰的错误信息，以便快速定位和修复配置问题。

#### 验收标准

1. WHEN 路径校验失败 THEN 系统 SHALL 返回包含具体错误原因的错误信息
2. WHEN 路径为空字符串 THEN 系统 SHALL 返回 "路径不能为空" 错误
3. WHEN 路径不以 `Device.` 开头 THEN 系统 SHALL 返回 "路径必须以 Device. 开头" 错误
4. WHEN 路径包含非法字符 THEN 系统 SHALL 返回 "路径包含非法字符" 错误并指出具体字符位置

### 需求 6

**用户故事:** 作为 USP Agent 开发者，我希望路径校验和匹配逻辑具有良好的性能，以便在高频率参数变化场景下不影响系统响应。

#### 验收标准

1. WHEN 执行路径校验 THEN 校验器 SHALL 在 O(n) 时间复杂度内完成（n 为路径长度）
2. WHEN 执行路径匹配 THEN NotifyListeners SHALL 使用前缀树或类似数据结构优化匹配效率
3. WHEN 存在大量订阅（超过 100 个）THEN 路径匹配 SHALL 保持稳定的响应时间

