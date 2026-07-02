## smzdmForGo

`smzdmForGo` 是一个基于 Go 的什么值得买监控与推送服务. 它会按配置搜索好价文章, 根据关键词、过滤词、评论数、值率和价格区间筛选商品, 去重后通过 Telegram Bot 推送. 项目同时提供 Web 面板, 用于维护商品规则、Telegram 配置和什么值得买签到 Cookie.

### 主要功能

- Web 面板维护商品监控规则和 Telegram Bot 配置
- 按关键词搜索什么值得买好价文章
- 支持每组关键词独立配置过滤词、最低评论数、最低值率、最低价和最高价
- 支持定时抓取与推送
- 使用 `pushed.json` 记录已推送文章, 避免重复推送
- 支持 Telegram Bot 消息推送
- 支持保存什么值得买 Cookie 并执行签到
- 本地默认使用 SQLite, 生产环境可使用 PostgreSQL
- 支持 Docker 和 Render Web Service 部署

### 项目结构

```text
.
├── main.go                 # Web 服务入口和定时任务
├── route*.go               # Web 页面与接口路由
├── config/config.yml       # 默认商品规则和 Telegram 配置
├── data/users.db           # 本地 SQLite 数据库
├── db/                     # SQLite/PostgreSQL 数据访问
├── smzdm/                  # 什么值得买文章抓取与过滤
├── push/                   # Telegram 推送
├── check_in/               # 什么值得买签到
├── template/               # Web 面板页面和默认 JSON
├── Dockerfile              # Docker 镜像构建
└── render.yaml             # Render 部署配置
```

### 本地运行

要求 Go 版本与 `go.mod` 保持一致.

```bash
go mod download
go run .
```

服务启动后默认监听 `9090`:

```text
http://localhost:9090
```

可通过健康检查确认服务状态:

```text
http://localhost:9090/health
```

如果环境变量 `PORT` 存在, 服务会监听该端口.

### 配置说明

程序启动时会读取 `config/config.yml`, 并在 Web 面板保存后把生产配置写入数据库. 推荐通过 Web 面板维护商品规则和 Telegram 配置.

当前配置结构示例:

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

字段含义:

- `keyWords`: 全局搜索关键词. 未配置 `keywordRules` 时使用.
- `filterWords`: 全局过滤词.
- `lowCommentNum`: 最低评论数.
- `lowWorthyNum`: 最低值率.
- `minPrice`: 最低价格, `0` 表示不限制.
- `maxPrice`: 最高价格, `0` 表示不限制.
- `satisfyNum`: 每次推送的商品数量上限.
- `keywordRules`: 独立关键词规则. 配置后会按每组规则搜索和过滤商品.
- `tickTime`: 商品监控执行间隔, 单位秒. 默认配置为 3 小时.
- `telegram`: Telegram Bot 推送配置.

### 数据库

本地运行默认使用 SQLite:

```text
data/users.db
```

如果设置了以下任一环境变量, 程序会改用 PostgreSQL:

- `DATABASE_URL`
- `SQL_DSN`
- `AXONHUB_DB_DSN`

生产环境可设置:

```bash
REQUIRE_DATABASE_URL=true
```

这样可以强制服务必须使用 PostgreSQL DSN 启动, 避免误用容器内 SQLite. 表名可通过环境变量覆盖:

- `SMZDM_USERS_TABLE`: 用户和签到信息表名. PostgreSQL 默认 `smzdm_users`, SQLite 默认 `users`.
- `SMZDM_SETTINGS_TABLE`: 应用配置表名. 默认 `smzdm_app_settings`.

### Docker 运行

构建镜像:

```bash
docker build -t smzdm-for-go .
```

本地 SQLite 方式运行:

```bash
docker run -d \
  --name smzdm-for-go \
  -p 9090:9090 \
  -v "$(pwd)/data:/opt/go/data" \
  smzdm-for-go
```

使用 PostgreSQL 运行:

```bash
docker run -d \
  --name smzdm-for-go \
  -p 9090:9090 \
  -e DATABASE_URL="postgres://user:password@host:5432/dbname?sslmode=require" \
  -e REQUIRE_DATABASE_URL=true \
  smzdm-for-go
```

### Render 部署

仓库包含 `render.yaml`, 当前配置会以 Docker Web Service 部署:

- 服务名: `smzdm-for-go`
- 运行时: Docker
- 健康检查: `/health`
- 默认时区: `Asia/Shanghai`
- 必填环境变量: `DATABASE_URL`
- 生产数据库强制开关: `REQUIRE_DATABASE_URL=true`

部署到 Render 时, 创建 PostgreSQL 数据库后把连接串填入 `DATABASE_URL`. Telegram Bot Token、Chat ID、关键词规则、过滤词和定时参数可以在 Web 面板保存.

### Web 接口

主要接口:

- `GET /`: Web 面板首页
- `GET /health`: 健康检查
- `GET /productConfig`: 读取商品监控配置
- `POST /productConfig`: 保存商品监控配置
- `POST /productSearch`: 按单条规则预览搜索结果
- `GET /conf`: 读取已保存的签到用户
- `POST /addConf`: 保存什么值得买签到 Cookie
- `POST /check`: 执行签到
- `GET /imageProxy`: 商品图片代理

### 运行产物

- `pushed.json`: 已推送文章记录, 用于去重. 程序会自动创建.
- `data/users.db`: 本地 SQLite 数据库, 保存签到用户和 Web 面板配置.

### 效果

![Web 面板](https://img.ggball.top/picGo/20231024210905.png)

![推送效果](https://img.ggball.top/picGo/image-20220419205742369.png)

![商品列表](https://img.ggball.top/picGo/image-20220419205914792.png)

![签到配置](https://img.ggball.top/picGo/20220428194347.png)
