# 实现计划

- [x] 1. 创建路径相关的数据模型和类型定义




  - [x] 1.1 在 `client/model/path.go` 中定义 PathInfo 结构体

    - 包含 FullPath, Segments, IsObject, IsParameter, HasWildcard, InstanceNums 字段
    - _需求: 3.1, 3.2, 3.3_
  - [x] 1.2 在 `client/model/path.go` 中定义 PathValidationError 错误类型

    - 包含 Path, Position, Reason 字段
    - 实现 Error() 方法
    - _需求: 5.1, 5.2, 5.3, 5.4_
  - [ ] 1.3 在 `client/model/path.go` 中定义常量和错误原因
    - MaxPathLength = 256




    - PathPrefix = "Device."
    - WildcardPlaceholder = "{i}"
    - 预定义错误原因常量
    - _需求: 2.4, 2.1_

- [ ] 2. 实现 PathValidator 路径校验器
  - [ ] 2.1 在 `client/repository/path_validator.go` 中创建 PathValidator 结构体
    - 实现 ValidatePath(path string) error 方法
    - 校验空路径、长度、前缀、非法字符、连续点
    - _需求: 2.1, 2.2, 2.3, 2.4, 3.4, 3.5, 5.1, 5.2, 5.3, 5.4_
  - [x]* 2.2 编写属性测试：有效路径接受

    - **Property 5: 有效路径接受**
    - **验证: 需求 2.1, 3.1, 3.2, 3.3**

  - [ ]* 2.3 编写属性测试：无效路径拒绝
    - **Property 6: 无效路径拒绝**

    - **验证: 需求 2.2, 2.3, 2.4, 3.4, 3.5**
  - [ ]* 2.4 编写属性测试：错误位置报告
    - **Property 7: 错误位置报告**




    - **验证: 需求 5.1, 5.4**

  - [ ] 2.5 实现 IsObjectPath(path string) bool 方法
    - 判断路径是否以 `.` 结尾
    - _需求: 3.1_
  - [ ] 2.6 实现 IsParameterPath(path string) bool 方法
    - 判断路径是否不以 `.` 结尾
    - _需求: 3.2_
  - [ ] 2.7 实现 ParsePath(path string) (*PathInfo, error) 方法
    - 解析路径为 PathInfo 结构
    - 识别实例编号和通配符
    - _需求: 3.1, 3.2, 3.3_


- [ ] 3. 实现 PathMatcher 路径匹配器
  - [ ] 3.1 在 `client/repository/path_matcher.go` 中创建 PathMatcher 结构体和 MatchResult 类型
    - 定义 MatchType 枚举：None, Exact, Prefix, Wildcard
    - _需求: 1.1, 1.2, 4.1_
  - [x] 3.2 实现 Match(subscriptionPath, changedPath string) MatchResult 方法

    - 按优先级检查：精确匹配 → 前缀匹配 → 通配符匹配
    - _需求: 1.1, 1.2, 1.3, 4.1, 4.3_



  - [-]* 3.3 编写属性测试：前缀匹配正确性

    - **Property 1: 前缀匹配正确性**

    - **验证: 需求 1.1, 1.4**
  - [ ]* 3.4 编写属性测试：精确匹配优先
    - **Property 2: 精确匹配优先**

    - **验证: 需求 1.2, 4.3**
  - [ ]* 3.5 编写属性测试：实例隔离
    - **Property 3: 实例隔离**
    - **验证: 需求 1.3, 4.2**
  - [ ] 3.6 实现 matchWildcard(pattern, path string) bool 私有方法
    - 将 {i} 替换为正则表达式 \d+ 进行匹配
    - _需求: 4.1_
  - [x]* 3.7 编写属性测试：通配符匹配



    - **Property 4: 通配符匹配**

    - **验证: 需求 4.1**

- [ ] 4. 检查点 - 确保所有测试通过
  - 确保所有测试通过，如有问题请询问用户。



- [ ] 5. 更新 ListenerManager 集成校验和匹配逻辑
  - [ ] 5.1 更新 `client/repository/listener_manager.go` 中的 ListenerManager 结构体
    - 添加 PathValidator 和 PathMatcher 字段
    - 更新 NewListenerManager 构造函数
    - _需求: 2.1_
  - [ ] 5.2 更新 AddListener 方法
    - 在添加监听器前调用 PathValidator.ValidatePath 校验路径
    - 校验失败时返回错误
    - _需求: 2.1, 2.2, 2.3, 2.4_
  - [ ] 5.3 重构 NotifyListeners 方法
    - 遍历所有订阅，使用 PathMatcher.Match 进行匹配
    - 按匹配类型排序：精确匹配 > 前缀匹配 > 通配符匹配
    - 触发所有匹配的监听器
    - _需求: 1.1, 1.2, 1.3, 1.4, 4.1, 4.2, 4.3_
  - [ ]* 5.4 编写 ListenerManager 集成测试
    - 测试 AddListener 路径校验
    - 测试 NotifyListeners 匹配逻辑
    - _需求: 1.1, 1.2, 1.3, 2.1, 2.2_

- [ ] 6. 更新 UseCase 层订阅处理
  - [ ] 6.1 更新 `client/usecase/subscription.go` 中的 extractSubscriptionParams 方法
    - 调用 PathValidator 校验 ReferenceList 参数
    - 返回详细的校验错误信息
    - _需求: 2.1, 2.2, 2.3, 2.4, 5.1_
  - [ ]* 6.2 编写 extractSubscriptionParams 单元测试
    - 测试有效和无效 ReferenceList 的处理
    - _需求: 2.1, 2.2_

- [ ] 7. 最终检查点 - 确保所有测试通过
  - 确保所有测试通过，如有问题请询问用户。

