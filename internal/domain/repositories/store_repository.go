package repositories

import "carlivebot/internal/domain/entities"

// StoreConfigRepository 门店配置读取（当前为文件，可扩展为 DB）
type StoreConfigRepository interface {
	GetByID(storeID string) (*entities.StoreConfig, error)
	ListStoreIDs() ([]string, error)
}
