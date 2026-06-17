package push

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"ggball.com/smzdm/file"
	"ggball.com/smzdm/smzdm"
)

var wxPusherEndpoint = "https://wxpusher.zjiecode.com/api/send/message"
var httpClient = http.DefaultClient

// 定义推送者，声明推送方法
type Pusher interface {
	Push(content string, contentType string)
}

type DingPusher struct {
	Token string
}

// 钉钉推送者实现推送方法
func (pusher DingPusher) PushDingDing(params interface{}) {
	if strings.TrimSpace(pusher.Token) == "" {
		return
	}
	Url, err := url.Parse("https://oapi.dingtalk.com/robot/send?access_token=" + pusher.Token)
	if err != nil {
		return
	}

	paramsJson, _ := json.Marshal(params)
	fmt.Println(string(paramsJson))
	urlPath := Url.String()
	resp, err := http.Post(urlPath, "application/json;charset=utf-8", bytes.NewBuffer([]byte(string(paramsJson))))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	fmt.Println(string(content))

}

func PushProducts(pro []smzdm.Product, conf file.Config) {
	PushProWithDingDing(pro, conf)
	PushProWithWxPusher(pro, conf)
}

func PushText(resText string, conf file.Config) {
	PushTextWithDingDing(resText, conf)
	PushTextWithWxPusher(resText, conf)
}

func PushTargetProducts(pro []smzdm.Product, conf file.Config, atMobiles []string) {
	PushTextWithDingDingWIthMoblie(pro, conf, atMobiles)
	PushTargetWithWxPusher(pro, conf)
}

// 推送商品到钉钉
func PushProWithDingDing(pro []smzdm.Product, conf file.Config) {
	if len(pro) == 0 || strings.TrimSpace(conf.DingdingToken) == "" {
		return
	}
	dingPusher := DingPusher{
		Token: conf.DingdingToken,
	}

	// 需要提前申明数组的容量
	limit := productLimit(len(pro), conf.SatisfyNum)
	links := make([]Link, 0, limit)

	for index := 0; index < limit; index++ {
		item := pro[index]
		link := Link{
			Title:      formatProductTitle(item),
			MessageURL: item.ArticleUrl,
			PicURL:     item.ArticlePic,
		}
		links = append(links, link)
	}
	fmt.Printf("links:%#v", links)

	feedCard := FeedCard{
		Links: links,
	}

	params := DingFeedCardParam{
		MsgType:  "feedCard",
		FeedCard: feedCard,
	}

	dingPusher.PushDingDing(params)
}

// 推送文字到钉钉
func PushTextWithDingDing(resText string, conf file.Config) {
	if strings.TrimSpace(resText) == "" || strings.TrimSpace(conf.DingdingToken) == "" {
		return
	}
	dingPusher := DingPusher{
		Token: conf.DingdingToken,
	}

	text := Text{
		Content: resText + "【什么值得买】",
	}

	params := DingTextParam{
		MsgType: "text",
		Texts:   text,
	}

	dingPusher.PushDingDing(params)
}

func PushTextWithWxPusher(resText string, conf file.Config) {
	if strings.TrimSpace(resText) == "" || !canPushWx(conf) {
		return
	}
	pushWxPusher(resText+"【什么值得买】", "什么值得买", "", conf)
}

// 推送文字到钉钉并@人
func PushTextWithDingDingWIthMoblie(pro []smzdm.Product, conf file.Config, atMobiles []string) {

	if len(pro) == 0 || strings.TrimSpace(conf.DingdingToken) == "" {
		return
	}

	dingPusher := DingPusher{
		Token: conf.DingdingToken,
	}

	title := "【好物到了】 \n"
	text := ""
	for _, item := range pro {
		text += formatProductMarkdown(item)
	}
	md := Markdown{Title: title, Text: text}
	params := DingMdParam{
		MsgType:  "markdown",
		Markdown: md,
	}

	textParams := DingTextParam{
		MsgType: "text",
		Texts:   Text{Content: title},
		At:      At{AtMobiles: atMobiles, IsAtAll: false},
	}

	dingPusher.PushDingDing(textParams)
	dingPusher.PushDingDing(params)
}

func PushProWithWxPusher(pro []smzdm.Product, conf file.Config) {
	if len(pro) == 0 || !canPushWx(conf) {
		return
	}
	limit := productLimit(len(pro), conf.SatisfyNum)
	content := buildProductsMarkdown(pro[:limit], "什么值得买好价")
	pushWxPusher(content, "什么值得买好价", firstProductURL(pro), conf)
}

func PushTargetWithWxPusher(pro []smzdm.Product, conf file.Config) {
	if len(pro) == 0 || !canPushWx(conf) {
		return
	}
	content := buildProductsMarkdown(pro, "好物到了")
	pushWxPusher(content, "好物到了", firstProductURL(pro), conf)
}

func pushWxPusher(content string, summary string, link string, conf file.Config) {
	contentType := conf.WxPusher.ContentType
	if contentType == 0 {
		contentType = 3
	}
	params := WxPusherMessageParam{
		AppToken:    conf.WxPusher.AppToken,
		Content:     content,
		Summary:     summary,
		ContentType: contentType,
		UIDs:        conf.WxPusher.UIDs,
		TopicIDs:    conf.WxPusher.TopicIDs,
		URL:         link,
	}

	paramsJson, _ := json.Marshal(params)
	resp, err := httpClient.Post(wxPusherEndpoint, "application/json;charset=utf-8", bytes.NewBuffer(paramsJson))
	if err != nil {
		fmt.Println("WxPusher push failed:", err)
		return
	}
	defer resp.Body.Close()

	contentBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("WxPusher response read failed:", err)
		return
	}
	fmt.Println(string(contentBytes))
}

func canPushWx(conf file.Config) bool {
	return conf.WxPusher.Enabled &&
		strings.TrimSpace(conf.WxPusher.AppToken) != "" &&
		(len(conf.WxPusher.UIDs) > 0 || len(conf.WxPusher.TopicIDs) > 0)
}

func buildProductsMarkdown(pro []smzdm.Product, title string) string {
	var builder strings.Builder
	builder.WriteString("## " + title + "\n\n")
	for _, item := range pro {
		if strings.TrimSpace(item.ArticlePic) != "" {
			builder.WriteString("![](" + item.ArticlePic + ")\n\n")
		}
		builder.WriteString(formatProductMarkdown(item))
		builder.WriteString("\n")
	}
	return builder.String()
}

func formatProductMarkdown(item smzdm.Product) string {
	return fmt.Sprintf("- [**%s**](%s): %s, 评论 %s, 值率 %s%%\n", item.ArticleTitle, item.ArticleUrl, item.ArticlePrice, item.ArticleComment, item.ArticleWorthy)
}

func formatProductTitle(item smzdm.Product) string {
	return item.ArticlePrice + "!【" + item.ArticleTitle + "】评论 " + item.ArticleComment + " 值率 " + item.ArticleWorthy + "%【什么值得买】"
}

func productLimit(length int, satisfyNum int) int {
	if satisfyNum <= 0 || satisfyNum > length {
		return length
	}
	return satisfyNum
}

func firstProductURL(pro []smzdm.Product) string {
	if len(pro) == 0 {
		return ""
	}
	return pro[0].ArticleUrl
}
