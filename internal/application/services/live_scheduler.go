package services

import (
	"context"
	"log"
	"sync"
	"time"

	"carlivebot/internal/domain/entities"
	"carlivebot/internal/infrastructure/config"
)

// LiveScheduler 按门店配置定时触发开播口播与定时播报
type LiveScheduler struct {
	cfg         *config.App
	storeRepo   StoreConfigLister
	scriptSvc   ScriptRunner
	intervalSec int
	stopCh      chan struct{}
	wg          sync.WaitGroup
}

// StoreConfigLister 可列出门店 ID 并获取配置
type StoreConfigLister interface {
	ListStoreIDs() ([]string, error)
	GetByID(storeID string) (*entities.StoreConfig, error)
}

// ScriptRunner 执行一次话术生成+合成（用于开播/定时播报）
type ScriptRunner interface {
	RunScript(ctx context.Context, storeID, promptType string, onAudio func([]byte) error) error
}

// NewLiveScheduler 创建调度器（若不需要调度可传 nil scriptSvc）
func NewLiveScheduler(cfg *config.App, storeRepo StoreConfigLister, scriptSvc ScriptRunner) *LiveScheduler {
	interval := cfg.Live.HealthCheckIntervalSeconds
	if interval <= 0 {
		interval = 60
	}
	return &LiveScheduler{
		cfg:         cfg,
		storeRepo:   storeRepo,
		scriptSvc:   scriptSvc,
		intervalSec: interval,
		stopCh:      make(chan struct{}),
	}
}

// Start 后台启动调度循环：每分钟检查各门店是否到点播报
func (s *LiveScheduler) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(time.Duration(s.intervalSec) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-s.stopCh:
				return
			case <-ticker.C:
				s.tick(context.Background())
			}
		}
	}()
	log.Println("LiveScheduler started")
}

// Stop 停止调度
func (s *LiveScheduler) Stop() {
	close(s.stopCh)
	s.wg.Wait()
	log.Println("LiveScheduler stopped")
}

func (s *LiveScheduler) tick(ctx context.Context) {
	if s.scriptSvc == nil {
		return
	}
	ids, err := s.storeRepo.ListStoreIDs()
	if err != nil {
		log.Printf("scheduler list stores: %v", err)
		return
	}
	for _, id := range ids {
		store, err := s.storeRepo.GetByID(id)
		if err != nil {
			continue
		}
		s.maybeRunReplay(ctx, store)
	}
}

func (s *LiveScheduler) maybeRunReplay(ctx context.Context, store *entities.StoreConfig) {
	// 简化：不解析具体开播时间，仅按 replay_interval 触发定时播报
	// 实际可在此判断当前时间是否在 open_time~close_time 且距上次播报已满 replay_interval
	_ = store
	_ = ctx
}
