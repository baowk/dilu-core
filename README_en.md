# dilu-core

[中文](README.md) 

dilu core package, essential basic configurations and encapsulations for dilu

System logs were using zap logger before version 1.0.0, starting from version 1.0.0, system logs use the standard library log/slog logger

## Features

- 🚀 **Core Configuration Management** - Unified application configuration structure
- 🔧 **Basic Encapsulation** - Encapsulation of commonly used components for simplified usage
- 📝 **Logging System** - Migration from zap to standard library slog (version 1.0.0+)
- ⚡ **High Performance** - Optimized core component implementation
- 🛠️ **Extensible** - Modular design for easy feature extension

### Configuration Structure

Main configuration structures include:

- `AppCfg` - Main application configuration
- `ServerCfg` - Server configuration
- `DBCfg` - Database configuration
- `LogCfg` - Logging configuration
- `CacheCfg` - Cache configuration
- `JWT` - JWT configuration

Welcome to submit Issues and Pull Requests to help improve the project.

## License

MIT License