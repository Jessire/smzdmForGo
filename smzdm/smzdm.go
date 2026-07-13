package smzdm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"ggball.com/smzdm/file"
)

type result struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	Data      Data   `json:"data"`
}

type Data struct {
	Rows  []Product `json:"rows"`
	Total int       `json:"total"`
}

type Product struct {
	ArticleTitle   string `json:"article_title"`
	ArticlePrice   string `json:"article_price"`
	ArticleWorthy  string `json:"article_worthy"`
	ArticleComment string `json:"article_comment"`
	ArticleId      string `json:"article_id"`
	ArticleDate    string `json:"publish_date_lt"`
	ArticlePic     string `json:"article_pic"`
	ArticleUrl     string `json:"article_url"`
	Referral       string `json:"article_referrals"`
}

// 全局配置
var globalConf = file.Config{}

// 推送信息文件地址
var pushedPath = "./pushed.json"

// 获取商品
//
//	@return []product 符合条件的商品集合
//	@return []product 符合自己条件的商品集合
func GetSatisfiedGoods(conf file.Config) ([]Product, []Product) {
	globalConf = conf
	fmt.Println("开始爬取符合条件商品。。")

	// 获取已推送文章id
	pushedMap := file.ReadPusedInfo(pushedPath)

	// 符合条件的商品集合
	var satisfyGoodsList []Product

	// 符合自己条件的商品集合
	var satisfyGoodsListBySelf []Product

	if len(conf.KeywordRules) > 0 {
		satisfyGoodsList = getSatisfiedGoodsByKeywordRules(conf, pushedMap)
	} else {
		satisfyGoodsList = getSatisfiedGoodsFromFeed(pushedMap)
	}

	// 根据评论数排序
	sort.SliceStable(satisfyGoodsList, func(a, b int) bool {
		return parseMetric(satisfyGoodsList[a].ArticleComment) > parseMetric(satisfyGoodsList[b].ArticleComment)
	})

	fmt.Println("结束爬取符合条件商品。。")

	//过滤出自己的商品
	satisfyGoodsListBySelf = filterMyselfProduct(satisfyGoodsList)

	// Do NOT mark as pushed here — only after Telegram send succeeds (see MarkPushed).
	return satisfyGoodsList, satisfyGoodsListBySelf
}

// MarkPushed records successfully delivered product IDs so they won't be re-sent.
func MarkPushed(goods []Product) {
	if len(goods) == 0 {
		return
	}
	pushedMap := file.ReadPusedInfo(pushedPath)
	savePushed(pushedMap, pushedPath, goods)
}

func getSatisfiedGoodsFromFeed(pushedMap map[string]interface{}) []Product {
	var satisfyGoodsList []Product
	page := 0
	for {

		var productList = []Product{}
		// Get the good list
		productList = GetGoods(page, "").Data.Rows

		// add satisfy good
		if len(productList) > 0 {
			rows := productList
			for i := 0; i < len(rows); i++ {
				good := rows[i]

				// 商品 包含 “k” 转换数字 默认给1000
				if strings.Contains(strings.ToLower(good.ArticleComment), "k") {
					good.ArticleComment = "1000"
				}

				if removeByFilterRules(good, pushedMap) {
					continue
				}

				if satisfy(good, satisfyGoodsList) {
					satisfyGoodsList = append(satisfyGoodsList, good)
				}

			}
		}

		// 页数+1
		page++
		// 延时2s
		time.Sleep(time.Duration(2) * time.Second)

		// 判断是否退出
		if shouldStop(len(satisfyGoodsList), page) {
			fmt.Println("退出")
			break
		}

	}

	return satisfyGoodsList
}

func getSatisfiedGoodsByKeywordRules(conf file.Config, pushedMap map[string]interface{}) []Product {
	seen := map[string]bool{}
	limit := conf.SatisfyNum * 3
	if limit < 24 {
		limit = 24
	}
	var satisfyGoodsList []Product
	for _, rule := range conf.KeywordRules {
		if !keywordRuleEnabled(rule) {
			continue
		}
		for _, keyword := range rule.Words {
			for _, good := range SearchGoods(keyword, rule, limit) {
				if removePushedOrOld(good, pushedMap) || seen[good.ArticleId] {
					continue
				}
				seen[good.ArticleId] = true
				satisfyGoodsList = append(satisfyGoodsList, good)
			}
		}
	}
	return satisfyGoodsList
}

