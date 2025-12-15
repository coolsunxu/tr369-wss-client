# 实现计划

- [x] 1. 定义新接口

  - [x] 1.1 在 client/model/client.go 中定义 DataRepository 接口


    - 包含 GetValue、GetParameters、SetValue、DeleteNode、Start 方法
    - _需求: 1.1, 4.1, 4.2_
  - [x] 1.2 在 client/model/client.go 中定义 ListenerManager 接口

    - 包含 AddListener、RemoveListener、ResetListener、NotifyListeners 方法
    - _需求: 2.1, 4.1, 4.2_

- [x] 2. 更新 Repository 实现

  - [x] 2.1 重命名 HandleDeleteRequest 为 DeleteNode


    - 更新方法签名和实现
    - _需求: 4.2_
  - [x] 2.2 添加 GetParameters 方法

    - 返回 TR181DataModel.Parameters
    - _需求: 1.1_
  - [x] 2.3 移除 ConstructGetResp、IsExistPath、GetNewInstance 方法

    - 这些方法将移至 UseCase 层
    - _需求: 1.1_
  - [ ]* 2.4 编写属性测试：数据操作行为等价性
    - **Property 1: 行为等价性**
    - **验证: 需求 3.1, 3.3**

- [x] 3. 更新 UseCase 层
  - [x] 3.1 修改 ClientUseCase 结构体


    - 将 ClientRepository 拆分为 DataRepo 和 ListenerMgr 两个字段
    - _需求: 5.1_
  - [x] 3.2 更新 NewClientUseCase 构造函数

    - 接受 DataRepository 和 ListenerManager 两个参数
    - _需求: 5.1, 5.2_
  - [x] 3.3 添加 constructGetResp 方法

    - 从 Repository 移至 UseCase
    - 调用 trtree.ConstructGetResp
    - _需求: 1.1_
  - [x] 3.4 添加 isExistPath 方法

    - 从 Repository 移至 UseCase
    - 调用 trtree.IsExistPath
    - _需求: 1.1_

  - [x] 3.5 添加 getNewInstance 方法
    - 从 Repository 移至 UseCase
    - 调用 trtree.GetNewInstance
    - _需求: 1.1_

  - [x] 3.6 更新所有调用点
    - HandleGetRequest 使用 constructGetResp
    - HandleSetRequest 使用 isExistPath
    - HandleAddRequest 使用 getNewInstance
    - HandleDeleteRequest 使用 DataRepo.DeleteNode
    - _需求: 3.1, 3.3_
  - [ ]* 3.7 编写属性测试：监听器管理行为等价性
    - **Property 2: 监听器管理等价性**
    - **验证: 需求 3.1**

- [x] 4. 更新依赖注入

  - [x] 4.1 更新 client/client.go 中的初始化代码


    - 创建 clientRepository 实例
    - 将同一实例作为 DataRepository 和 ListenerManager 传入 UseCase
    - _需求: 5.1, 5.3_
  - [x] 4.2 更新 main.go 中的初始化代码（如需要）

    - 确保依赖注入正确
    - _需求: 3.2_

- [x] 5. 检查点 - 确保所有测试通过

  - 确保所有测试通过，如有问题请询问用户。

- [x] 6. 清理和验证
  - [x] 6.1 移除 model/client.go 中旧接口的冗余方法声明

    - 确保接口定义简洁
    - _需求: 4.3_

  - [x] 6.2 验证编译通过

    - 运行 go build 确保无编译错误
    - _需求: 3.1_
  - [x] 6.3 验证功能正常

    - 确保 USP 消息处理流程正常工作
    - _需求: 3.3_

- [x] 7. 最终检查点 - 确保所有测试通过


  - 确保所有测试通过，如有问题请询问用户。
