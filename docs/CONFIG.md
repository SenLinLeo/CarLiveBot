# 配置项说明

## 配置来源与优先级

1. **configs/config.yaml**：默认值  
2. **环境变量**：覆盖 YAML，键名为大写+下划线（如 `ARK_API_KEY`、`HTTP_PORT`）  
3. **密钥**：仅通过环境变量设置，不写入 YAML

## 全局配置（config.yaml）

| 路径 | 说明 | 示例 |
|------|------|------|
| server.port | HTTP 服务端口 | "8080" |
| server.config_path | 配置根目录（含 stores、compliance） | "configs" |
| ark.chat_url | 方舟 Chat API 地址 | 见文档 |
| ark.model_id | 模型 ID | doubao-seed-2-0-lite-260215 |
| ark.temperature | 生成温度 0~2 | 0.2 |
| ark.max_tokens | 最大生成长度 | 2048 |
| ark.stream | 是否流式 | true |
| ark.thinking_enabled | 是否开启深度思考 | true |
| tts.wss_url | TTS Websocket 地址 | 见文档 |
| tts.voice_type | 音色 ID | zh_female_cancan_mars_bigtts |
| tts.speed_ratio | 语速 0.1~2.0 | 1.0 |
| tts.loudness_ratio | 音量 0.5~2.0 | 1.0 |
| tts.encoding | 音频编码 | mp3 |
| tts.model | TTS 模型 | seed-tts-1.1 |
| compliance.label_text | AI 虚拟人标注文案 | 本直播由AI虚拟人提供服务 |
| compliance.forbidden_words_path | 禁用词文件路径 | configs/compliance/forbidden_words.txt |
| live.replay_interval_minutes | 定时播报间隔（分钟） | 30 |
| live.health_check_interval_seconds | 健康检查间隔 | 60 |
| live.auto_reconnect_interval_seconds | 断流重连间隔 | 5 |

## 环境变量（密钥与覆盖）

| 变量名 | 说明 |
|--------|------|
| ARK_API_KEY | 方舟大模型 API Key（必填，调用 Chat 时使用） |
| TTS_APP_ID | 火山引擎 TTS 应用 ID |
| TTS_ACCESS_TOKEN | 火山引擎 TTS access_token |
| TTS_WSS_URL | 覆盖 tts.wss_url |
| HTTP_PORT | 覆盖 server.port |
| CONFIG_PATH | 覆盖 server.config_path |

## 门店配置（configs/stores/{store_id}.yaml）

| 字段 | 说明 |
|------|------|
| store_id | 门店唯一 ID |
| name | 门店名称 |
| promo_cars | 主推车型（用于话术变量） |
| promo_offers | 优惠说明 |
| lead_link | 留资链接 |
| system_prompt | 系统话术（AI 角色、技能、合规要求） |
| schedule.open_time | 每日开播时间 |
| schedule.close_time | 每日下播时间 |
| schedule.replay_interval_minutes | 定时播报间隔 |
| tts_voice_type | 可选，覆盖全局音色 |

新增门店：复制 `example_store.yaml` 为 `{store_id}.yaml` 并修改上述字段即可。
