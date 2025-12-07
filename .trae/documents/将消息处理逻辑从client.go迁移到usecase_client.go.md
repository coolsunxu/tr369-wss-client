1. 修改client.go：
   - 导入client/usecase包
   - 在WSClient结构体中添加clientUseCase字段
   - 修改NewWSClient函数，初始化clientUseCase
   - 修改messageHandler方法，将消息处理部分调用clientUseCase.HandleMessage
   - 移除或调整原有的handleProtobufMessage等方法

2. 确保usecase/client.go中的方法能够正确处理消息，并返回适当的响应

3. 验证修改后的代码能够正常编译和运行

这个计划将把消息处理的核心逻辑从client.go迁移到usecase/client.go，实现更好的代码分离和分层设计。