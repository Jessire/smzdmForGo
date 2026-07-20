package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"ggball.com/smzdm/file"
	"ggball.com/smzdm/smzdm"
)

type productSearchRequest struct {
	Rule  keywordRuleConfigRequest `json:"rule"`
	Limit int                      `json:"limit"`
}

type discoverySearchRequest struct {
	Type      string                   `json:"type"`
	GlobalHot globalHotConfigRequest   `json:"globalHot"`
	Rule      keywordRuleConfigRequest `json:"rule"`
	Limit     int                      `json:"limit"`
}

type productSearchResponse struct {
	Keyword     string                 `json:"keyword"`
	OpenURL     string                 `json:"openUrl"`
	Items       []productSearchProduct `json:"items"`
	Pushable    int                    `json:"pushable"`
	AlreadyPush int                    `json:"alreadyPushed"`
	TooOld      int                    `json:"tooOld"`
}

type productSearchProduct struct {
	Title     string `json:"title"`
	Price     string `json:"price"`
	Worthy    string `json:"worthy"`
	Comment   string `json:"comment"`
	Pic       string `json:"pic"`
	URL       string `json:"url"`
	Referral  string `json:"referral"`
	Date      string `json:"date"`
	PushState string `json:"pushState"` // ok | already_pushed | too_old | bad_date
}

func ProductSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var req productSearchRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, fmt.Errorf("解析搜索条件失败: %v", err))
		return
	}

	words := cleanWords(req.Rule.Words)
	if len(words) == 0 {
		writeError(w, fmt.Errorf("关键词不能为空"))
		return
	}
	keyword := words[0]
	rule := searchRuleFromRequest(req.Rule)
	products := smzdm.SearchGoods(keyword, rule, req.Limit)
	response := productSearchPayload(keyword, keywordSearchURL(keyword), products)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code": "0",
		"msg":  "",
		"data": response,
	})
}

func DiscoverySearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var req discoverySearchRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, fmt.Errorf("解析搜索条件失败: %v", err))
		return
	}

	kind := strings.ToLower(strings.TrimSpace(req.Type))
	if kind != "hot" && kind != "author" {
		writeError(w, fmt.Errorf("不支持的发现规则类型"))
		return
	}

	globalHot := file.GlobalHotConfig{
		Enabled:              req.GlobalHot.Enabled,
		WindowHours:          normalizeGlobalHotWindow(req.GlobalHot.WindowHours),
		MinCommentNum:        normalizeGlobalHotCommentFloor(req.GlobalHot.MinCommentNum),
		HotKeywords:          cleanWords(req.GlobalHot.HotKeywords),
		FollowAuthorsEnabled: req.GlobalHot.FollowAuthorsEnabled,
		FollowedAuthors:      cleanWords(req.GlobalHot.FollowedAuthors),
		AuthorKeywords:       cleanWords(req.GlobalHot.AuthorKeywords),
		FilterWords:          cleanWords(req.GlobalHot.FilterWords),
		LowCommentNum:        nonNegativeInt(req.GlobalHot.LowCommentNum),
		LowWorthyNum:         nonNegativeInt(req.GlobalHot.LowWorthyNum),
		MinPrice:             nonNegativeFloat(req.GlobalHot.MinPrice),
		MaxPrice:             nonNegativeFloat(req.GlobalHot.MaxPrice),
	}
	common := searchRuleFromRequest(req.Rule)
	if len(common.FilterWords) == 0 && common.LowCommentNum == nil && common.LowWorthyNum == nil && common.MinPrice == nil && common.MaxPrice == nil {
		common = discoveryCommonRuleFromRequest(globalHot)
	}
	products := smzdm.SearchDiscoveryGoods(globalHot, kind, common, req.Limit)
	label := "全站热门"
	if kind == "author" {
		label = "关注作者"
	}
	response := productSearchPayload(label, "", products)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code": "0",
		"msg":  "",
		"data": response,
	})
}

func discoveryCommonRuleFromRequest(globalHot file.GlobalHotConfig) file.KeywordRule {
	rule := file.KeywordRule{FilterWords: cleanWords(globalHot.FilterWords)}
	if globalHot.LowCommentNum > 0 {
		value := globalHot.LowCommentNum
		rule.LowCommentNum = &value
	}
	if globalHot.LowWorthyNum > 0 {
		value := globalHot.LowWorthyNum
		rule.LowWorthyNum = &value
	}
	if globalHot.MinPrice > 0 {
		value := globalHot.MinPrice
		rule.MinPrice = &value
	}
	if globalHot.MaxPrice > 0 {
		value := globalHot.MaxPrice
		rule.MaxPrice = &value
	}
	return rule
}

func productSearchPayload(keyword, openURL string, products []smzdm.Product) productSearchResponse {
	pushedMap := smzdm.LoadPushedMap()
	maxAgeDays := smzdm.NormalizeMaxArticleAgeDays(currentConfig().MaxArticleAgeDays)
	items := make([]productSearchProduct, 0, len(products))
	pushable, already, tooOld := 0, 0, 0
	for _, product := range products {
		state := smzdm.PushSkipReason(product, pushedMap, maxAgeDays)
		if state == "" {
			state = "ok"
			pushable++
		} else if state == "already_pushed" {
			already++
		} else if state == "too_old" || state == "bad_date" {
			tooOld++
		}
		items = append(items, productSearchProduct{
			Title:     product.ArticleTitle,
			Price:     product.ArticlePrice,
			Worthy:    product.ArticleWorthy,
			Comment:   product.ArticleComment,
			Pic:       product.ArticlePic,
			URL:       product.ArticleUrl,
			Referral:  product.Referral,
			Date:      product.ArticleDate,
			PushState: state,
		})
	}

	return productSearchResponse{
		Keyword:     keyword,
		OpenURL:     openURL,
		Items:       items,
		Pushable:    pushable,
		AlreadyPush: already,
		TooOld:      tooOld,
	}
}

func searchRuleFromRequest(req keywordRuleConfigRequest) file.KeywordRule {
	rule := file.KeywordRule{
		Words:       cleanWords(req.Words),
		FilterWords: cleanWords(req.FilterWords),
	}
	if req.LowCommentNum > 0 {
		value := req.LowCommentNum
		rule.LowCommentNum = &value
	}
	if req.LowWorthyNum > 0 {
		value := req.LowWorthyNum
		rule.LowWorthyNum = &value
	}
	if req.MinPrice > 0 {
		value := req.MinPrice
		rule.MinPrice = &value
	}
	if req.MaxPrice > 0 {
		value := req.MaxPrice
		rule.MaxPrice = &value
	}
	return rule
}

func keywordSearchURL(keyword string) string {
	return "https://search.smzdm.com/?c=faxian&s=" + url.QueryEscape(keyword)
}
