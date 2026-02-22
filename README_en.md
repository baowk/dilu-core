# dilu-core

dilu core package, essential basic configurations and encapsulations for dilu

System logs were using zap logger before version 1.0.0, starting from version 1.0.0, system logs use the standard library log/slog logger

## Features

- 🚀 **Core Configuration Management** - Unified application configuration structure
- 🔧 **Basic Encapsulation** - Encapsulation of commonly used components for simplified usage
- 📝 **Logging System** - Migration from zap to standard library slog (version 1.0.0+)
- ⚡ **High Performance** - Optimized core component implementation
- 🛠️ **Extensible** - Modular design for easy feature extension

## Directory Structure

```
dilu-core/
├── common/          # Common utilities and constants
│   ├── consts/      # Constant definitions
│   └── utils/       # Utility functions
├── config/          # Configuration management
│   ├── config.go    # Main configuration structure
│   ├── db.go        # Database configuration
│   ├── log.go       # Logging configuration
│   ├── redis.go     # Redis configuration
│   └── ...          # Other configuration modules
├── core/            # Core components
│   ├── base/        # Base structures
│   ├── cache/       # Cache interfaces
│   ├── errs/        # Error handling
│   ├── logger/      # Logging system
│   ├── application.go # Application entry
│   ├── db.go        # Database initialization
│   └── version.go   # Version management
└── resources/       # Resource files
    └── config.dev.yaml # Development environment configuration example
```

## Quick Start

### Installation

```bash
go get github.com/baowk/dilu-core
```

### Basic Usage

```go
package main

import (
    "github.com/baowk/dilu-core/core"
    "github.com/baowk/dilu-core/config"
)

func main() {
    // Initialize application
    core.Init()
    
    // Get configuration
    cfg := core.Cfg
    
    // Use configuration
    println("Server Name:", cfg.Server.Name)
    println("Server Port:", cfg.Server.Port)
}
```

### Configuration Structure

Main configuration structures include:

- `AppCfg` - Main application configuration
- `ServerCfg` - Server configuration
- `DBCfg` - Database configuration
- `LogCfg` - Logging configuration
- `CacheCfg` - Cache configuration
- `JWT` - JWT configuration

## Version Changes

### 1.0.0 (Latest)
- 🔄 Migrated logging system from zap to standard library log/slog
- 🐛 Fixed known issues
- 📈 Performance optimizations

### 0.x.x (Historical)
- 🚀 Initial release
- 📦 Basic functionality implementation
- 🔧 Core component encapsulation

## Contribution Guide

Welcome to submit Issues and Pull Requests to help improve the project.

### Development Environment Setup

1. Clone repository
```bash
git clone https://github.com/baowk/dilu-core.git
cd dilu-core
```

2. Install dependencies
```bash
go mod tidy
```

3. Run tests
```bash
go test ./...
```

## License

MIT License