// Package config 提供基于 viper 的配置加载功能
package config

import (
	"tr369-wss-client/internal/infrastructure/logging"

	"github.com/spf13/viper"
)

// GlobalConfig 全局配置实例
var GlobalConfig = Config{
	DataRefreshConfig: &DataRefreshConfig{},
	WebsocketConfig:   &WebsocketConfig{},
}

// InitConfig 使用 viper 初始化配置
func InitConfig(configPath string, logger logging.Logger) error {
	// 设置文件路径
	viper.SetConfigFile(configPath)

	// 设置文件格式
	viper.SetConfigType("json")

	// 读取文件内容
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("读取配置文件失败 %s: %s", configPath, err)
		return err
	}

	// 反序列化参数到全局变量中
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		logger.Fatal("反序列化配置失败 %s: %s", configPath, err)
		return err
	}

	// 验证配置
	if err := GlobalConfig.Validate(); err != nil {
		logger.Error("配置验证失败: %s", err)
		return err
	}

	return nil
}
