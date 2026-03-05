# 车播智盒（CarLiveBot）

抖音本地生活 + 汽车垂类 + AI 数字人直播代运营技术端：文字处理 → 语音合成 → 直播调度，业务全可配置，支持多门店。

## 快速开始

1. 复制环境变量：`cp .env.example .env`，填入 `ARK_API_KEY`、TTS 等密钥。
2. 加载依赖：`go mod tidy`。
3. 运行：`go run ./cmd/carlivebot`。
4. 管理 API：`http://localhost:8080/health`、`http://localhost:8080/api/v1/stores` 等。

详见 [docs/DEPLOY.md](docs/DEPLOY.md)、[docs/CONFIG.md](docs/CONFIG.md)、[docs/LIVE_COMPANION.md](docs/LIVE_COMPANION.md)。

## 目录结构

- `cmd/carlivebot`：入口。
- `configs`：全局与门店配置、合规词表。
- `internal/domain`：实体与仓库接口。
- `internal/application`：服务、命令、查询。
- `internal/infrastructure`：配置、Ark Chat、TTS、DB。
- `internal/interface/api/rest`：管理 API。
- `pkg/compliance`：合规校验。

## 技术栈

Go 1.21+，Echo，Viper，火山引擎方舟 Chat API，火山引擎 TTS V3 Websocket。
