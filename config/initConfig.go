package config

import (
	"gin-study/model" // 这里的路径必须和你go.mod里的module名完全一致

	"github.com/spf13/viper"
)

// 这是解决硬编码方面的东西的
// 1. 用到model包：定义全局配置变量，类型是model.AppConfig
var GlobalConfig model.AppConfig

// 2. 用到viper包：写InitConfig函数，调用viper的方法
func InitConfig() error {
	// 这些都是viper的方法，用到了导入的viper包
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// 把读取到的配置绑定到GlobalConfig（用到了model包）
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return err
	}

	return nil
}
