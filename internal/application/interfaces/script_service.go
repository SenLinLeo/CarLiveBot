package interfaces

import "context"

// ScriptService 话术服务：合规过滤 + Chat 流式 + TTS 流式
type ScriptService interface {
	// GenerateAndSynthesize 根据门店 system prompt 与用户输入，流式生成并合成语音，音频通过 onAudio 回调
	GenerateAndSynthesize(ctx context.Context, storeID, userInput string, onAudio func([]byte) error) error
	// GenerateOnly 仅生成文本（流式），用于测试或仅要文案
	GenerateOnly(ctx context.Context, storeID, userInput string, onChunk func(text string) error) error
}
