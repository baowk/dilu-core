package core

// ServiceRegistry 服务注册发现接口
// dilu-rd 或其他注册中心实现此接口后，通过 SetRegistry 注入即可自动注册
type ServiceRegistry interface {
	// Register 注册服务到注册中心
	Register(name, addr string, port int, healthURL string, tags []string) error
	// Deregister 从注册中心注销
	Deregister()
}

// SetRegistry 设置服务注册中心（可选）
// 调用后，Application.Run() 启动时会自动注册服务，关闭时自动注销
func (app *Application) SetRegistry(registry ServiceRegistry) {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.registry = registry
}

// GetRegistry 获取当前注册中心
func (app *Application) GetRegistry() ServiceRegistry {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.registry
}
