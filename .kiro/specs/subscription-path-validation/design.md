# 设计文档

## 概述

本设计文档描述了订阅路径匹配和参数校验功能的优化方案。核心改进包括：
1. 将 `NotifyListeners` 的精确匹配升级为层级前缀匹配
2. 新增 `PathValidator` 组件用于校验 ReferenceList 参数
3. 支持 `{i}` 实例占位符的通配匹配

## 架构

### 当前架构

```
┌─────────────────────────────────────────────────────────┐
│                   ListenerManager                        │
├─────────────────────────────────────────────────────────┤
│  Listeners: map[string][]Listener                       │
│  (精确匹配: paramName == subscriptionPath)              │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                   NotifyListeners                        │
│  listeners, exists := Listeners[paramName]              │
│  (仅触发完全匹配的监听器)                                │
└─────────────────────────────────────────────────────────┘
```

### 优化后架构

```
┌─────────────────────────────────────────────────────────┐
│                    PathValidator                         │
├─────────────────────────────────────────────────────────┤
│  - ValidatePath(path string) error                      │
│  - IsObjectPath(path string) bool                       │
│  - IsParameterPath(path string) bool                    │
│  - ParsePath(path string) (*PathInfo, error)            │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                   ListenerManager                        │
├─────────────────────────────────────────────────────────┤
│  Listeners: map[string][]Listener                       │
│  PathValidator: *PathValidator                          │
├─────────────────────────────────────────────────────────┤
│  - AddListener(path, listener) error                    │
│    (校验路径后添加)                                      │
│  - NotifyListeners(paramName, value)                    │
│    (层级前缀匹配 + 通配符匹配)                           │
└─────────────────────────────────────────────────────────┘
```

## 组件和接口

### 1. PathValidator 组件

负责 TR181 Path Name 规范的校验：

```go
// PathValidator 路径校验器
type PathValidator struct{}

// PathInfo 路径解析结果
type PathInfo struct {
    FullPath     string   // 完整路径
    Segments     []string // 路径段列表
    IsObject     bool     // 是否为对象路径（以.结尾）
    IsParameter  bool     // 是否为参数路径
    HasWildcard  bool     // 是否包含通配符 {i}
    InstanceNums []int    // 实例编号列表
}

// PathValidationError 路径校验错误
type PathValidationError struct {
    Path     string
    Position int    // 错误位置（-1 表示整体错误）
    Reason   string
}

func (e *PathValidationError) Error() string {
    if e.Position >= 0 {
        return fmt.Sprintf("路径校验失败 [%s] 位置 %d: %s", e.Path, e.Position, e.Reason)
    }
    return fmt.Sprintf("路径校验失败 [%s]: %s", e.Path, e.Reason)
}

// ValidatePath 校验路径是否符合 TR181 Path Name 规范
func (v *PathValidator) ValidatePath(path string) error

// IsObjectPath 判断是否为对象路径
func (v *PathValidator) IsObjectPath(path string) bool

// IsParameterPath 判断是否为参数路径
func (v *PathValidator) IsParameterPath(path string) bool

// ParsePath 解析路径信息
func (v *PathValidator) ParsePath(path string) (*PathInfo, error)
```

### 2. PathMatcher 组件

负责路径匹配逻辑：

```go
// PathMatcher 路径匹配器
type PathMatcher struct{}

// MatchResult 匹配结果
type MatchResult struct {
    Matched      bool
    MatchType    MatchType // Exact, Prefix, Wildcard
    MatchedPath  string
}

type MatchType int

const (
    MatchTypeNone MatchType = iota
    MatchTypeExact     // 精确匹配
    MatchTypePrefix    // 前缀匹配（订阅路径是变化路径的前缀）
    MatchTypeWildcard  // 通配符匹配（{i} 匹配任意实例）
)

// Match 检查变化路径是否匹配订阅路径
// subscriptionPath: 订阅的路径（可能包含 {i}）
// changedPath: 发生变化的参数路径
func (m *PathMatcher) Match(subscriptionPath, changedPath string) MatchResult

// IsPrefix 检查 subscriptionPath 是否是 changedPath 的前缀
func (m *PathMatcher) IsPrefix(subscriptionPath, changedPath string) bool

// MatchWildcard 检查带通配符的路径匹配
func (m *PathMatcher) MatchWildcard(subscriptionPath, changedPath string) bool
```

### 3. 更新后的 ListenerManager 接口

```go
// ListenerManager 定义监听器管理接口
type ListenerManager interface {
    // AddListener 添加参数变化监听器
    // 校验 paramName 是否符合 Path Name 规范
    AddListener(paramName string, listener tr181Model.Listener) error

    // RemoveListener 移除指定参数的监听器
    RemoveListener(paramName string) error

    // ResetListener 重置所有监听器
    ResetListener() error

    // NotifyListeners 通知匹配参数路径的所有监听器
    // 支持层级前缀匹配和通配符匹配
    NotifyListeners(paramName string, value interface{})
}
```

