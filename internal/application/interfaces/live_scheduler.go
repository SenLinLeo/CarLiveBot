package interfaces

// LiveScheduler 直播调度：按门店配置的开播/定时播报触发话术
type LiveScheduler interface {
	// Start 启动调度（阻塞或后台）
	Start()
	// Stop 停止调度
	Stop()
}
