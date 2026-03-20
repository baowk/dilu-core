package core

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// ErrorHandler 错误处理器接口
type ErrorHandler interface {
	HandleError(c *gin.Context, err error)
	HandlePanic(c *gin.Context, recovered interface{})
}

// DefaultErrorHandler 默认错误处理器
type DefaultErrorHandler struct {
	logger zerolog.Logger
}

// NewDefaultErrorHandler 创建默认错误处理器
func NewDefaultErrorHandler(logger zerolog.Logger) *DefaultErrorHandler {
	return &DefaultErrorHandler{logger: logger}
}

// HandleError 处理普通错误
func (eh *DefaultErrorHandler) HandleError(c *gin.Context, err error) {
	// 记录错误日志
	eh.logger.Error().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.String()).
		Err(err).
		Str("client_ip", c.ClientIP()).
		Msg("Request error")

	// 返回错误响应（不暴露内部错误详情）
	if bizErr, ok := err.(*BusinessError); ok {
		c.JSON(bizErr.Code, gin.H{
			"code":    bizErr.Code,
			"message": bizErr.Message,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
	}
}

// HandlePanic 处理panic
func (eh *DefaultErrorHandler) HandlePanic(c *gin.Context, recovered interface{}) {
	// 记录panic日志
	stack := string(debug.Stack())
	eh.logger.Error().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.String()).
		Str("recovered", fmt.Sprintf("%v", recovered)).
		Str("stack", stack).
		Str("client_ip", c.ClientIP()).
		Msg("Panic recovered")

	// 返回错误响应
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"message": "Internal Server Error",
		"error":   "An unexpected error occurred",
	})
}

// RecoveryMiddleware panic恢复中间件
func RecoveryMiddleware(errorHandler ErrorHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				errorHandler.HandlePanic(c, recovered)
				c.Abort()
			}
		}()
		c.Next()
	}
}

// ErrorMiddleware 错误处理中间件
func ErrorMiddleware(errorHandler ErrorHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// 检查是否有错误
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				errorHandler.HandleError(c, err.Err)
			}
		}
	}
}

// BusinessError 业务错误类型
type BusinessError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (be *BusinessError) Error() string {
	return fmt.Sprintf("[%d] %s", be.Code, be.Message)
}

// NewBusinessError 创建业务错误
func NewBusinessError(code int, message string, detail string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// ValidationError 验证错误类型
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", ve.Field, ve.Message)
}

// ValidationErrors 验证错误集合
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (ves *ValidationErrors) Error() string {
	return fmt.Sprintf("validation failed with %d errors", len(ves.Errors))
}

// AddError 添加验证错误
func (ves *ValidationErrors) AddError(field, message string, value interface{}) {
	ves.Errors = append(ves.Errors, ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
}

// HasErrors 检查是否有错误
func (ves *ValidationErrors) HasErrors() bool {
	return len(ves.Errors) > 0
}

// ValidationErrorHandler 验证错误处理器
type ValidationErrorHandler struct {
	logger zerolog.Logger
}

func NewValidationErrorHandler(logger zerolog.Logger) *ValidationErrorHandler {
	return &ValidationErrorHandler{logger: logger}
}

func (veh *ValidationErrorHandler) HandleValidationError(c *gin.Context, errs *ValidationErrors) {
	veh.logger.Warn().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.String()).
		Interface("errors", errs.Errors).
		Str("client_ip", c.ClientIP()).
		Msg("Validation errors")

	c.JSON(http.StatusBadRequest, gin.H{
		"code":    http.StatusBadRequest,
		"message": "Validation failed",
		"errors":  errs.Errors,
	})
}

// TimeoutMiddleware 超时中间件
func TimeoutMiddleware(timeoutSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建带超时的上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(timeoutSeconds)*time.Second)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		// 使用带缓冲的 channel，防止 goroutine 在超时后泄漏
		done := make(chan struct{}, 1)
		go func() {
			c.Next()
			done <- struct{}{}
		}()

		select {
		case <-done:
			// 正常完成
		case <-ctx.Done():
			// 超时：先 Abort 阻止后续 handler 写响应，再写超时响应
			c.Abort()
			c.JSON(http.StatusRequestTimeout, gin.H{
				"code":    http.StatusRequestTimeout,
				"message": "Request timeout",
			})
		}
	}
}