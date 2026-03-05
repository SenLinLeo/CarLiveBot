package config

import (
	"os"
	"path/filepath"

	"carlivebot/internal/domain/entities"
)

// StoreConfigRepo 从 configs/stores/*.yaml 加载门店配置
type StoreConfigRepo struct {
	configPath string
	cache      map[string]*entities.StoreConfig
}

// NewStoreConfigRepo 创建门店配置仓库
func NewStoreConfigRepo(configPath string) *StoreConfigRepo {
	if configPath == "" {
		configPath = "configs"
	}
	return &StoreConfigRepo{configPath: configPath, cache: make(map[string]*entities.StoreConfig)}
}

// GetByID 按 store_id 加载并转为 domain 实体
func (r *StoreConfigRepo) GetByID(storeID string) (*entities.StoreConfig, error) {
	if c, ok := r.cache[storeID]; ok {
		return c, nil
	}
	cfg, err := LoadStoreConfig(storeID, r.configPath)
	if err != nil {
		return nil, err
	}
	e := toEntity(cfg)
	r.cache[storeID] = e
	return e, nil
}

// ListStoreIDs 列出 configs/stores 下所有 .yaml 对应的 store_id
func (r *StoreConfigRepo) ListStoreIDs() ([]string, error) {
	dir := filepath.Join(r.configPath, "stores")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if len(name) > 5 && name[len(name)-5:] == ".yaml" {
			ids = append(ids, name[:len(name)-5])
		}
	}
	return ids, nil
}

func toEntity(c *StoreConfig) *entities.StoreConfig {
	return &entities.StoreConfig{
		StoreID:       c.StoreID,
		Name:          c.Name,
		PromoCars:     c.PromoCars,
		PromoOffers:   c.PromoOffers,
		LeadLink:      c.LeadLink,
		SystemPrompt:  c.SystemPrompt,
		OpenTime:      c.Schedule.OpenTime,
		CloseTime:     c.Schedule.CloseTime,
		ReplayMinutes: c.Schedule.ReplayIntervalMinutes,
		TTSVoiceType:  c.TTSVoiceType,
	}
}
