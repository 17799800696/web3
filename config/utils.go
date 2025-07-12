package config

import (
	"log"

	"gorm.io/gorm"
)

// GetDBWithLog 获取数据库连接并记录日志
func GetDBWithLog() *gorm.DB {
	db := GetDB()
	log.Println("使用全局配置获取数据库连接成功")
	return db
}

// GetDSNWithLog 获取数据库连接字符串并记录日志
func GetDSNWithLog() string {
	dsn := GetDSN()
	log.Printf("当前数据库连接字符串: %s", dsn)
	return dsn
}

// PrintConfig 打印当前数据库配置
func PrintConfig() {
	config := GetDatabaseConfig()
	log.Printf("当前数据库配置: Host=%s, Port=%d, User=%s, DBName=%s",
		config.Host, config.Port, config.User, config.DBName)
}
