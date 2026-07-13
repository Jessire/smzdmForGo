package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"ggball.com/smzdm/db"
	"ggball.com/smzdm/file"
	"ggball.com/smzdm/push"
	"ggball.com/smzdm/smzdm"
)

var conf = file.Config{}
var confMu sync.RWMutex
var checks = []file.CheckInfo{}
var userDbPath = "data/users.db"
var productScheduleChanged = make(chan struct{}, 1)

func main() {

	go cronForProduct()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	fmt.Println("启动web服务，监听" + port + "端口")
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func cronForProduct() {
	for {
		tickTime := currentConfig().TickTime
		if tickTime <= 0 {
			tickTime = 10800
		}
		timer := time.NewTimer(time.Duration(tickTime) * time.Second)
		select {
		case <-timer.C:
			requestSmzdm()
		case <-productScheduleChanged:
			// Config saved: cancel the pending wait, scan+push once now, then restart the interval.
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			requestSmzdm()
		}
	}
}

var productScanMu sync.Mutex

// 推送商品任务
func requestSmzdm() {
	productScanMu.Lock()
	defer productScanMu.Unlock()

	// 搜索商品（已推送过 / 超过约2天 的会在这里被去掉）
	satisfyGoodsList, _ := smzdm.GetSatisfiedGoods(currentConfig())
	if len(satisfyGoodsList) == 0 {
		log.Printf("商品扫描完成：无新命中，跳过推送")
		// Surface to UI so "save but no Telegram" is explainable.
		age := currentConfig().MaxArticleAgeDays
		reason := "没有新命中（可能已推送过或未过门槛）"
		if age > 0 {
			reason = fmt.Sprintf("没有新命中（可能已推送过、超过近%d天或未过门槛）", age)
		}
		push.AddLog(push.LogEntry{
			Title:  "扫描完成 · 无新推送",
			Status: "skip",
			Reason: reason,
		})
		return
	}
	// 推送商品（成功后才写入已推送去重）
	push.PushProducts(satisfyGoodsList, currentConfig())
	time.Sleep(1 * time.Second)
}

func init() {

	// 读取项目根目录的配置文件
	conf = file.ReadConf("")
	loadSavedProductConfig()
	checks = file.ReadCheckInfoJsonToCheck()

	// 配置路由
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/conf", ReadCheckInfoHandler)
	http.HandleFunc("/addConf", AddCheckInfoHandler)
	http.HandleFunc("/check", CheckInHandler)
	http.HandleFunc("/productConfig", ProductConfigHandler)
	http.HandleFunc("/productSearch", ProductSearchHandler)
	http.HandleFunc("/imageProxy", ImageProxyHandler)
	http.HandleFunc("/telegramTest", TelegramTestHandler)
	http.HandleFunc("/telegramAvatar", TelegramAvatarHandler)
	http.HandleFunc("/pushLogs", PushLogsHandler)
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/health/db", HealthDBHandler)
	http.HandleFunc("/html/", HtmlHandler)

	// Lightweight DB keepalive so free Aiven Postgres is less likely to power off.
	go keepaliveDatabase()
}

func currentConfig() file.Config {
	confMu.RLock()
	defer confMu.RUnlock()
	return conf
}

func setCurrentConfig(next file.Config) {
	confMu.Lock()
	defer confMu.Unlock()
	conf = next
}

func notifyProductScheduleChanged() {
	select {
	case productScheduleChanged <- struct{}{}:
	default:
	}
}

func loadSavedProductConfig() {
	database, err := db.NewDB(userDbPath)
	if err != nil {
		log.Printf("读取数据库配置失败: %v", err)
		return
	}
	defer database.Close()
	if err := database.InitTables(); err != nil {
		log.Printf("初始化数据库配置表失败: %v", err)
		return
	}
	next, err := database.GetProductConfig(conf)
	if err != nil {
		log.Printf("读取商品规则配置失败: %v", err)
		return
	}
	setCurrentConfig(next)
}
