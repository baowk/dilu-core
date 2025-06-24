package logger

import (
	"context"
	"log/slog"
	"sync"
)

// LevelHandler 将不同级别的日志分发到对应的 slog.Handler
type LevelHandler struct {
	handlers map[slog.Level][]slog.Handler
	mu       sync.Mutex
}

func NewLevelHandler() *LevelHandler {
	return &LevelHandler{
		handlers: make(map[slog.Level][]slog.Handler),
	}
}

// AddHandler 添加级别对应的处理器
func (h *LevelHandler) AddHandler(level slog.Level, handler slog.Handler) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[level] = append(h.handlers[level], handler)
}

// Enabled 检查当前级别是否有对应的处理器
func (h *LevelHandler) Enabled(_ context.Context, level slog.Level) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, ok := h.handlers[level]
	return ok
}

// Handle 分发日志记录到对应的处理器
func (h *LevelHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	handlers, ok := h.handlers[r.Level]
	h.mu.Unlock()
	if !ok {
		return nil
	}
	for _, handler := range handlers {
		handler.Handle(ctx, r)

	}
	return nil
}

// WithAttrs 传递属性到所有子处理器
func (h *LevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.mu.Lock()
	defer h.mu.Unlock()
	newHandlers := make(map[slog.Level][]slog.Handler)
	for level, handlers := range h.handlers {
		for _, handler := range handlers {
			newHandlers[level] = append(newHandlers[level], handler.WithAttrs(attrs))
		}
	}
	return &LevelHandler{handlers: newHandlers}
}

// WithGroup 传递组到所有子处理器
func (h *LevelHandler) WithGroup(name string) slog.Handler {
	h.mu.Lock()
	defer h.mu.Unlock()
	newHandlers := make(map[slog.Level][]slog.Handler)
	for level, handlers := range h.handlers {
		for _, handler := range handlers {
			newHandlers[level] = append(newHandlers[level], handler.WithGroup(name))
		}
	}
	return &LevelHandler{handlers: newHandlers}
}
