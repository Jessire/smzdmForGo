package main

import (
	"encoding/json"
	"fmt"
	"ggball.com/smzdm/check_in"
	"ggball.com/smzdm/db"
	"ggball.com/smzdm/file"
	"ggball.com/smzdm/push"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type addCheckInfoRequest struct {
	Remark string `json:"remark"`
	Cookie string `json:"cookie"`
}

type checkInRequest struct {
	ID     int    `json:"id"`
	Remark string `json:"remark"`
	Cookie string `json:"cookie"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		t, err := template.ParseFiles("template/html/index.html")
		if err != nil {
			log.Println(err)
		}
		t.Execute(w, nil)
	}

}

func HtmlHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	t, err := template.ParseFiles("template/" + r.URL.Path + ".html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}

func ReadCheckInfoHandler(w http.ResponseWriter, r *http.Request) {
	checks, err := readCheckInfos()
	if err != nil {
		writeError(w, err)
		return
	}
	jsonByte, _ := json.Marshal(checks)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(wrapDataWithResult(string(jsonByte))))
	// fmt.Println(wrapDataWithResult(string(jsonByte)))
}

func AddCheckInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 读取添加的数据
	body, _ := ioutil.ReadAll(r.Body)
	var req addCheckInfoRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, fmt.Errorf("解析新增配置失败: %v", err))
		return
	}
	if req.Cookie == "" {
		writeError(w, fmt.Errorf("cookie 不能为空"))
		return
	}

	database, err := openUserDB()
	if err != nil {
		writeError(w, err)
		return
	}
	defer database.Close()

	user := db.User{
		Name:     req.Remark,
		Token:    req.Cookie,
		Platform: "smzdm",
	}
	if err := database.AddUser(&user); err != nil {
		writeError(w, err)
		return
	}
	w.Write([]byte(wrapDataWithResult("\"" + "添加成功" + "\"")))
	// fmt.Println(checks)
}

func CheckInHandler(w http.ResponseWriter, r *http.Request) {
	// 读取添加的数据
	body, _ := ioutil.ReadAll(r.Body)
	reqs, err := parseCheckInRequests(body)
	if err != nil {
		writeError(w, fmt.Errorf("解析签到请求失败: %v", err))
		return
	}
	if len(reqs) == 0 {
		writeError(w, fmt.Errorf("签到请求不能为空"))
		return
	}
	checkInfo := reqs[0]
	fmt.Println("checkInfo:", checkInfo)
	if checkInfo.Cookie == "" {
		writeError(w, fmt.Errorf("cookie 不能为空"))
		return
	}

	checker, err := check_in.NewCheckIn(userDbPath)
	if err != nil {
		writeError(w, err)
		return
	}
	defer checker.Close()
	checker.SetConfig(currentConfig(), checks)
	msg, err := checker.CheckInUser(db.User{
		ID:       int64(checkInfo.ID),
		Name:     checkInfo.Remark,
		Token:    checkInfo.Cookie,
		Platform: "smzdm",
	})
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(wrapDataWithResult("\"" + msg + "\"")))

}

func ProductConfigHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"code": "0",
			"msg":  "",
			"data": productConfigFromConfig(currentConfig()),
		})
	case http.MethodPost:
		body, _ := ioutil.ReadAll(r.Body)
		var req productConfigRequest
		if err := json.Unmarshal(body, &req); err != nil {
			writeError(w, fmt.Errorf("解析商品规则失败: %v", err))
			return
		}
		next := req.applyTo(currentConfig())
		file.ApplyEnvOverrides(&next)

		database, err := openUserDB()
		if err != nil {
			writeError(w, err)
			return
		}
		defer database.Close()
		if err := database.SaveProductConfig(next); err != nil {
			writeError(w, err)
			return
		}
		setCurrentConfig(next)
		notifyProductScheduleChanged()
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"code": "0",
			"msg":  "保存成功",
			"data": productConfigFromConfig(next),
		})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func TelegramTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var req telegramTestRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, fmt.Errorf("解析 Telegram 测试配置失败: %v", err))
		return
	}
	if err := sendTelegramTest(req); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code": "0",
		"msg":  "测试发送成功",
		"data": map[string]interface{}{},
	})
}

func PushLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code": "0",
		"msg":  "",
		"data": push.RecentLogs(),
	})
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// HealthDBHandler pings Postgres so free-tier providers see activity.
func HealthDBHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := pingDatabase(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		writeJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"status": "db_down",
			"error":  err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","db":"up"}`))
}

func pingDatabase() error {
	database, err := db.NewDB(userDbPath)
	if err != nil {
		return err
	}
	defer database.Close()
	return database.Ping()
}

func keepaliveDatabase() {
	// First ping shortly after boot, then every 6 hours.
	time.Sleep(30 * time.Second)
	for {
		if err := pingDatabase(); err != nil {
			log.Printf("db keepalive failed: %v", err)
		} else {
			log.Printf("db keepalive ok")
		}
		time.Sleep(6 * time.Hour)
	}
}

func readCheckInfoJson() []byte {
	// 打开json文件
	jsonFile, err := os.Open("template/json/checkInfo.json")

	// 最好要处理以下错误
	if err != nil {
		fmt.Println(err)
	}

	// 要记得关闭
	defer jsonFile.Close()

	jsonByte, _ := ioutil.ReadAll(jsonFile)
	return jsonByte
}

func readCheckInfos() ([]file.CheckInfo, error) {
	database, err := openUserDB()
	if err != nil {
		return nil, err
	}
	defer database.Close()

	users, err := database.GetAllUsers()
	if err != nil {
		return nil, err
	}

	checks := make([]file.CheckInfo, 0, len(users))
	for _, user := range users {
		checks = append(checks, file.CheckInfo{
			Id:         int(user.ID),
			LastTIme:   user.LastTime,
			Remark:     user.Name,
			LastMsg:    user.LastMsg,
			LastResult: user.LastResult,
			Cookie:     user.Token,
		})
	}
	return checks, nil
}

func openUserDB() (*db.DB, error) {
	database, err := db.NewDB(userDbPath)
	if err != nil {
		return nil, err
	}
	if err := database.InitTables(); err != nil {
		database.Close()
		return nil, err
	}
	return database, nil
}

func deserializeJson(CheckInfoJson string) []file.CheckInfo {
	// fmt.Println("CheckInfoJson:", CheckInfoJson)
	jsonAsBytes := []byte(CheckInfoJson)
	checks := make([]file.CheckInfo, 0)
	err := json.Unmarshal(jsonAsBytes, &checks)
	// fmt.Printf("%#v", checks)
	if err != nil {
		panic(err)
	}
	return checks
}

func parseCheckInRequests(body []byte) ([]checkInRequest, error) {
	var reqs []checkInRequest
	if err := json.Unmarshal(body, &reqs); err == nil {
		return reqs, nil
	}

	var req checkInRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	return []checkInRequest{req}, nil
}

func wrapDataWithResult(data string) string {

	result := `
	{"code":"0",
	"msg":   "",
	"count": "10",
	"data":  ` + data + `}`

	return result
}

func writeError(w http.ResponseWriter, err error) {
	log.Println(err)
	writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
		"code":  "1",
		"msg":   err.Error(),
		"count": "0",
		"data":  []interface{}{},
	})
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