## 数据模型

### PathInfo 结构

```go
// PathInfo 路径解析结果
type PathInfo struct {
    FullPath     string   // 完整路径，如 "Device.WiFi.Radio.1.Enabled"
    Segments     []string // 路径段，如 ["Device", "WiFi", "Radio", "1", "Enabled"]
    IsObject     bool     // 是否为对象路径（以.结尾）
    IsParameter  bool     // 是否为参数路径（不以.结尾）
    HasWildcard  bool     // 是否包含通配符 {i}
    InstanceNums []int    // 实例编号位置索引
}
```

### 校验规则常量

```go
const (
    // MaxPathLength 路径最大长度
    MaxPathLength = 256
    
    // PathPrefix 路径必须以此开头
    PathPrefix = "Device."
    
    // WildcardPlaceholder 实例通配符
    WildcardPlaceholder = "{i}"
)

// 非法字符集合（控制字符 0x00-0x1F，除了常见的空白字符）
var illegalChars = []rune{
    '\x00', '\x01', '\x02', '\x03', '\x04', '\x05', '\x06', '\x07',
    '\x08', '\x0B', '\x0C', '\x0E', '\x0F', '\x10', '\x11', '\x12',
    '\x13', '\x14', '\x15', '\x16', '\x17', '\x18', '\x19', '\x1A',
    '\x1B', '\x1C', '\x1D', '\x1E', '\x1F',
}
```



## 正确性属性

*属性是一种特征或行为，应该在系统的所有有效执行中保持为真——本质上是关于系统应该做什么的形式化陈述。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### Property 1: 前缀匹配正确性

*对于任何* 有效的订阅路径 S 和变化路径 C，如果 S 是 C 的前缀（S 以 `.` 结尾且 C 以 S 开头），则 NotifyListeners(C, value) 应触发 S 的监听器。

**验证: 需求 1.1, 1.4**

### Property 2: 精确匹配优先

*对于任何* 变化路径 C，如果同时存在精确匹配的订阅和前缀匹配的订阅，NotifyListeners 应同时触发两者的监听器，且精确匹配的监听器先于前缀匹配的监听器执行。

**验证: 需求 1.2, 4.3**

### Property 3: 实例隔离

*对于任何* 订阅路径 S 包含具体实例编号（如 `Device.X.Y.1.`）和变化路径 C 包含不同实例编号（如 `Device.X.Y.2.Z`），NotifyListeners(C, value) 不应触发 S 的监听器。

**验证: 需求 1.3, 4.2**

### Property 4: 通配符匹配

*对于任何* 包含 `{i}` 占位符的订阅路径 S 和变化路径 C，如果将 S 中的 `{i}` 替换为 C 中对应位置的实例编号后 S 是 C 的前缀，则 NotifyListeners(C, value) 应触发 S 的监听器。

**验证: 需求 4.1**

### Property 5: 有效路径接受

*对于任何* 以 `Device.` 开头、不包含非法字符、长度不超过 256 字符、不包含连续点的路径 P，ValidatePath(P) 应返回 nil（无错误）。

**验证: 需求 2.1, 3.1, 3.2, 3.3**

### Property 6: 无效路径拒绝

*对于任何* 不以 `Device.` 开头、或包含非法字符、或长度超过 256 字符、或包含连续点的路径 P，ValidatePath(P) 应返回非 nil 错误。

**验证: 需求 2.2, 2.3, 2.4, 3.4, 3.5**

### Property 7: 错误位置报告

*对于任何* 包含非法字符的路径 P，ValidatePath(P) 返回的错误应包含非法字符的位置信息，且该位置应与实际非法字符位置一致。

**验证: 需求 5.1, 5.4**

## 错误处理

### 错误类型

```go
// PathValidationError 路径校验错误
type PathValidationError struct {
    Path     string // 原始路径
    Position int    // 错误位置（-1 表示整体错误）
    Reason   string // 错误原因
}

// 预定义错误原因
const (
    ErrReasonEmpty           = "路径不能为空"
    ErrReasonInvalidPrefix   = "路径必须以 Device. 开头"
    ErrReasonTooLong         = "路径长度超过 256 字符限制"
    ErrReasonIllegalChar     = "路径包含非法字符"
    ErrReasonConsecutiveDots = "路径包含连续的点"
    ErrReasonInvalidSegment  = "路径段格式无效"
)
```

### 错误处理策略

1. **ValidatePath**: 返回 `*PathValidationError`，包含具体错误位置和原因
2. **AddListener**: 先调用 ValidatePath，校验失败时返回错误，不添加监听器
3. **NotifyListeners**: 不返回错误，匹配失败时静默跳过

## 测试策略

### 单元测试

- 测试 PathValidator 的各种边界情况
- 测试 PathMatcher 的匹配逻辑
- 测试 ListenerManager 的集成行为

