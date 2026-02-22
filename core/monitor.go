package core

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"log/slog"
)

// Metrics 性能指标
type Metrics struct {
	StartTime     time.Time     `json:"start_time"`
	Uptime        time.Duration `json:"uptime"`
	Goroutines    int           `json:"goroutines"`
	MemAllocated  uint64        `json:"mem_allocated_bytes"`
	MemSys        uint64        `json:"mem_sys_bytes"`
	MemGC         uint64        `json:"mem_gc_bytes"`
	NumGC         uint32        `json:"num_gc"`
	LastGC        time.Time     `json:"last_gc"`
	DBConnections map[string]int `json:"db_connections"`
	CacheHits     uint64        `json:"cache_hits"`
	CacheMisses   uint64        `json:"cache_misses"`
}

// Monitor 监控器
type Monitor struct {
	app           *Application
	metrics       *Metrics
	mu            sync.RWMutex
	cacheHitCount uint64
	cacheMissCount uint64
}

// NewMonitor 创建监控器
func NewMonitor(app *Application) *Monitor {
	return &Monitor{
		app:     app,
		metrics: &Metrics{StartTime: time.Now()},
	}
}

// CollectMetrics 收集性能指标
func (m *Monitor) CollectMetrics() *Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	m.metrics.Uptime = time.Since(m.metrics.StartTime)
	m.metrics.Goroutines = runtime.NumGoroutine()
	m.metrics.MemAllocated = ms.Alloc
	m.metrics.MemSys = ms.Sys
	m.metrics.MemGC = ms.NextGC
	m.metrics.NumGC = ms.NumGC
	if ms.LastGC > 0 {
		m.metrics.LastGC = time.Unix(0, int64(ms.LastGC))
	}

	// 收集数据库连接信息
	m.metrics.DBConnections = make(map[string]int)
	m.app.mu.RLock()
	for key, db := range m.app.databases {
		if sqlDB, err := db.DB(); err == nil {
			stats := sqlDB.Stats()
			m.metrics.DBConnections[key] = stats.OpenConnections
		}
	}
	m.app.mu.RUnlock()

	// 收集缓存统计
	m.metrics.CacheHits = m.cacheHitCount
	m.metrics.CacheMisses = m.cacheMissCount

	return m.metrics
}

// IncrementCacheHit 增加缓存命中计数
func (m *Monitor) IncrementCacheHit() {
	m.mu.Lock()
	m.cacheHitCount++
	m.mu.Unlock()
}

// IncrementCacheMiss 增加缓存未命中计数
func (m *Monitor) IncrementCacheMiss() {
	m.mu.Lock()
	m.cacheMissCount++
	m.mu.Unlock()
}

// HealthCheck 健康检查
type HealthCheck struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"` // "healthy", "unhealthy", "degraded"
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"`
}

// HealthChecker 健康检查接口
type HealthChecker interface {
	Check(ctx context.Context) *HealthCheck
	Name() string
}

// HealthService 健康服务
type HealthService struct {
	checkers []HealthChecker
	logger   *slog.Logger
}

// NewHealthService 创建健康服务
func NewHealthService(logger *slog.Logger) *HealthService {
	return &HealthService{
		checkers: make([]HealthChecker, 0),
		logger:   logger,
	}
}

// RegisterChecker 注册健康检查器
func (hs *HealthService) RegisterChecker(checker HealthChecker) {
	hs.checkers = append(hs.checkers, checker)
	hs.logger.Info("Registered health checker", "name", checker.Name())
}

// CheckAll 执行所有健康检查
func (hs *HealthService) CheckAll(ctx context.Context) []*HealthCheck {
	results := make([]*HealthCheck, len(hs.checkers))
	
	for i, checker := range hs.checkers {
		result := checker.Check(ctx)
		results[i] = result
		
		if result.Status != "healthy" {
			hs.logger.Warn("Health check failed", 
				"name", result.Name, 
				"status", result.Status, 
				"message", result.Message,
				"error", result.Error)
		}
	}
	
	return results
}

// GetOverallStatus 获取整体健康状态
func (hs *HealthService) GetOverallStatus(results []*HealthCheck) string {
	status := "healthy"
	
	for _, result := range results {
		if result.Status == "unhealthy" {
			return "unhealthy"
		}
		if result.Status == "degraded" {
			status = "degraded"
		}
	}
	
	return status
}

// DatabaseHealthChecker 数据库健康检查器
type DatabaseHealthChecker struct {
	app *Application
}

func NewDatabaseHealthChecker(app *Application) *DatabaseHealthChecker {
	return &DatabaseHealthChecker{app: app}
}

func (dhc *DatabaseHealthChecker) Name() string {
	return "database"
}

func (dhc *DatabaseHealthChecker) Check(ctx context.Context) *HealthCheck {
	check := &HealthCheck{
		Name:      dhc.Name(),
		Timestamp: time.Now(),
	}

	dhc.app.mu.RLock()
	defer dhc.app.mu.RUnlock()

	if len(dhc.app.databases) == 0 {
		check.Status = "degraded"
		check.Message = "no databases configured"
		return check
	}

	healthyCount := 0
	totalCount := len(dhc.app.databases)

	for key, db := range dhc.app.databases {
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.PingContext(ctx); err != nil {
				check.Status = "unhealthy"
				check.Error = fmt.Sprintf("database %s ping failed: %v", key, err)
				return check
			}
			healthyCount++
		} else {
			check.Status = "unhealthy"
			check.Error = fmt.Sprintf("failed to get sql.DB for %s: %v", key, err)
			return check
		}
	}

	if healthyCount == totalCount {
		check.Status = "healthy"
		check.Message = fmt.Sprintf("all %d databases healthy", totalCount)
	} else {
		check.Status = "degraded"
		check.Message = fmt.Sprintf("only %d/%d databases healthy", healthyCount, totalCount)
	}

	return check
}

// CacheHealthChecker 缓存健康检查器
type CacheHealthChecker struct {
	app *Application
}

func NewCacheHealthChecker(app *Application) *CacheHealthChecker {
	return &CacheHealthChecker{app: app}
}

func (chc *CacheHealthChecker) Name() string {
	return "cache"
}

func (chc *CacheHealthChecker) Check(ctx context.Context) *HealthCheck {
	check := &HealthCheck{
		Name:      chc.Name(),
		Timestamp: time.Now(),
	}

	cache := chc.app.GetCache()
	if cache == nil {
		check.Status = "unhealthy"
		check.Error = "cache not initialized"
		return check
	}

	// 测试缓存连接
	testKey := "health_check_" + time.Now().Format("20060102150405")
	testValue := "ok"
	
	if err := cache.Set(testKey, testValue, time.Minute); err != nil {
		check.Status = "unhealthy"
		check.Error = fmt.Sprintf("failed to set cache value: %v", err)
		return check
	}

	if val, err := cache.Get(testKey); err != nil {
		check.Status = "unhealthy"
		check.Error = fmt.Sprintf("failed to get cache value: %v", err)
		return check
	} else if val != testValue {
		check.Status = "degraded"
		check.Message = "cache value mismatch"
		return check
	}

	// 清理测试数据
	cache.Del(testKey)

	check.Status = "healthy"
	check.Message = fmt.Sprintf("cache type: %s", cache.Type())
	return check
}