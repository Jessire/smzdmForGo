package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"ggball.com/smzdm/file"
	"ggball.com/smzdm/smzdm"
)

type productSearchRequest struct {
	Rule  keywordRuleConfigRequest `json:"rule"`
	Limit int                      `json:"limit"`
}

type productSearchResponse struct {
	Keyword string                 `json:"keyword"`
	OpenURL string                 `json:"openUrl"`
	Items   []productSearchProduct `json:"items"`
}

type productSearchProduct struct {
	Title    string `json:"title"`
	Price    string `json:"price"`
	Worthy   string `json:"worthy"`
	Comment  string `json:"comment"`
	Pic      string `json:"pic"`
	URL      string `json:"url"`
	Referral string `json:"referral"`
	Date     string `json:"date"`
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
	items := make([]productSearchProduct, 0, len(products))
	for _, product := range products {
		items = append(items, productSearchProduct{
			Title:    product.ArticleTitle,
			Price:    product.ArticlePrice,
			Worthy:   product.ArticleWorthy,
			Comment:  product.ArticleComment,
			Pic:      product.ArticlePic,
			URL:      product.ArticleUrl,
			Referral: product.Referral,
			Date:     product.ArticleDate,
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code": "0",
		"msg":  "",
		"data": productSearchResponse{
			Keyword: keyword,
			OpenURL: keywordSearchURL(keyword),
			Items:   items,
		},
	})
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
