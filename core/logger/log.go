package logger

import (
	"context"
	"log/slog"
)

// 统一日志调用接口
// 所有项目应 import "github.com/baowk/dilu-core/core/logger" 而非直接 import "log/slog"
// 未来如需替换日志库（zap、zerolog 等），只需修改此文件

// Debug 记录 Debug 级别日志
func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// Info 记录 Info 级别日志
func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

// Warn 记录 Warn 级别日志
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

// Error 记录 Error 级别日志
func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// DebugContext 带 context 的 Debug 日志（用于链路追踪等场景）
func DebugContext(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, msg, args...)
}

// InfoContext 带 context 的 Info 日志
func InfoContext(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

// WarnContext 带 context 的 Warn 日志
func WarnContext(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

// ErrorContext 带 context 的 Error 日志
func ErrorContext(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

// With 创建带固定字段的子 logger（如 "module"="auth"）
func With(args ...any) *slog.Logger {
	return slog.Default().With(args...)
}

// Default 返回当前默认 logger 实例
func Default() *slog.Logger {
	return slog.Default()
}
