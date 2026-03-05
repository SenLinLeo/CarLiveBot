# 车播智盒（CarLiveBot）部署说明

## 环境要求

- Go 1.21+
- 云服务器/物理机建议：4 核 CPU、8G 内存、50G SSD，成都地域优先（低延迟）
- 上行带宽 ≥ 5Mbps，固定公网 IP 更佳

## 部署步骤

### 1. 克隆与依赖

```bash
cd CarLiveBot
go mod tidy
```

### 2. 配置环境变量

复制示例并填写密钥（**不要提交 .env 到仓库**）：

```bash
cp .env.example .env
# 编辑 .env，至少填写：
# ARK_API_KEY=你的方舟API密钥
# TTS_APP_ID=你的TTS应用ID
# TTS_ACCESS_TOKEN=你的TTS access_token
```

加载到当前 shell 再启动（或使用 systemd 的 EnvironmentFile）：

```bash
export $(grep -v '^#' .env | xargs)
go run ./cmd/carlivebot
```

或直接导出后运行：

```bash
export ARK_API_KEY=xxx
export TTS_APP_ID=xxx
export TTS_ACCESS_TOKEN=xxx
go run ./cmd/carlivebot
```

### 3. 门店配置

在 `configs/stores/` 下为每个门店新建 `{store_id}.yaml`，参考 `configs/stores/example_store.yaml`。  
修改 `configs/config.yaml` 可调整全局默认（模型、TTS 音色、合规文案等）。

### 4. 验证

- 健康检查：`curl http://localhost:8080/health`
- 门店列表：`curl http://localhost:8080/api/v1/stores`
- 话术生成（需有效 store_id）：  
  `curl -X POST http://localhost:8080/api/v1/script/generate -H "Content-Type: application/json" -d '{"store_id":"example_store","user_input":"凯美瑞多少钱？"}'`

### 5. 生产运行

建议使用 systemd 或 supervisor，并通过环境变量注入密钥。示例 systemd unit：

```ini
[Unit]
Description=CarLiveBot
After=network.target

[Service]
Type=simple
WorkingDirectory=/path/to/CarLiveBot
EnvironmentFile=/path/to/CarLiveBot/.env
ExecStart=/path/to/CarLiveBot/carlivebot
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

## 与直播伴侣对接

见 [LIVE_COMPANION.md](LIVE_COMPANION.md)。
