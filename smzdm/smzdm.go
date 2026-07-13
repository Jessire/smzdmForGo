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

	// 保存推送商品，去重使用
	savePushed(pushedMap, pushedPath, satisfyGoodsList)

	return satisfyGoodsList, satisfyGoodsListBySelf
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
	// 根据已推送文章id map 判断是否需要去除，如果已经推送过的，则去除
	_, b := pushedMap[good.ArticleId]
	if b {
		// fmt.Println(good.ArticleTitle + "文章已存在,不予添加")
		return true
	}

	// 文章时间小于昨天 去除
	// var timeLayoutStr = "2006-01-02 15:04:05" //go中的时间格式化必须是这个时间
	nTime := time.Now()
	// 前天
	beforeYesDate := nTime.AddDate(0, 0, -2)
	if good.ArticleDate == "" {
		return false
	}
	dateInt64, err1 := strconv.ParseInt(good.ArticleDate, 10, 64)

	if err1 != nil {
		return true
	}
	arDate := time.Unix(dateInt64, 0)
	// fmt.Println("文章时间：" + arDate.Format(timeLayoutStr) + "昨天时间：" + beforeYesDate.Format(timeLayoutStr))
	if arDate.Before(beforeYesDate) {
		return true
	}

	return false
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

	// Ensure local keyword/regex filtering works even when caller only passed keyword.
	if len(rule.Words) == 0 {
		rule.Words = []string{keyword}
	}

	seeds := searchSeeds(keyword)
	if len(seeds) == 0 {
		seeds = []string{keyword}
	}
	// Cap API fan-out; longer seeds are already sorted first.
	if len(seeds) > 5 {
		seeds = seeds[:5]
	}

	results := make([]Product, 0, limit)
	seen := map[string]bool{}
	for _, seed := range seeds {
		if len(results) >= limit {
			break
		}
		for page := 0; page < 2 && len(results) < limit; page++ {
			rows := GetGoods(page, seed).Data.Rows
			if len(rows) == 0 {
				break
			}
			for _, item := range rows {
				if item.ArticleId != "" && seen[item.ArticleId] {
					continue
				}
				if !searchProductMatches(item, rule) {
					continue
				}
				if item.ArticleId != "" {
					seen[item.ArticleId] = true
				}
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
	// Re-check keyword/regex locally: SMZDM API only accepts plain search seeds.
	if len(rule.Words) > 0 && !containsAnyWord(good, rule.Words) {
		return false
	}
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
	word = strings.TrimSpace(word)
	if word == "" {
		return false
	}

	// Keywords/filters support regex (case-insensitive). Invalid patterns fall back to literal contains.
	if re, err := regexp.Compile("(?i)" + word); err == nil {
		return re.MatchString(good.ArticleTitle) ||
			re.MatchString(good.ArticlePrice) ||
			re.MatchString(good.Referral)
	}

	lw := strings.ToLower(word)
	return strings.Contains(strings.ToLower(good.ArticleTitle), lw) ||
		strings.Contains(strings.ToLower(good.ArticlePrice), lw) ||
		strings.Contains(strings.ToLower(good.Referral), lw)
}

// regexMetaPattern detects characters that change plain-text matching when used as regex.
var regexMetaPattern = regexp.MustCompile(`[.\\+*?()|\[\]{}^$]`)

// literalSeedPattern extracts a plain search seed for SMZDM API from a regex/keyword.
var literalSeedPattern = regexp.MustCompile(`[\p{L}\p{N}]+`)

// simpleGroupPattern matches one non-nested (a|b|c) group for seed expansion.
var simpleGroupPattern = regexp.MustCompile(`\(([^()]+)\)`)

func hasRegexMeta(s string) bool {
	return regexMetaPattern.MatchString(s)
}

// searchSeeds returns plain keywords suitable for SMZDM API search.
// Examples:
//
//	显示器            → [显示器]
//	显示器|屏幕       → [显示器, 屏幕]
//	(香|大)米         → [香米, 大米]
//	(4K|8K).{0,6}显示器 → [显示器]
func searchSeeds(keyword string) []string {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil
	}
	if !hasRegexMeta(keyword) {
		return []string{keyword}
	}

	seen := map[string]bool{}
	seeds := make([]string, 0, 4)
	add := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" || seen[s] {
			return
		}
		if hasRegexMeta(s) {
			s = extractLiteralSeed(s)
			if s == "" || seen[s] {
				return
			}
		}
		seen[s] = true
		seeds = append(seeds, s)
	}

	for _, part := range splitRegexAlternatives(keyword) {
		for _, expanded := range expandSimpleGroups(part) {
			add(expanded)
		}
	}

	// Prefer longer seeds first (better SMZDM recall than single characters).
	sort.SliceStable(seeds, func(i, j int) bool {
		return len([]rune(seeds[i])) > len([]rune(seeds[j]))
	})

	if len(seeds) == 0 {
		if seed := extractLiteralSeed(keyword); seed != "" {
			return []string{seed}
		}
		return []string{keyword}
	}
	return seeds
}

// expandSimpleGroups expands non-nested (a|b) groups into concrete strings.
// "(香|大)米" → ["香米", "大米"]. Nested/complex patterns are left as-is for literal extraction.
func expandSimpleGroups(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	loc := simpleGroupPattern.FindStringSubmatchIndex(s)
	if loc == nil {
		return []string{s}
	}

	before := s[:loc[0]]
	inside := s[loc[2]:loc[3]]
	after := s[loc[1]:]

	// Only expand pure alternations of literals (no nested groups / character classes).
	if strings.ContainsAny(inside, "()[]{}\\.+*?^$") {
		return []string{s}
	}

	alts := strings.Split(inside, "|")
	out := make([]string, 0, len(alts))
	for _, alt := range alts {
		alt = strings.TrimSpace(alt)
		if alt == "" {
			continue
		}
		combined := before + alt + after
		out = append(out, expandSimpleGroups(combined)...)
	}
	if len(out) == 0 {
		return []string{s}
	}
	return out
}

func extractLiteralSeed(s string) string {
	matches := literalSeedPattern.FindAllString(s, -1)
	if len(matches) == 0 {
		return ""
	}
	best := matches[0]
	bestLen := len([]rune(best))
	for _, m := range matches[1:] {
		n := len([]rune(m))
		if n > bestLen {
			best = m
			bestLen = n
		}
	}
	return best
}

// splitRegexAlternatives splits on top-level | (outside parentheses).
func splitRegexAlternatives(s string) []string {
	var parts []string
	depth := 0
	start := 0
	for i, r := range s {
		switch r {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case '|':
			if depth == 0 {
				parts = append(parts, strings.TrimSpace(s[start:i]))
				start = i + len("|")
			}
		}
	}
	parts = append(parts, strings.TrimSpace(s[start:]))
	return parts
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
