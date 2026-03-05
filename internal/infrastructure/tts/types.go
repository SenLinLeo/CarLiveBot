package tts

// TTSRequest V3 单向流式请求体（JSON）
type TTSRequest struct {
	App     TTSApp     `json:"app"`
	User    TTSUser    `json:"user"`
	Audio   TTSAudio   `json:"audio"`
	Request TTSReqBody `json:"request"`
}

type TTSApp struct {
	AppID   string `json:"appid"`
	Token   string `json:"token"`
	Cluster string `json:"cluster"`
}

type TTSUser struct {
	UID string `json:"uid"`
}

type TTSAudio struct {
	VoiceType      string  `json:"voice_type"`
	SpeedRatio     float64 `json:"speed_ratio,omitempty"`
	LoudnessRatio  float64 `json:"loudness_ratio,omitempty"`
	Encoding       string  `json:"encoding,omitempty"`
}

type TTSReqBody struct {
	ReqID      string `json:"reqid"`
	Text       string `json:"text"`
	Operation  string `json:"operation"` // "submit"
	Model      string `json:"model,omitempty"`
	WithTimestamp int  `json:"with_timestamp,omitempty"`
}
