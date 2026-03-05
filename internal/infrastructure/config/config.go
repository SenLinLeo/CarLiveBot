package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// App 聚合全局配置
type App struct {
	Server    ServerConfig
	Ark       ArkConfig
	TTS       TTSConfig
	Compliance ComplianceConfig
	Live      LiveConfig
}

type ServerConfig struct {
	Port       string `mapstructure:"port"`
	ConfigPath string `mapstructure:"config_path"`
}

type ArkConfig struct {
	APIKey         string  `mapstructure:"api_key"`
	ChatURL        string  `mapstructure:"chat_url"`
	ModelID        string  `mapstructure:"model_id"`
	Temperature    float64 `mapstructure:"temperature"`
	MaxTokens      int     `mapstructure:"max_tokens"`
	Stream         bool    `mapstructure:"stream"`
	ThinkingEnabled bool   `mapstructure:"thinking_enabled"`
}

type TTSConfig struct {
	AppID          string  `mapstructure:"app_id"`
	AccessToken    string  `mapstructure:"access_token"`
	WSSURL         string  `mapstructure:"wss_url"`
	VoiceType      string  `mapstructure:"voice_type"`
	SpeedRatio     float64 `mapstructure:"speed_ratio"`
	LoudnessRatio  float64 `mapstructure:"loudness_ratio"`
	Encoding       string  `mapstructure:"encoding"`
	Model          string  `mapstructure:"model"`
}

type ComplianceConfig struct {
	LabelText          string `mapstructure:"label_text"`
	ForbiddenWordsPath string `mapstructure:"forbidden_words_path"`
}

type LiveConfig struct {
	ReplayIntervalMinutes       int `mapstructure:"replay_interval_minutes"`
	HealthCheckIntervalSeconds  int `mapstructure:"health_check_interval_seconds"`
	AutoReconnectIntervalSeconds int `mapstructure:"auto_reconnect_interval_seconds"`
}

// StoreConfig 单门店配置（来自 configs/stores/{id}.yaml）
type StoreConfig struct {
	StoreID       string         `mapstructure:"store_id"`
	Name          string         `mapstructure:"name"`
	PromoCars     string         `mapstructure:"promo_cars"`
	PromoOffers   string         `mapstructure:"promo_offers"`
	LeadLink      string         `mapstructure:"lead_link"`
	SystemPrompt  string         `mapstructure:"system_prompt"`
	Schedule      StoreSchedule `mapstructure:"schedule"`
	TTSVoiceType  string         `mapstructure:"tts_voice_type"` // 可选覆盖
}

type StoreSchedule struct {
	OpenTime             string `mapstructure:"open_time"`
	CloseTime            string `mapstructure:"close_time"`
	ReplayIntervalMinutes int    `mapstructure:"replay_interval_minutes"`
}

// Load 从 configs/config.yaml + 环境变量加载配置，密钥仅从环境变量读
func Load() (*App, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("configs")
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		v.AddConfigPath(p)
	}
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	// 环境变量覆盖（大写+下划线）
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfg := &App{}
	bindServer(v, cfg)
	bindArk(v, cfg)
	bindTTS(v, cfg)
	bindCompliance(v, cfg)
	bindLive(v, cfg)

	// 密钥仅环境变量
	if k := os.Getenv("ARK_API_KEY"); k != "" {
		cfg.Ark.APIKey = k
	}
	if k := os.Getenv("TTS_APP_ID"); k != "" {
		cfg.TTS.AppID = k
	}
	if k := os.Getenv("TTS_ACCESS_TOKEN"); k != "" {
		cfg.TTS.AccessToken = k
	}
	if p := os.Getenv("HTTP_PORT"); p != "" {
		cfg.Server.Port = p
	}

	return cfg, nil
}

func bindServer(v *viper.Viper, c *App) {
	c.Server.Port = getString(v, "server.port", "8080")
	c.Server.ConfigPath = getString(v, "server.config_path", "configs")
}

func bindArk(v *viper.Viper, c *App) {
	c.Ark.ChatURL = getString(v, "ark.chat_url", "https://ark.cn-beijing.volces.com/api/v3/chat/completions")
	c.Ark.ModelID = getString(v, "ark.model_id", "doubao-seed-2-0-lite-260215")
	c.Ark.Temperature = getFloat64(v, "ark.temperature", 0.2)
	c.Ark.MaxTokens = getInt(v, "ark.max_tokens", 2048)
	c.Ark.Stream = getBool(v, "ark.stream", true)
	c.Ark.ThinkingEnabled = getBool(v, "ark.thinking_enabled", true)
}

func bindTTS(v *viper.Viper, c *App) {
	c.TTS.WSSURL = getString(v, "tts.wss_url", "wss://openspeech.bytedance.com/api/v3/tts/unidirectional/stream")
	c.TTS.VoiceType = getString(v, "tts.voice_type", "zh_female_cancan_mars_bigtts")
	c.TTS.SpeedRatio = getFloat64(v, "tts.speed_ratio", 1.0)
	c.TTS.LoudnessRatio = getFloat64(v, "tts.loudness_ratio", 1.0)
	c.TTS.Encoding = getString(v, "tts.encoding", "mp3")
	c.TTS.Model = getString(v, "tts.model", "seed-tts-1.1")
}

func bindCompliance(v *viper.Viper, c *App) {
	c.Compliance.LabelText = getString(v, "compliance.label_text", "本直播由AI虚拟人提供服务")
	c.Compliance.ForbiddenWordsPath = getString(v, "compliance.forbidden_words_path", "configs/compliance/forbidden_words.txt")
}

func bindLive(v *viper.Viper, c *App) {
	c.Live.ReplayIntervalMinutes = getInt(v, "live.replay_interval_minutes", 30)
	c.Live.HealthCheckIntervalSeconds = getInt(v, "live.health_check_interval_seconds", 60)
	c.Live.AutoReconnectIntervalSeconds = getInt(v, "live.auto_reconnect_interval_seconds", 5)
}

func getString(v *viper.Viper, key, def string) string {
	v.SetDefault(key, def)
	return v.GetString(key)
}
func getInt(v *viper.Viper, key string, def int) int {
	v.SetDefault(key, def)
	return v.GetInt(key)
}
func getFloat64(v *viper.Viper, key string, def float64) float64 {
	v.SetDefault(key, def)
	return v.GetFloat64(key)
}
func getBool(v *viper.Viper, key string, def bool) bool {
	v.SetDefault(key, def)
	return v.GetBool(key)
}

// LoadStoreConfig 按 store_id 加载门店配置（从 configs/stores/{store_id}.yaml）
func LoadStoreConfig(storeID, configPath string) (*StoreConfig, error) {
	if configPath == "" {
		configPath = "configs"
	}
	path := filepath.Join(configPath, "stores", storeID+".yaml")
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read store config %s: %w", path, err)
	}
	var s StoreConfig
	if err := v.Unmarshal(&s); err != nil {
		return nil, fmt.Errorf("unmarshal store config: %w", err)
	}
	return &s, nil
}