### 属性测试

使用 Go 的 `testing/quick` 包进行属性测试：

- 配置每个属性测试运行至少 100 次迭代
- 每个属性测试需标注对应的正确性属性：`**Feature: subscription-path-validation, Property {number}: {property_text}**`

### 测试生成器

```go
// 生成有效的 TR181 路径
func generateValidPath(rand *rand.Rand) string

// 生成无效路径（不以 Device. 开头）
func generateInvalidPrefixPath(rand *rand.Rand) string

// 生成包含非法字符的路径
func generatePathWithIllegalChar(rand *rand.Rand) string

// 生成超长路径
func generateTooLongPath(rand *rand.Rand) string

// 生成带通配符的订阅路径
func generateWildcardPath(rand *rand.Rand) string
```

## 实现细节

### 路径校验算法

```go
func (v *PathValidator) ValidatePath(path string) error {
    // 1. 检查空路径
    if path == "" {
        return &PathValidationError{Path: path, Position: -1, Reason: ErrReasonEmpty}
    }
    
    // 2. 检查长度
    if len(path) > MaxPathLength {
        return &PathValidationError{Path: path, Position: MaxPathLength, Reason: ErrReasonTooLong}
    }
    
    // 3. 检查前缀
    if !strings.HasPrefix(path, PathPrefix) {
        return &PathValidationError{Path: path, Position: 0, Reason: ErrReasonInvalidPrefix}
    }
    
    // 4. 检查非法字符
    for i, r := range path {
        if isIllegalChar(r) {
            return &PathValidationError{Path: path, Position: i, Reason: ErrReasonIllegalChar}
        }
    }
    
    // 5. 检查连续点
    if strings.Contains(path, "..") {
        pos := strings.Index(path, "..")
        return &PathValidationError{Path: path, Position: pos, Reason: ErrReasonConsecutiveDots}
    }
    
    return nil
}
```

### 路径匹配算法

```go
func (m *PathMatcher) Match(subscriptionPath, changedPath string) MatchResult {
    // 1. 精确匹配
    if subscriptionPath == changedPath {
        return MatchResult{Matched: true, MatchType: MatchTypeExact, MatchedPath: subscriptionPath}
    }
    
    // 2. 前缀匹配（订阅路径以 . 结尾）
    if strings.HasSuffix(subscriptionPath, ".") && strings.HasPrefix(changedPath, subscriptionPath) {
        return MatchResult{Matched: true, MatchType: MatchTypePrefix, MatchedPath: subscriptionPath}
    }
    
    // 3. 通配符匹配
    if strings.Contains(subscriptionPath, WildcardPlaceholder) {
        if m.matchWildcard(subscriptionPath, changedPath) {
            return MatchResult{Matched: true, MatchType: MatchTypeWildcard, MatchedPath: subscriptionPath}
        }
    }
    
    return MatchResult{Matched: false, MatchType: MatchTypeNone}
}

func (m *PathMatcher) matchWildcard(pattern, path string) bool {
    // 将 {i} 替换为正则表达式 \d+
    regexPattern := strings.ReplaceAll(regexp.QuoteMeta(pattern), `\{i\}`, `\d+`)
    if strings.HasSuffix(pattern, ".") {
        // 前缀匹配模式
        regexPattern = "^" + regexPattern
    } else {
        // 精确匹配模式
        regexPattern = "^" + regexPattern + "$"
    }
    matched, _ := regexp.MatchString(regexPattern, path)
    return matched
}
```

### NotifyListeners 优化实现

```go
func (lm *ListenerManager) NotifyListeners(paramName string, value interface{}) {
    var matchedListeners []struct {
        listener  tr181Model.Listener
        matchType MatchType
    }
    
    // 遍历所有订阅，找出匹配的监听器
    for subPath, listeners := range lm.TR181DataModel.Listeners {
        result := lm.matcher.Match(subPath, paramName)
        if result.Matched {
            for _, listener := range listeners {
                matchedListeners = append(matchedListeners, struct {
                    listener  tr181Model.Listener
                    matchType MatchType
                }{listener, result.MatchType})
            }
        }
    }
    
    // 按匹配类型排序：精确匹配 > 前缀匹配 > 通配符匹配
    sort.Slice(matchedListeners, func(i, j int) bool {
        return matchedListeners[i].matchType < matchedListeners[j].matchType
    })
    
    // 触发监听器
    for _, ml := range matchedListeners {
        go ml.listener.Listener(ml.listener.SubscriptionId, value)
    }
}
```

## 文件组织

```
client/
├── repository/
│   ├── listener_manager.go      # 更新 NotifyListeners 实现
│   └── path_validator.go        # 新增：路径校验器
├── model/
│   └── path.go                  # 新增：PathInfo, PathValidationError 等类型
└── usecase/
    └── subscription.go          # 更新：添加订阅时调用路径校验
```

