package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"ggball.com/smzdm/check_in"
	"ggball.com/smzdm/db"
	"ggball.com/smzdm/file"
	"ggball.com/smzdm/push"
	"ggball.com/smzdm/smzdm"
	"github.com/robfig/cron"
)

var conf = file.Config{}
var confMu sync.RWMutex
var checks = []file.CheckInfo{}
var userDbPath = "data/users.db"
var checkInCron *cron.Cron
var checkInCronMu sync.Mutex

func main() {

	go cronForProduct()
	go cronForCheckIn()

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
		time.Sleep(time.Duration(tickTime) * time.Second)
		requestSmzdm()
	}
}

// 每天定时打卡任务开启
func cronForCheckIn() {
	resetCheckInCron(currentConfig().Cron)
}

func resetCheckInCron(schedule string) {
	checkInCronMu.Lock()
	defer checkInCronMu.Unlock()

	if checkInCron != nil {
		checkInCron.Stop()
	}

	c := cron.New()
	if err := c.AddFunc(schedule, func() {
		chekIn, err := check_in.NewCheckIn(userDbPath)
		if err != nil {
			log.Fatal("Failed to initialize check-in service:", err)
		}
		chekIn.SetConfig(currentConfig(), checks)
		chekIn.CheckInAllUsers()
	}); err != nil {
		log.Printf("签到 Cron 配置无效: %v", err)
		return
	}
	c.Start()
	checkInCron = c
}

func validateCronSchedule(schedule string) error {
	if len(strings.Fields(schedule)) != 6 {
		return fmt.Errorf("必须是 6 段 cron 表达式")
	}
	c := cron.New()
	return c.AddFunc(schedule, func() {})
}

// 推送商品任务
func requestSmzdm() {
	// 搜索商品
	satisfyGoodsList, satisfyGoodsMyselfList := smzdm.GetSatisfiedGoods(currentConfig())
	if len(satisfyGoodsList) == 0 {
		return
	}
	// 推送商品
	push.PushProducts(satisfyGoodsList, currentConfig())
	// 推送自己关注的商品
	atMobiles := []string{"13217913287"}
	push.PushTargetProducts(satisfyGoodsMyselfList, currentConfig(), atMobiles)
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
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/html/", HtmlHandler)
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
