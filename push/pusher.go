package push

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"ggball.com/smzdm/file"
	"ggball.com/smzdm/smzdm"
)

var telegramEndpointBase = "https://api.telegram.org"
var httpClient = http.DefaultClient
var logMu sync.Mutex
var recentLogs []LogEntry

type LogEntry struct {
	Title  string `json:"title"`
	Status string `json:"status"`
	Reason string `json:"reason"`
	Image  string `json:"image"`
	Time   string `json:"time"`
}

func PushProducts(pro []smzdm.Product, conf file.Config) {
	PushProWithTelegram(pro, conf)
}

func PushText(resText string, conf file.Config) {
	PushTextWithTelegram(resText, conf)
}

func PushTargetProducts(pro []smzdm.Product, conf file.Config, atMobiles []string) {
	PushTargetWithTelegram(pro, conf)
}

func PushProWithTelegram(pro []smzdm.Product, conf file.Config) {
	if len(pro) == 0 {
		return
	}
	if !canPushTelegram(conf) {
		AddLog(LogEntry{
			Title:  "扫描有命中但未推送",
			Status: "skip",
			Reason: "Telegram 未启用或未配置 Bot Token / Chat ID",
		})
		return
	}

	// SatisfyNum = max products per Telegram message. Overflow is sent in follow-up messages.
	batchSize := conf.SatisfyNum
	if batchSize <= 0 {
		batchSize = 5
	}
	totalBatches := (len(pro) + batchSize - 1) / batchSize
	delivered := make([]smzdm.Product, 0, len(pro))

	for batchNo := 0; batchNo < totalBatches; batchNo++ {
		start := batchNo * batchSize
		end := start + batchSize
		if end > len(pro) {
			end = len(pro)
		}
		batch := pro[start:end]

		title := "什么值得买好价"
		if totalBatches > 1 {
			title = fmt.Sprintf("什么值得买好价（%d/%d）", batchNo+1, totalBatches)
		}
		content := buildProductsHTML(batch, title)
		err := sendTelegram(content, conf)
		for _, product := range batch {
			if err != nil {
				logProduct(product, "fail", err.Error())
				continue
			}
			reason := "已推送"
			if totalBatches > 1 {
				reason = fmt.Sprintf("已推送 · 第%d/%d批", batchNo+1, totalBatches)
			}
			logProduct(product, "success", reason)
			delivered = append(delivered, product)
		}
		if err != nil {
			fmt.Println("Telegram push failed:", err)
		}
		// Mild delay between batches to avoid Telegram flood limits.
		if batchNo+1 < totalBatches {
			time.Sleep(400 * time.Millisecond)
		}
	}

	// Only mark successfully delivered IDs so a failed send can be retried next scan.
	if len(delivered) > 0 {
		smzdm.MarkPushed(delivered)
	}
}

func PushTargetWithTelegram(pro []smzdm.Product, conf file.Config) {
	if len(pro) == 0 || !canPushTelegram(conf) {
		return
	}
	content := buildProductsHTML(pro, "好物到了")
	pushTelegram(content, conf)
}

func PushTextWithTelegram(resText string, conf file.Config) {
	if strings.TrimSpace(resText) == "" || !canPushTelegram(conf) {
		return
	}
	pushTelegram(html.EscapeString(resText)+"\n\n什么值得买", conf)
}

func SendTelegramTest(conf file.Config) error {
	conf.Telegram.Enabled = true
	err := sendTelegram("<b>什么值得买提醒</b>\nTelegram 测试发送成功。", conf)
	if err != nil {
		AddLog(LogEntry{Title: "Telegram 测试发送", Status: "fail", Reason: err.Error()})
		return err
	}
	AddLog(LogEntry{Title: "Telegram 测试发送", Status: "success", Reason: "已推送"})
	return nil
}

func pushTelegram(content string, conf file.Config) {
	if err := sendTelegram(content, conf); err != nil {
		fmt.Println("Telegram push failed:", err)
	}
}

func RecentLogs() []LogEntry {
	logMu.Lock()
	defer logMu.Unlock()
	result := make([]LogEntry, len(recentLogs))
	copy(result, recentLogs)
	return result
}

func AddLog(entry LogEntry) {
	if strings.TrimSpace(entry.Time) == "" {
		entry.Time = time.Now().Format("15:04:05")
	}
	if strings.TrimSpace(entry.Status) == "" {
		entry.Status = "success"
	}
	logMu.Lock()
	defer logMu.Unlock()
	recentLogs = append([]LogEntry{entry}, recentLogs...)
	if len(recentLogs) > 80 {
		recentLogs = recentLogs[:80]
	}
}

func logProduct(product smzdm.Product, status string, reason string) {
	title := strings.TrimSpace(product.ArticleTitle)
	if title == "" {
		title = "商品命中通知"
	}
	AddLog(LogEntry{
		Title:  title,
		Status: status,
		Reason: reason,
		Image:  product.ArticlePic,
	})
}

func sendTelegram(content string, conf file.Config) error {
	if strings.TrimSpace(conf.Telegram.BotToken) == "" {
		return fmt.Errorf("bot token 不能为空")
	}
	if strings.TrimSpace(conf.Telegram.ChatID) == "" {
		return fmt.Errorf("chat id 不能为空")
	}
	params := TelegramMessageParam{
		ChatID:                conf.Telegram.ChatID,
		Text:                  content,
		ParseMode:             telegramParseMode(conf.Telegram.ParseMode),
		DisableWebPagePreview: conf.Telegram.DisableWebPagePreview,
	}

	paramsJson, _ := json.Marshal(params)
	endpoint := strings.TrimRight(telegramEndpointBase, "/") + "/bot" + strings.TrimSpace(conf.Telegram.BotToken) + "/sendMessage"
	resp, err := httpClient.Post(endpoint, "application/json;charset=utf-8", bytes.NewBuffer(paramsJson))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	contentBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram response status %d: %s", resp.StatusCode, string(contentBytes))
	}
	fmt.Println(string(contentBytes))
	return nil
}

func canPushTelegram(conf file.Config) bool {
	return conf.Telegram.Enabled &&
		strings.TrimSpace(conf.Telegram.BotToken) != "" &&
		strings.TrimSpace(conf.Telegram.ChatID) != ""
}

func buildProductsHTML(pro []smzdm.Product, title string) string {
	var builder strings.Builder
	builder.WriteString("<b>" + html.EscapeString(title) + "</b>\n")
	for _, item := range pro {
		builder.WriteString("\n")
		builder.WriteString(formatProductHTML(item))
	}
	return builder.String()
}

func formatProductHTML(item smzdm.Product) string {
	title := html.EscapeString(item.ArticleTitle)
	link := html.EscapeString(item.ArticleUrl)
	price := html.EscapeString(item.ArticlePrice)
	comment := html.EscapeString(item.ArticleComment)
	worthy := html.EscapeString(item.ArticleWorthy)
	return fmt.Sprintf("<a href=\"%s\">%s</a>\n%s, 评论 %s, 值率 %s%%\n", link, title, price, comment, worthy)
}

func telegramParseMode(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "HTML"
	}
	return value
}

// productLimit is kept for tests/helpers: clamp batch size to available items.
func productLimit(length int, satisfyNum int) int {
	if satisfyNum <= 0 || satisfyNum > length {
		return length
	}
	return satisfyNum
}