// GetGoods 获取商品集合
//
//	@param offset
//	@return result 商品集合
func GetGoods(page int, keword string) result {

	var res result

	params := url.Values{}
	Url, err := url.Parse("https://api.smzdm.com/v1/list")
	if err != nil {
		return res
	}
	params.Set("keyword", keword)
	// score 值率排序  time 时间排序
	params.Set("order", "time")
	params.Set("type", "good_price")
	params.Set("offset", strconv.Itoa(page*100))
	params.Set("limit", "100")

	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Println(urlPath)
	resp, err := http.Get(urlPath)
	if err != nil {
		return res
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))

	_ = json.Unmarshal(body, &res)
	// fmt.Printf("%#v", res)
	return res

}

// 根据条件 判断是否应该停止爬取
func shouldStop(length int, page int) bool {
	fmt.Println("length:" + strconv.Itoa(length) + "\n\r page:" + strconv.Itoa(page))
	//  判断数量是否超过【符合商品个数】 且 page > 20
	return length > globalConf.SatisfyNum || page > 100

}

// 根据过滤规则，去除商品
func removeByFilterRules(good Product, pushedMap map[string]interface{}) bool {
	var noNeed = false
	// 1. 文章名称 包含过滤字符 一概不要
	for j := 0; j < len(globalConf.FilterWords); j++ {
		if containsWord(good, globalConf.FilterWords[j]) {
			noNeed = true
			break
		}
	}

	return noNeed || removePushedOrOld(good, pushedMap)
}

func removePushedOrOld(good Product, pushedMap map[string]interface{}) bool {
	return PushSkipReason(good, pushedMap) != ""
}

// LoadPushedMap loads the on-disk set of already-delivered article IDs.
func LoadPushedMap() map[string]interface{} {
	return file.ReadPusedInfo(pushedPath)
}

// PushSkipReason explains why a product would not be Telegram-pushed.
// Empty string means eligible for push (subject to rule match already applied).
// Values: already_pushed | too_old | bad_date
func PushSkipReason(good Product, pushedMap map[string]interface{}) string {
	if good.ArticleId != "" {
		if _, ok := pushedMap[good.ArticleId]; ok {
			return "already_pushed"
		}
	}

	// Only keep roughly last 2 days for background push (preview used to skip this).
	if good.ArticleDate == "" {
		return ""
	}
	dateInt64, err := strconv.ParseInt(good.ArticleDate, 10, 64)
	if err != nil {
		return "bad_date"
	}
	arDate := time.Unix(dateInt64, 0)
	if arDate.Before(time.Now().AddDate(0, 0, -2)) {
		return "too_old"
	}
	return ""
}

// 根据规则判断符合规则的商品
func satisfy(good Product, satisfyGoodsList []Product) bool {
	if len(globalConf.KeywordRules) > 0 {
		return matchesPersonalRules(good)
	}
	if !priceInRange(good, globalConf.MinPrice, globalConf.MaxPrice) {
		return false
	}

	articleComment := parseMetric(good.ArticleComment)
	articleWorthy := parseMetric(good.ArticleWorthy)

	// 评论，值率满足要求 则添加商品
	if articleComment >= globalConf.LowCommentNum || articleWorthy >= globalConf.LowWorthyNum {
		fmt.Printf("appear satisfy good: %#v", good)
		return true
	}

	return false
}

func SearchGoods(keyword string, rule file.KeywordRule, limit int) []Product {
	keyword = strings.TrimSpace(keyword)
	if limit <= 0 {
		limit = 20
	}
	if keyword == "" {
		return []Product{}
	}

	// SMZDM list API is plain keyword search (not the website ranking).
	// Official web results often surface mid-ranked deals (e.g. 今日必买) that
	// sit past the first 1–2 pages of /v1/list; with comment/worthy filters most
	// early rows are discarded, so we must scan deeper to fill the preview.
	// Each page is 100 rows → 8 pages ≈ 800 candidates max.
	maxPages := 8
	if limit > 40 {
		maxPages = 12
	}

	results := make([]Product, 0, limit)
	for page := 0; page < maxPages && len(results) < limit; page++ {
		rows := GetGoods(page, keyword).Data.Rows
		if len(rows) == 0 {
			break
		}
		for _, item := range rows {
			if searchProductMatches(item, rule) {
				results = append(results, item)
				if len(results) >= limit {
					break
				}
			}
		}
	}
	return results
}

// 保存推送商品，去重使用
func savePushed(pushedMap map[string]interface{}, pushedPath string, satisfyGoodsList []Product) {
	tempMap := make(map[string]interface{})

	for index, value := range satisfyGoodsList {
		tempMap[value.ArticleId] = index
	}
	file.WritePushedInfo(tempMap, pushedMap, pushedPath)
}

