package logger

import (
	"context"
	"log/slog"
	"sync"
)

// LevelHandler 将不同级别的日志分发到对应的 slog.Handler
type LevelHandler struct {
	handlers map[slog.Level][]slog.Handler
	mu       sync.RWMutex
}

func NewLevelHandler() *LevelHandler {
	return &LevelHandler{
		handlers: make(map[slog.Level][]slog.Handler),
	}
}

// AddHandler 添加级别对应的处理器（仅在初始化阶段调用）
func (h *LevelHandler) AddHandler(level slog.Level, handler slog.Handler) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[level] = append(h.handlers[level], handler)
}

// Enabled 检查当前级别是否有对应的处理器
func (h *LevelHandler) Enabled(_ context.Context, level slog.Level) bool {
	h.mu.RLock()
	_, ok := h.handlers[level]
	h.mu.RUnlock()
	return ok
}

// Handle 分发日志记录到对应的处理器
func (h *LevelHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.RLock()
	handlers, ok := h.handlers[r.Level]
	h.mu.RUnlock()
	if !ok {
		return nil
	}
	for _, handler := range handlers {
		if err := handler.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

// WithAttrs 传递属性到所有子处理器
func (h *LevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.mu.RLock()
	defer h.mu.RUnlock()
	newHandlers := make(map[slog.Level][]slog.Handler, len(h.handlers))
	for level, handlers := range h.handlers {
		newH := make([]slog.Handler, 0, len(handlers))
		for _, handler := range handlers {
			newH = append(newH, handler.WithAttrs(attrs))
		}
		newHandlers[level] = newH
	}
	return &LevelHandler{handlers: newHandlers}
}

// WithGroup 传递组到所有子处理器
func (h *LevelHandler) WithGroup(name string) slog.Handler {
	h.mu.RLock()
	defer h.mu.RUnlock()
	newHandlers := make(map[slog.Level][]slog.Handler, len(h.handlers))
	for level, handlers := range h.handlers {
		newH := make([]slog.Handler, 0, len(handlers))
		for _, handler := range handlers {
			newH = append(newH, handler.WithGroup(name))
		}
		newHandlers[level] = newH
	}
	return &LevelHandler{handlers: newHandlers}
}
