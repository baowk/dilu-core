# dilu-core

[English](README_en.md) 

dilu基础包，dilu必要的基础配置和封装

1.0.0版本之前，系统日志为zap日志，1.0.0版本开始，系统日志为标准库log/slog日志

## 功能特性

- 🚀 **核心配置管理** - 提供统一的应用配置结构
- 🔧 **基础封装** - 对常用组件进行封装，简化使用
- 📝 **日志系统** - 从zap迁移到标准库slog（1.0.0+版本）
- ⚡ **高性能** - 优化的核心组件实现
- 🛠️ **易扩展** - 模块化设计，便于功能扩展

## 目录结构

```
dilu-core/
├── common/          # 通用工具和常量
│   ├── consts/      # 常量定义
│   └── utils/       # 工具函数
├── config/          # 配置管理
│   ├── config.go    # 主配置结构
│   ├── db.go        # 数据库配置
│   ├── log.go       # 日志配置
│   ├── redis.go     # Redis配置
│   └── ...          # 其他配置模块
├── core/            # 核心组件
│   ├── base/        # 基础结构
│   ├── cache/       # 缓存接口
│   ├── errs/        # 错误处理
│   ├── logger/      # 日志系统
│   ├── application.go # 应用入口
│   ├── db.go        # 数据库初始化
│   └── version.go   # 版本管理
└── resources/       # 资源文件
    └── config.dev.yaml # 开发环境配置示例
```

## 快速开始

### 安装

```bash
go get github.com/baowk/dilu-core
```

### 基本使用

```go
package main

import (
    "github.com/baowk/dilu-core/core"
    "github.com/baowk/dilu-core/config"
)

func main() {
    // 初始化应用
    core.Init()
    
    // 获取配置
    cfg := core.Cfg
    
    // 使用配置
    println("Server Name:", cfg.Server.Name)
    println("Server Port:", cfg.Server.Port)
}
```

### 配置结构

主要配置结构包括：

- `AppCfg` - 应用主配置
- `ServerCfg` - 服务器配置
- `DBCfg` - 数据库配置
- `LogCfg` - 日志配置
- `CacheCfg` - 缓存配置
- `JWT` - JWT配置

## 版本变更

### 1.0.0 (最新版本)
- 🔄 将日志系统从zap迁移至标准库log/slog
- 🐛 修复已知问题
- 📈 性能优化

### 0.x.x (历史版本)
- 🚀 初始版本发布
- 📦 基础功能实现
- 🔧 核心组件封装

## 贡献指南

欢迎提交Issue和Pull Request来帮助改进项目。

### 开发环境设置

1. 克隆仓库
```bash
git clone https://github.com/baowk/dilu-core.git
cd dilu-core
```

2. 安装依赖
```bash
go mod tidy
```

3. 运行测试
```bash
go test ./...
```

## 许可证

MIT License