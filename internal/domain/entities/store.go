package entities

import "github.com/google/uuid"

// StoreConfig 门店配置（业务视图）
type StoreConfig struct {
	StoreID       string
	Name          string
	PromoCars     string
	PromoOffers   string
	LeadLink      string
	SystemPrompt  string
	OpenTime      string
	CloseTime     string
	ReplayMinutes int
	TTSVoiceType  string
}

// LiveRoom 直播间（与门店 1:1 或 1:N，当前简化为一店一直播间）
type LiveRoom struct {
	ID        uuid.UUID
	StoreID   string
	Name      string
	Status    string // idle, live, error
	UpdatedAt int64
}

// Lead 留资记录
type Lead struct {
	ID        uuid.UUID
	StoreID   string
	Phone     string
	Extra     string
	CreatedAt int64
}
