package rest

import (
	"carlivebot/internal/application/interfaces"
	"carlivebot/internal/infrastructure/config"
)

// Deps 注入 REST 所需依赖（Config 必填，Script 可选）
type Deps struct {
	Config *config.App
	Script interfaces.ScriptService
}
