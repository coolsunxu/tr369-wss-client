# 实现计划

- [x] 1. 创建 helpers.go 辅助函数文件




  - [ ] 1.1 创建 helpers.go 文件并添加包声明和导入
    - 添加 package usecase 声明

    - 添加必要的 import 语句
    - _需求: 5.3_
  - [ ] 1.2 迁移辅助函数到 helpers.go
    - 迁移 extractParamSettings 泛型函数
    - 迁移 constructGetResp 方法




    - 迁移 isExistPath 方法
    - 迁移 getNewInstance 方法

    - _需求: 5.3_

- [ ] 2. 创建 notification.go 通知发送文件
  - [ ] 2.1 创建 notification.go 文件并添加包声明和导入
    - 添加 package usecase 声明
    - 添加必要的 import 语句
    - _需求: 3.1, 3.3_
  - [x] 2.2 迁移通知发送函数到 notification.go




    - 迁移 sendNotification 方法


    - 迁移 HandleValueChange 方法
    - 迁移 HandleObjectCreation 方法
    - 迁移 HandleObjectDeletion 方法
    - 迁移 notifyValueChange 方法
    - 迁移 notifyObjectCreation 方法
    - 迁移 notifyObjectDeletion 方法
    - _需求: 3.1, 3.2, 3.3_

- [ ] 3. 创建 subscription.go 订阅管理文件
  - [ ] 3.1 创建 subscription.go 文件并添加包声明和导入
    - 添加 package usecase 声明
    - 添加必要的 import 语句（包括 regexp）



    - _需求: 2.1, 2.3_

  - [x] 3.2 迁移订阅管理函数到 subscription.go

    - 迁移 subscriptionInstanceRegex 正则表达式变量
    - 迁移 HandleAddLocalAgentSubscription 方法
    - 迁移 HandleDeleteLocalAgentSubscription 方法
    - 迁移 HandleSubscription 方法
    - 迁移 isSubscriptionPath 方法
    - 迁移 isSubscriptionParentPath 方法
    - 迁移 isSubscriptionInstancePath 方法
    - 迁移 extractSubscriptionParams 方法
    - 迁移 deleteAllSubscriptions 方法
    - 迁移 deleteSingleSubscription 方法




    - 迁移 getSubscriptionReferenceList 方法
    - _需求: 2.1, 2.2, 2.3_

- [ ] 4. 创建 message_handler.go 消息处理文件
  - [x] 4.1 创建 message_handler.go 文件并添加包声明和导入

    - 添加 package usecase 声明
    - 添加必要的 import 语句
    - _需求: 1.1, 1.3_



  - [ ] 4.2 迁移消息处理函数到 message_handler.go
    - 迁移 HandleGetRequest 方法
    - 迁移 HandleSetRequest 方法
    - 迁移 HandleAddRequest 方法
    - 迁移 HandleDeleteRequest 方法
    - 迁移 HandleOperateRequest 方法
    - 迁移 HandleNotifyResp 方法
    - 迁移 SendOperateCompleteNotify 方法
    - 迁移 handleObjectCreationSideEffects 方法
    - 迁移 handleObjectDeletionPreEffects 方法
    - _需求: 1.1, 1.2, 1.3_

- [ ] 5. 精简主文件 client.go
  - [ ] 5.1 移除已迁移的函数，保留核心内容
    - 保留 ClientUseCase 结构体定义
    - 保留 NewClientUseCase 构造函数
    - 保留 HandleMessage 消息分发入口
    - 保留 HandleMTPMsgTransmit 消息传输方法
    - 移除所有已迁移到其他文件的函数
    - _需求: 5.1, 5.2_
  - [ ] 5.2 清理 import 语句
    - 移除不再需要的 import
    - 保留核心功能所需的 import
    - _需求: 4.3_

- [ ] 6. 检查点 - 确保编译通过
  - 确保所有测试通过，如有问题请询问用户。

- [ ]* 6.1 编写属性测试：拆分前后行为等价性
  - **Property 1: 拆分前后行为等价性**
  - **验证: 需求 4.1, 4.2**

