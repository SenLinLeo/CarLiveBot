package tts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const defaultDialTimeout = 15 * time.Second

// Client 火山引擎 TTS V3 Websocket 客户端
type Client struct {
	wssURL      string
	accessToken string
	appID       string
	voiceType   string
	speedRatio  float64
	loudness    float64
	encoding    string
	model       string
	dialer      websocket.Dialer
}

// NewClient 创建 TTS 客户端
func NewClient(wssURL, accessToken, appID, voiceType string, speedRatio, loudness float64, encoding, model string) *Client {
	if wssURL == "" {
		wssURL = "wss://openspeech.bytedance.com/api/v3/tts/unidirectional/stream"
	}
	if encoding == "" {
		encoding = "mp3"
	}
	return &Client{
		wssURL:      wssURL,
		accessToken: accessToken,
		appID:       appID,
		voiceType:   voiceType,
		speedRatio:  speedRatio,
		loudness:    loudness,
		encoding:    encoding,
		model:       model,
		dialer:      websocket.Dialer{HandshakeTimeout: defaultDialTimeout},
	}
}

// SynthesizeStream 将文本合成为语音流，通过 fn 回调每段二进制音频；reqid 每次唯一
func (c *Client) SynthesizeStream(ctx context.Context, text string, fn func(data []byte) error) error {
	reqid := uuid.New().String()
	req := TTSRequest{
		App: TTSApp{
			AppID:   c.appID,
			Token:   "carlivebot_token",
			Cluster: "volcano_tts",
		},
		User: TTSUser{UID: "carlivebot_001"},
		Audio: TTSAudio{
			VoiceType:     c.voiceType,
			SpeedRatio:    c.speedRatio,
			LoudnessRatio: c.loudness,
			Encoding:      c.encoding,
		},
		Request: TTSReqBody{
			ReqID:     reqid,
			Text:      text,
			Operation: "submit",
			Model:     c.model,
		},
	}
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("tts marshal: %w", err)
	}

	// Authorization: Bearer; {token} 注意分号+空格
	header := http.Header{}
	header.Set("Authorization", "Bearer; "+c.accessToken)

	u, err := url.Parse(c.wssURL)
	if err != nil {
		return fmt.Errorf("tts url: %w", err)
	}

	conn, _, err := c.dialer.DialContext(ctx, u.String(), header)
	if err != nil {
		return fmt.Errorf("tts dial: %w", err)
	}
	defer conn.Close()

	if err := conn.WriteMessage(websocket.TextMessage, body); err != nil {
		return fmt.Errorf("tts write: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		_, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return nil
			}
			return fmt.Errorf("tts read: %w", err)
		}
		if len(data) == 0 {
			continue
		}
		if err := fn(data); err != nil {
			return err
		}
	}
}
