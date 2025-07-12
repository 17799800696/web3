# 全局数据库配置使用指南

## 🎯 问题解决

你之前遇到的问题是在多个文件中重复配置数据库连接字符串：
```go
dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
```

## ✅ 解决方案

现在你可以在项目的任何地方使用全局配置：

### 1. 基本使用
```go
import "github.com/test/init_project/config"

func yourFunction() {
    // 获取数据库连接
    db := config.GetDB()
    
    // 直接使用，无需重复配置
    // ... 你的数据库操作
}
```

### 2. 完整示例
参考 `examples/global_config_example.go` 文件，展示了完整的转账操作。

### 3. 配置修改
如果需要修改配置，在程序启动时调用：
```go
config.SetDatabaseConfig(config.DatabaseConfig{
    Host:     "your-host",
    Port:     3306,
    User:     "your-user", 
    Password: "your-password",
    DBName:   "your-database",
})
```

## 🚀 优势

1. **统一管理**: 所有数据库连接使用相同配置
2. **单例模式**: 确保只有一个连接实例
3. **易于维护**: 修改配置只需一个地方
4. **类型安全**: 使用结构体而非字符串
5. **自动管理**: 自动处理连接创建和关闭

## 📁 文件结构

```
config/
├── database.go    # 全局配置核心文件
└── utils.go       # 工具函数

examples/
└── global_config_example.go  # 使用示例

README_global_config.md       # 详细文档
```

## 🧪 测试

运行示例程序验证配置：
```bash
go run examples/global_config_example.go
```

现在你可以在项目的任何地方使用 `config.GetDB()` 来获取数据库连接，无需重复配置连接字符串！ 