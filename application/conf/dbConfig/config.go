package dbConfig

import (
	"github.com/spf13/viper"
	"log"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

var DBConfig DatabaseConfig

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("conf/dbConfig")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Mysql读取配置错误: %v", err)
	}

	// UnmarshalKey()将配置中的"database"部分解析到dbConfig.DBConfig结构体中
	if err := viper.UnmarshalKey("database", &DBConfig); err != nil {
		log.Fatalf("解析数据库配置错误: %v", err)
	}
}
