# 全局数据库配置使用指南

## 概述

这个项目提供了一个全局数据库配置系统，让你可以在项目的任何地方轻松使用数据库连接，而不需要重复配置连接字符串。

## 配置结构

全局配置位于 `config/database.go` 文件中，包含以下功能：

### 1. 默认配置
```go
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
```

### 2. 主要函数

#### 获取数据库连接
```go
db := config.GetDB()
```

#### 获取连接字符串
```go
dsn := config.GetDSN()
```

#### 设置自定义配置
```go
config.SetDatabaseConfig(DatabaseConfig{
    Host:     "localhost",
    Port:     3306,
    User:     "myuser",
    Password: "mypassword",
    DBName:   "mydatabase",
})
```

#### 关闭数据库连接
```go
config.CloseDB()
```

## 使用方法

### 1. 基本使用

在任何Go文件中，你只需要导入config包并使用GetDB()函数：

```go
package main

import (
    "log"
    "github.com/test/init_project/config"
)

func main() {
    // 获取数据库连接
    db := config.GetDB()
    
    // 使用数据库连接
    // ... 你的数据库操作代码
}
```

### 2. 完整示例

参考 `examples/global_config_example.go` 文件，展示了如何使用全局配置进行转账操作。

### 3. 修改配置

如果你需要修改数据库配置，可以在程序启动时调用：

```go
config.SetDatabaseConfig(config.DatabaseConfig{
    Host:     "your-host",
    Port:     3306,
    User:     "your-user",
    Password: "your-password",
    DBName:   "your-database",
})
```

## 优势

1. **统一配置**: 所有数据库连接使用相同的配置
2. **单例模式**: 确保只有一个数据库连接实例
3. **易于维护**: 修改配置只需要在一个地方
4. **类型安全**: 使用结构体而不是字符串配置
5. **自动连接管理**: 自动处理连接的创建和关闭

## 注意事项

1. 确保在使用前已经正确配置了数据库连接参数
2. 在程序结束时可以调用 `config.CloseDB()` 关闭连接
3. 配置是全局的，修改会影响所有使用该配置的地方
4. 使用单例模式，确保线程安全

## 运行示例

```bash
# 运行示例程序
go run examples/global_config_example.go
```

这样你就可以在项目的任何地方使用 `config.GetDB()` 来获取数据库连接，而不需要重复配置连接字符串了！ 