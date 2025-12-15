# 实现HandleDeleteLocalAgentSubscription函数逻辑（修正版）

## 目标

补充`HandleDeleteLocalAgentSubscription`函数的逻辑，实现以下功能：

1. 判断`requestPath`是否是`tr181Model.DeviceLocalAgentSubscription`的父节点（包括各级父节点）
2. 判断`requestPath`是否符合`Device.LocalAgent.Subscription.*.`的形式，其中`*`只能为正整数

## 实现方案

### 1. 父节点判断

* `tr181Model.DeviceLocalAgentSubscription`的值为`"Device.LocalAgent.Subscription."`

* 其各级父节点包括：

  * `"Device.LocalAgent."`

  * `"Device."`

  * 根节点（空字符串，这里不考虑）

* 使用字符串前缀匹配或直接比较判断

### 2. 路径格式判断

* 使用正则表达式`^Device\.LocalAgent\.Subscription\.([1-9]\d*)\.$`

* 该正则表达式确保：

  * 以`Device.LocalAgent.Subscription.`开头

  * 后跟一个正整数（`[1-9]\d*`）

  * 以`.`结尾

### 3. 函数实现

```go
func (uc *ClientUseCase) HandleDeleteLocalAgentSubscription(requestPath string) error {
    // 为空则直接返回
    if requestPath == "" {
        return nil
    }

    // 判断是否是各级父节点
    switch requestPath {
    case "Device.":
        // 处理Device.父节点逻辑
        logger.Infof("Handle delete subscription for parent path: %s", requestPath)
        return nil
    case "Device.LocalAgent.":
        // 处理Device.LocalAgent.父节点逻辑
        logger.Infof("Handle delete subscription for parent path: %s", requestPath)
        return nil
    }

    // 判断是否符合 Device.LocalAgent.Subscription.*. 格式，*为正整数
    matched, err := regexp.MatchString(`^Device\.LocalAgent\.Subscription\.([1-9]\d*)\.$`, requestPath)
    if err != nil {
        return err
    }
    
    if matched {
        // 符合格式，处理删除逻辑
        logger.Infof("Handle delete subscription for path: %s", requestPath)
        // 这里需要添加实际的删除订阅逻辑，比如调用repository的RemoveListener
    }

    return nil
}
```

### 4. 注意事项

* 需要导入`regexp`包

* 需要添加实际的删除订阅逻辑，比如调用repository的RemoveListener方法

* 需要添加必要的日志记录

## 预期效果

* 当`requestPath`是`Device.`或`Device.LocalAgent.`时，函数能正确识别为父节点

* 当`requestPath`符合`Device.LocalAgent.Subscription.*.`格式且`*`为正整数时，函数能正确识别

* 其他情况下，函数不做处理

## 实现步骤

1. 修改`HandleDeleteLocalAgentSubscription`函数，添加上述逻辑
2. 确保导入`regexp`包
3. 添加必要的日志记录
4. 实现实际的删除订阅逻辑

