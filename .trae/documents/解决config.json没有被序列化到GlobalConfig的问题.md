1. 修改`main.go`中的配置文件路径，从`./config/config.json`改为`./config.json`
2. 确保`config.json`文件存在于根目录下
3. 验证配置加载逻辑是否正确

这个修改将使程序能够读取根目录下的`config.json`文件，并将其内容正确序列化到GlobalConfig中。