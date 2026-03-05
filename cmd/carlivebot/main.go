package main

import (
	"log"
	"os"

	"carlivebot/internal/application/services"
	"carlivebot/internal/infrastructure/ark"
	"carlivebot/internal/infrastructure/config"
	"carlivebot/internal/infrastructure/tts"
	"carlivebot/internal/interface/api/rest"
	"carlivebot/pkg/compliance"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// 合规词过滤（文件不存在时跳过，不阻塞启动）
	var compl *compliance.Filter
	if path := cfg.Compliance.ForbiddenWordsPath; path != "" {
		if f, err := compliance.NewFilter(path); err == nil {
			compl = f
		}
	}

	storeRepo := config.NewStoreConfigRepo(cfg.Server.ConfigPath)
	arkClient := ark.NewClient(cfg.Ark.APIKey, cfg.Ark.ChatURL, cfg.Ark.ModelID)
	ttsClient := tts.NewClient(
		cfg.TTS.WSSURL,
		cfg.TTS.AccessToken,
		cfg.TTS.AppID,
		cfg.TTS.VoiceType,
		cfg.TTS.SpeedRatio,
		cfg.TTS.LoudnessRatio,
		cfg.TTS.Encoding,
		cfg.TTS.Model,
	)
	scriptSvc := services.NewScriptService(cfg, storeRepo, compl, arkClient, ttsClient)

	e := echo.New()
	deps := &rest.Deps{Config: cfg, Script: scriptSvc}
	rest.RegisterRoutes(e, deps)

	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("CarLiveBot starting on %s", addr)
	if err := e.Start(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
	os.Exit(0)
}
