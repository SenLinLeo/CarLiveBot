package services

import (
	"context"
	"fmt"
	"strings"

	"carlivebot/internal/domain/entities"
	"carlivebot/internal/infrastructure/ark"
	"carlivebot/internal/infrastructure/config"
	"carlivebot/internal/infrastructure/tts"
	"carlivebot/pkg/compliance"
)

// ScriptService 话术服务实现
type ScriptService struct {
	cfg        *config.App
	storeRepo  StoreConfigRepo
	compliance *compliance.Filter
	arkClient  *ark.Client
	ttsClient  *tts.Client
}

// StoreConfigRepo 门店配置仓库接口（避免依赖 infrastructure 具体类型）
type StoreConfigRepo interface {
	GetByID(storeID string) (*entities.StoreConfig, error)
}

// NewScriptService 创建话术服务
func NewScriptService(
	cfg *config.App,
	storeRepo StoreConfigRepo,
	compl *compliance.Filter,
	arkClient *ark.Client,
	ttsClient *tts.Client,
) *ScriptService {
	return &ScriptService{
		cfg:        cfg,
		storeRepo:  storeRepo,
		compliance: compl,
		arkClient:  arkClient,
		ttsClient:  ttsClient,
	}
}

// GenerateOnly 仅流式生成文本
func (s *ScriptService) GenerateOnly(ctx context.Context, storeID, userInput string, onChunk func(text string) error) error {
	if s.compliance != nil && !s.compliance.Allowed(userInput) {
		return fmt.Errorf("compliance: input contains forbidden word")
	}
	store, err := s.storeRepo.GetByID(storeID)
	if err != nil {
		return fmt.Errorf("store config: %w", err)
	}
	req := &ark.ChatRequest{
		Messages:    buildMessages(store.SystemPrompt, userInput),
		Stream:      true,
		Temperature: s.cfg.Ark.Temperature,
		MaxTokens:   s.cfg.Ark.MaxTokens,
	}
	if s.cfg.Ark.ThinkingEnabled {
		req.Thinking = &ark.Thinking{Type: "enabled"}
	}
	return s.arkClient.ChatStream(ctx, req, onChunk)
}

// GenerateAndSynthesize 流式生成并合成语音
func (s *ScriptService) GenerateAndSynthesize(ctx context.Context, storeID, userInput string, onAudio func([]byte) error) error {
	if s.compliance != nil && !s.compliance.Allowed(userInput) {
		return fmt.Errorf("compliance: input contains forbidden word")
	}
	store, err := s.storeRepo.GetByID(storeID)
	if err != nil {
		return fmt.Errorf("store config: %w", err)
	}

	var buf strings.Builder
	req := buildChatRequest(s, store.SystemPrompt, userInput)
	err = s.arkClient.ChatStream(ctx, req, func(content string) error {
		buf.WriteString(content)
		return nil
	})
	if err != nil {
		return err
	}
	text := strings.TrimSpace(buf.String())
	if text == "" {
		return nil
	}
	// 按句或 300 字分段送 TTS，此处简化为整段一次
	return s.ttsClient.SynthesizeStream(ctx, text, onAudio)
}

// RunScript 供调度器调用：promptType 为 "open"（开播口播）或 "replay"（定时播报）
func (s *ScriptService) RunScript(ctx context.Context, storeID, promptType string, onAudio func([]byte) error) error {
	var userInput string
	switch promptType {
	case "open":
		userInput = "请做开播口播：用一句话介绍本店和今日优惠，并提醒观众本直播由AI虚拟人提供服务。"
	case "replay":
		userInput = "请做定时播报：用两句话重复今日主推车型和优惠，并引导用户点击链接留资。"
	default:
		userInput = "请简短介绍今日优惠并引导留资。"
	}
	return s.GenerateAndSynthesize(ctx, storeID, userInput, onAudio)
}

func buildChatRequest(s *ScriptService, systemPrompt, userInput string) *ark.ChatRequest {
	msgs := buildMessages(systemPrompt, userInput)
	req := &ark.ChatRequest{
		Model:       s.arkClient.Model(),
		Messages:    msgs,
		Stream:      true,
		Temperature: s.cfg.Ark.Temperature,
		MaxTokens:   s.cfg.Ark.MaxTokens,
	}
	if s.cfg.Ark.ThinkingEnabled {
		req.Thinking = &ark.Thinking{Type: "enabled"}
	}
	return req
}

func buildMessages(systemPrompt, userInput string) []ark.ChatMessage {
	msgs := []ark.ChatMessage{{Role: "system", Content: systemPrompt}}
	if userInput != "" {
		msgs = append(msgs, ark.ChatMessage{Role: "user", Content: userInput})
	}
	return msgs
}
