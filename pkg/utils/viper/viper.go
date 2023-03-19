package viper

import (
	v "github.com/spf13/viper"
	"log"
)

// Config 公有变量,获取Viper
type Config struct {
	Viper *v.Viper
}

// Init 初始化Viper配置
func Init(ConfigName string) Config {
	config := Config{Viper: v.New()}
	viper := config.Viper
	viper.SetConfigType("yaml")     //设置配置文件类型
	viper.SetConfigName(ConfigName) //设置配置文件名
	viper.AddConfigPath("./config") //设置配置文件路径
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")
	//读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("errno is %+v", err)
	}
	return config
}
