package dao

import (
	"CipProject/application/conf/dbConfig"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var _DB *gorm.DB

func init() {
	// 从配置加载连接
	cfg := dbConfig.DBConfig

	// 连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName)

	// 建立连接
	var err error
	_DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 配置日志，也可通过Db.Debug()模式开启日志
	})
	if err != nil {
		log.Fatalln("db connected error", err)
	}

	// 使用连接池
	db, err := _DB.DB()     // 获得连接池
	db.SetMaxOpenConns(100) // 设置连接池最大连接数100
	db.SetMaxIdleConns(20)  // 设置连接池最大空闲连接数20
	if err != nil {
		log.Fatalln("db connected error", err)
	}
}

func GetDB() *gorm.DB {
	return _DB
}
