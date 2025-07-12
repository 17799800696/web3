package config

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DatabaseConfig 数据库配置结构
type DatabaseConfig struct {
	Host      string
	Port      int
	User      string
	Password  string
	DBName    string
	Charset   string
	ParseTime bool
	Loc       string
}

// 默认数据库配置
var defaultConfig = DatabaseConfig{
	Host:      "127.0.0.1",
	Port:      3306,
	User:      "root",
	Password:  "123456",
	DBName:    "test",
	Charset:   "utf8mb4",
	ParseTime: true,
	Loc:       "Local",
}

// 全局数据库连接实例
var (
	db   *gorm.DB
	once sync.Once
)

// GetDSN 获取数据库连接字符串
func GetDSN() string {
	config := GetDatabaseConfig()
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.Charset,
		config.ParseTime,
		config.Loc,
	)
}

// GetDatabaseConfig 获取数据库配置
func GetDatabaseConfig() DatabaseConfig {
	return defaultConfig
}

// SetDatabaseConfig 设置数据库配置
func SetDatabaseConfig(config DatabaseConfig) {
	defaultConfig = config
}

// GetDB 获取数据库连接实例（单例模式）
func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		db, err = gorm.Open(mysql.Open(GetDSN()), &gorm.Config{})
		if err != nil {
			log.Fatal("数据库连接失败:", err)
		}
		log.Println("数据库连接成功")
	})
	return db
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}
