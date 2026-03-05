package request

// ScriptGenerateRequest 话术生成请求
type ScriptGenerateRequest struct {
	StoreID   string `json:"store_id"`
	UserInput string `json:"user_input"`
}