// 过滤自己的商品
func filterMyselfProduct(satisfyGoodsList []Product) []Product {

	var satisfyGoodsListBySelf []Product

	for _, value := range satisfyGoodsList {
		if matchesPersonalRules(value) {
			fmt.Printf("appear myself satisfy good: %#v", value)
			satisfyGoodsListBySelf = append(satisfyGoodsListBySelf, value)
		}
	}
	return satisfyGoodsListBySelf

}

func matchesPersonalRules(good Product) bool {
	if len(globalConf.KeywordRules) == 0 {
		return containsAnyWord(good, globalConf.KeyWords)
	}

	for _, rule := range globalConf.KeywordRules {
		if !keywordRuleEnabled(rule) {
			continue
		}
		if len(rule.Words) == 0 || !containsAnyWord(good, rule.Words) {
			continue
		}
		if containsAnyWord(good, rule.FilterWords) {
			continue
		}

		minPrice := globalConf.MinPrice
		if rule.MinPrice != nil {
			minPrice = *rule.MinPrice
		}
		maxPrice := globalConf.MaxPrice
		if rule.MaxPrice != nil {
			maxPrice = *rule.MaxPrice
		}
		if !priceInRange(good, minPrice, maxPrice) {
			continue
		}

		lowCommentNum := globalConf.LowCommentNum
		if rule.LowCommentNum != nil {
			lowCommentNum = *rule.LowCommentNum
		}
		lowWorthyNum := globalConf.LowWorthyNum
		if rule.LowWorthyNum != nil {
			lowWorthyNum = *rule.LowWorthyNum
		}
		if parseMetric(good.ArticleComment) < lowCommentNum {
			continue
		}
		if parseMetric(good.ArticleWorthy) < lowWorthyNum {
			continue
		}
		return true
	}
	return false
}

func searchProductMatches(good Product, rule file.KeywordRule) bool {
	if containsAnyWord(good, rule.FilterWords) {
		return false
	}

	minPrice := 0.0
	if rule.MinPrice != nil {
		minPrice = *rule.MinPrice
	}
	maxPrice := 0.0
	if rule.MaxPrice != nil {
		maxPrice = *rule.MaxPrice
	}
	if !priceInRange(good, minPrice, maxPrice) {
		return false
	}

	if rule.LowCommentNum != nil && parseMetric(good.ArticleComment) < *rule.LowCommentNum {
		return false
	}
	if rule.LowWorthyNum != nil && parseMetric(good.ArticleWorthy) < *rule.LowWorthyNum {
		return false
	}
	return true
}

func keywordRuleEnabled(rule file.KeywordRule) bool {
	return rule.Enabled == nil || *rule.Enabled
}

func priceInRange(good Product, minPrice float64, maxPrice float64) bool {
	if minPrice <= 0 && maxPrice <= 0 {
		return true
	}
	price, ok := parsePrice(good.ArticlePrice)
	if !ok {
		return false
	}
	if minPrice > 0 && price < minPrice {
		return false
	}
	if maxPrice > 0 && price > maxPrice {
		return false
	}
	return true
}

func containsAnyWord(good Product, words []string) bool {
	for _, word := range words {
		if containsWord(good, word) {
			return true
		}
	}
	return false
}

func containsWord(good Product, word string) bool {
	word = strings.TrimSpace(strings.ToLower(word))
	if word == "" {
		return false
	}
	return strings.Contains(strings.ToLower(good.ArticleTitle), word) ||
		strings.Contains(strings.ToLower(good.ArticlePrice), word) ||
		strings.Contains(strings.ToLower(good.Referral), word)
}

var numberPattern = regexp.MustCompile(`\d+(?:,\d{3})*(?:\.\d+)?`)

func parseMetric(value string) int {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return 0
	}
	matches := numberPattern.FindAllString(value, -1)
	if len(matches) == 0 {
		return 0
	}
	number, err := strconv.ParseFloat(strings.ReplaceAll(matches[0], ",", ""), 64)
	if err != nil {
		return 0
	}
	if strings.Contains(value, "万") {
		number *= 10000
	} else if strings.Contains(value, "k") {
		number *= 1000
	}
	return int(number)
}

func parsePrice(value string) (float64, bool) {
	value = strings.TrimSpace(value)
	match := numberPattern.FindString(value)
	if match == "" {
		return 0, false
	}
	price, err := strconv.ParseFloat(strings.ReplaceAll(match, ",", ""), 64)
	return price, err == nil
}
