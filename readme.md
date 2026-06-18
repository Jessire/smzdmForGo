## 什么值得买文章推送器
-----

### 目的

**smzdmForGo** 是一个基于 Go 的信息监控、去重与通知工具，用于爬取和过滤「什么值得买」上的商品文章，帮助用户及时获取心仪商品的折扣信息。同时本项目也是 Go 语言学习的实践项目。

### 已实现
- [x] 自定义文章提取规则
- [x] 推送文章
- [x]  去重文章
- [x]  定时推送
- [x] 设定关键字，爬取含关键词的商品
- [x] 利用github Action 自动编译，部署到个人服务器
- [x] 每天定时打卡
- [x] 支持 Telegram Bot 推送
- [x] 支持关键词级过滤词、评论数、值率、价格规则
- [x] 支持 Web 面板保存商品规则、定时参数和通知配置
  
### 待实现
- [ ] 配置server酱
- [X] 配置签到

### 使用步骤
下载整个代码 window平台直接运行`smzdm.exe`，切勿挪动exe文件，会导致读不到配置
推荐直接打开 Web 面板维护关键词、过滤词、阈值、Telegram Bot 和签到账号配置。保存后的生产配置会写入数据库。
1. **配置式：**
修改以下配置，保存配置，再运行`smzdm.exe`即可
```yml
# 搜索关键词
keyWord: 
- 信小兔
- 零食

# 最低评论数
lowCommentNum: 0
# 最低值率
lowWorthyNum: 0
# 最低价格, 0 表示不限制
minPrice: 0
# 最高价格, 0 表示不限制
maxPrice: 0
# 满意商品数量
satisfyNum: 10
# 过滤词
filterWords: 
- "榴莲"
- "唯品会"
- "牛奶"
- "电脑"

# 每组关键词的独立规则
keywordRules:
  - words:
      - 显示器
    filterWords:
      - 二手
    lowCommentNum: 5
    lowWorthyNum: 20
    minPrice: 300
    maxPrice: 2000

# 定时任务多长执行一次 单位秒 默认 12个小时
tickTime: 43200
# Telegram Bot 推送
telegram:
  enabled: true
  botToken: "123456:ABC_xxxxx"
  chatId: "123456789"
  parseMode: "HTML"
  disableWebPagePreview: false

# 签到时间(默认早上8:30)
cron: "0 30 8 ? * *"

# 签到需要的cookie
cookie: "XXXX"

```

Render 部署只需要配置数据库连接和 `REQUIRE_DATABASE_URL=true`。Telegram Bot Token、Chat ID、关键词规则、过滤词和定时参数都在 Web 面板里保存。
2. **docker方式**

- 执行 `docker pull registry.cn-hangzhou.aliyuncs.com/ggball/smzdm_for_go:latest`
- 再创建配置目录`D:\\documents\\config`（我这里实在win下操作的），将`config/config.yml`文件 放入创建好的配置文件夹
![20220606095621](https://img.ggball.top/picGo/20220606095621.png)
- 最后执行`docker run -d --name smzdmForgo -v D:\\documents\\config:/opt/go/config registry.cn-hangzhou.aliyuncs.com/ggball/smzdm_for_go:latest`

> -v :前是宿主机目录，:后是容器目录

3. 源码启动
```go 
go run .\main.go .\route.go
```

**配置签到**
可手动和定时签到
[配置地址](http://1.15.141.114:9090/)
![20220810003120](https://img.ggball.top/picGo/20220810003120.png)





### 效果
![20231024210905](https://img.ggball.top/picGo/20231024210905.png)
![image-20220419205742369](https://img.ggball.top/picGo/image-20220419205742369.png)

![image-20220419205914792](https://img.ggball.top/picGo/image-20220419205914792.png)

![20220428194347](https://img.ggball.top/picGo/20220428194347.png)
