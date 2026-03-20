package logger

import (
	"context"

	"github.com/rs/zerolog"
)

// 所有项目应 import "github.com/baowk/dilu-core/core/logger" 而非直接 import zerolog
// 未来如需替换日志库，只需修改此文件

func Debug() *zerolog.Event { return Log.Debug() }
func Info() *zerolog.Event  { return Log.Info() }
func Warn() *zerolog.Event  { return Log.Warn() }
func Error() *zerolog.Event { return Log.Error() }
func Fatal() *zerolog.Event { return Log.Fatal() }

// Ctx 从 context 中取出 logger（用于链路追踪）
func Ctx(ctx context.Context) *zerolog.Logger { return zerolog.Ctx(ctx) }

// With 创建带固定字段的子 logger（如 logger.With().Str("module","auth").Logger()）
func With() zerolog.Context { return Log.With() }

// Default 返回全局 logger 实例
func Default() zerolog.Logger { return Log }
