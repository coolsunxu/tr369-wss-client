package config

import (
	logger "tr369-wss-client/log"

	"github.com/spf13/viper"
)

func InitConfig(configPath string) error {

	// 设置文件路径
	viper.SetConfigFile(configPath)

	// 设置文件格式
	viper.SetConfigType("json")

	// 读取文件内容
	if err := viper.ReadInConfig(); err != nil {
		logger.Warnf("read config file %s err: %s ", configPath, err)
	}

	// 反序列化参数到全局变量中
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		logger.Fatalf("unmarshal error config file %s err: %s ", configPath, err)
		return err
	}

	return nil
}
