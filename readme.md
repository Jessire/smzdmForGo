# smzdmForGo

`smzdmForGo` 是一个什么值得买好价监控和 Telegram 推送服务. 程序会按关键词规则搜索什么值得买文章, 过滤不需要的商品, 对已推送内容做去重, 并按固定间隔推送命中的好价信息.

项目提供一个 Web 面板, 用来维护商品规则、搜索预览和 Telegram 通知配置.

## 界面预览

![Web 面板](docs/web-panel.png)

## 功能

- 按关键词搜索什么值得买好价文章
- 支持多组独立关键词规则
- 支持过滤词、最低评论数、最低值率、价格区间
- 支持搜索预览
- 支持 Telegram Bot 推送
- 支持定时监控和去重推送
- 支持 SQLite 本地运行
- 支持 PostgreSQL 生产部署
- 支持 Docker 和 Render 部署

## 快速开始

项目需要 Go 版本与 `go.mod` 保持一致.

```bash
go mod download
go run .
```

启动后访问:

```text
http://localhost:9090
```

健康检查:

```text
http://localhost:9090/health
```

默认端口是 `9090`. 如果设置了 `PORT` 环境变量, 服务会监听对应端口.

## 配置

程序启动时读取 `config/config.yml`. Web 面板保存后的配置会写入数据库, 后续启动会优先使用数据库中的配置.

示例:

```yaml
keyWords:
  - 显示器
  - 面包

lowCommentNum: 1
lowWorthyNum: 6
minPrice: 0
maxPrice: 0
satisfyNum: 5
filterWords:
  - 过期
  - 售罄

keywordRules:
  - enabled: true
    words:
      - 显示器
    filterWords:
      - 二手
      - 支架
    lowCommentNum: 5
    lowWorthyNum: 20
    minPrice: 300
    maxPrice: 2000

tickTime: 10800

telegram:
  enabled: false
  botToken: ""
  chatId: ""
  parseMode: "HTML"
  disableWebPagePreview: false
```

主要字段:

- `keyWords`: 默认关键词列表.
- `keywordRules`: 独立关键词规则. 配置后按规则逐组搜索和过滤.
- `filterWords`: 过滤词.
- `lowCommentNum`: 最低评论数.
- `lowWorthyNum`: 最低值率.
- `minPrice`: 最低价格, `0` 表示不限制.
- `maxPrice`: 最高价格, `0` 表示不限制.
- `satisfyNum`: 单次推送数量上限.
- `tickTime`: 定时任务间隔, 单位秒.
- `telegram`: Telegram Bot 推送配置.

## 数据存储

本地默认使用 SQLite:

```text
data/users.db
```

运行时还会生成:

```text
pushed.json
```

`pushed.json` 用于记录已推送文章, 避免重复推送.

设置下面任一环境变量后, 程序会使用 PostgreSQL:

- `DATABASE_URL`
- `SQL_DSN`
- `AXONHUB_DB_DSN`

生产环境建议同时设置:

```bash
REQUIRE_DATABASE_URL=true
```

可选表名环境变量:

- `SMZDM_USERS_TABLE`
- `SMZDM_SETTINGS_TABLE`

## Docker

构建镜像:

```bash
docker build -t smzdm-for-go .
```

运行:

```bash
docker run -d --name smzdm-for-go -p 9090:9090 smzdm-for-go
```

使用 PostgreSQL:

```bash
docker run -d \
  --name smzdm-for-go \
  -p 9090:9090 \
  -e DATABASE_URL="postgres://user:password@host:5432/dbname?sslmode=require" \
  -e REQUIRE_DATABASE_URL=true \
  smzdm-for-go
```

## Render

仓库包含 `render.yaml`, 可作为 Render Web Service 部署.

当前部署配置:

- Runtime: Docker
- Health Check Path: `/health`
- Branch: `master`
- Timezone: `Asia/Shanghai`
- Required env: `DATABASE_URL`
- Required env: `REQUIRE_DATABASE_URL=true`

部署后在 Web 面板配置关键词规则和 Telegram Bot 信息.

## 目录说明

```text
.
├── main.go             # 服务入口和定时任务
├── route*.go           # Web 页面和接口
├── config/             # 默认配置
├── db/                 # SQLite/PostgreSQL 访问
├── docs/               # README 图片资源
├── smzdm/              # 什么值得买搜索和过滤
├── push/               # Telegram 推送
├── template/           # Web 面板页面
├── Dockerfile
└── render.yaml
```

## 许可证

见 `LICENSE`.
