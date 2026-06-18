package main

import (
	"strings"

	"ggball.com/smzdm/file"
)

type productConfigRequest struct {
	KeyWords      []string                   `json:"keyWords"`
	FilterWords   []string                   `json:"filterWords"`
	LowCommentNum int                        `json:"lowCommentNum"`
	LowWorthyNum  int                        `json:"lowWorthyNum"`
	MinPrice      float64                    `json:"minPrice"`
	MaxPrice      float64                    `json:"maxPrice"`
	SatisfyNum    int                        `json:"satisfyNum"`
	TickTime      int                        `json:"tickTime"`
	Cron          string                     `json:"cron"`
	Telegram      telegramConfigRequest      `json:"telegram"`
	KeywordRules  []keywordRuleConfigRequest `json:"keywordRules"`
}

type telegramConfigRequest struct {
	Enabled               bool   `json:"enabled"`
	BotToken              string `json:"botToken"`
	ChatID                string `json:"chatId"`
	ParseMode             string `json:"parseMode"`
	DisableWebPagePreview bool   `json:"disableWebPagePreview"`
}

type keywordRuleConfigRequest struct {
	Enabled       *bool    `json:"enabled,omitempty"`
	Words         []string `json:"words"`
	FilterWords   []string `json:"filterWords"`
	LowCommentNum int      `json:"lowCommentNum"`
	LowWorthyNum  int      `json:"lowWorthyNum"`
	MinPrice      float64  `json:"minPrice"`
	MaxPrice      float64  `json:"maxPrice"`
}

func productConfigFromConfig(conf file.Config) productConfigRequest {
	rules := make([]keywordRuleConfigRequest, 0, len(conf.KeywordRules)+len(conf.KeyWords))
	seenWords := map[string]bool{}
	for _, rule := range conf.KeywordRules {
		enabled := keywordRuleEnabled(rule)
		item := keywordRuleConfigRequest{
			Enabled:     &enabled,
			Words:       cleanWords(rule.Words),
			FilterWords: cleanWords(rule.FilterWords),
		}
		for _, word := range item.Words {
			seenWords[word] = true
		}
		if rule.LowCommentNum != nil {
			item.LowCommentNum = *rule.LowCommentNum
		}
		if rule.LowWorthyNum != nil {
			item.LowWorthyNum = *rule.LowWorthyNum
		}
		if rule.MinPrice != nil {
			item.MinPrice = *rule.MinPrice
		}
		if rule.MaxPrice != nil {
			item.MaxPrice = *rule.MaxPrice
		}
		rules = append(rules, item)
	}
	for _, word := range cleanWords(conf.KeyWords) {
		if seenWords[word] {
			continue
		}
		enabled := true
		rules = append(rules, keywordRuleConfigRequest{
			Enabled:       &enabled,
			Words:         []string{word},
			FilterWords:   cleanWords(conf.FilterWords),
			LowCommentNum: conf.LowCommentNum,
			LowWorthyNum:  conf.LowWorthyNum,
			MinPrice:      conf.MinPrice,
			MaxPrice:      conf.MaxPrice,
		})
	}
	return productConfigRequest{
		KeyWords:      cleanWords(conf.KeyWords),
		FilterWords:   cleanWords(conf.FilterWords),
		LowCommentNum: conf.LowCommentNum,
		LowWorthyNum:  conf.LowWorthyNum,
		MinPrice:      conf.MinPrice,
		MaxPrice:      conf.MaxPrice,
		SatisfyNum:    conf.SatisfyNum,
		TickTime:      conf.TickTime,
		Cron:          conf.Cron,
		Telegram: telegramConfigRequest{
			Enabled:               conf.Telegram.Enabled,
			BotToken:              conf.Telegram.BotToken,
			ChatID:                conf.Telegram.ChatID,
			ParseMode:             normalizedParseMode(conf.Telegram.ParseMode),
			DisableWebPagePreview: conf.Telegram.DisableWebPagePreview,
		},
		KeywordRules: rules,
	}
}

func (req productConfigRequest) applyTo(conf file.Config) file.Config {
	conf.FilterWords = []string{}
	conf.LowCommentNum = nonNegativeInt(req.LowCommentNum)
	conf.LowWorthyNum = nonNegativeInt(req.LowWorthyNum)
	conf.MinPrice = nonNegativeFloat(req.MinPrice)
	conf.MaxPrice = nonNegativeFloat(req.MaxPrice)
	conf.SatisfyNum = req.SatisfyNum
	if conf.SatisfyNum <= 0 {
		conf.SatisfyNum = 5
	}
	conf.TickTime = req.TickTime
	if conf.TickTime <= 0 {
		conf.TickTime = 10800
	}
	conf.Cron = strings.TrimSpace(req.Cron)
	if conf.Cron == "" {
		conf.Cron = "0 10 10 ? * *"
	}
	conf.Telegram = file.Telegram{
		Enabled:               req.Telegram.Enabled,
		BotToken:              strings.TrimSpace(req.Telegram.BotToken),
		ChatID:                strings.TrimSpace(req.Telegram.ChatID),
		ParseMode:             normalizedParseMode(req.Telegram.ParseMode),
		DisableWebPagePreview: req.Telegram.DisableWebPagePreview,
	}

	ruleRequests := req.KeywordRules
	if len(ruleRequests) == 0 {
		for _, word := range cleanWords(req.KeyWords) {
			enabled := true
			ruleRequests = append(ruleRequests, keywordRuleConfigRequest{
				Enabled:       &enabled,
				Words:         []string{word},
				FilterWords:   cleanWords(req.FilterWords),
				LowCommentNum: req.LowCommentNum,
				LowWorthyNum:  req.LowWorthyNum,
				MinPrice:      req.MinPrice,
				MaxPrice:      req.MaxPrice,
			})
		}
	}

	rules := make([]file.KeywordRule, 0, len(ruleRequests))
	for _, ruleReq := range ruleRequests {
		words := cleanWords(ruleReq.Words)
		if len(words) == 0 {
			continue
		}
		enabled := true
		if ruleReq.Enabled != nil {
			enabled = *ruleReq.Enabled
		}
		rule := file.KeywordRule{
			Enabled:     &enabled,
			Words:       words,
			FilterWords: cleanWords(ruleReq.FilterWords),
		}
		if ruleReq.LowCommentNum > 0 {
			value := ruleReq.LowCommentNum
			rule.LowCommentNum = &value
		}
		if ruleReq.LowWorthyNum > 0 {
			value := ruleReq.LowWorthyNum
			rule.LowWorthyNum = &value
		}
		if ruleReq.MinPrice > 0 {
			value := ruleReq.MinPrice
			rule.MinPrice = &value
		}
		if ruleReq.MaxPrice > 0 {
			value := ruleReq.MaxPrice
			rule.MaxPrice = &value
		}
		rules = append(rules, rule)
	}
	conf.KeywordRules = rules
	conf.KeyWords = flattenRuleWords(rules)
	if len(conf.KeywordRules) > 0 {
		first := conf.KeywordRules[0]
		conf.FilterWords = cleanWords(first.FilterWords)
		if first.LowCommentNum != nil {
			conf.LowCommentNum = *first.LowCommentNum
		}
		if first.LowWorthyNum != nil {
			conf.LowWorthyNum = *first.LowWorthyNum
		}
		if first.MinPrice != nil {
			conf.MinPrice = *first.MinPrice
		}
		if first.MaxPrice != nil {
			conf.MaxPrice = *first.MaxPrice
		}
	}
	return conf
}

func cleanWords(words []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(words))
	for _, word := range words {
		for _, part := range strings.FieldsFunc(word, splitWordList) {
			part = strings.TrimSpace(part)
			if part == "" || seen[part] {
				continue
			}
			seen[part] = true
			result = append(result, part)
		}
	}
	return result
}

func splitWordList(r rune) bool {
	return r == ',' || r == '，' || r == '\n' || r == '\r' || r == '\t'
}

func nonNegativeInt(value int) int {
	if value < 0 {
		return 0
	}
	return value
}

func nonNegativeFloat(value float64) float64 {
	if value < 0 {
		return 0
	}
	return value
}

func normalizedParseMode(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "HTML"
	}
	return value
}

func keywordRuleEnabled(rule file.KeywordRule) bool {
	return rule.Enabled == nil || *rule.Enabled
}

func flattenRuleWords(rules []file.KeywordRule) []string {
	words := []string{}
	for _, rule := range rules {
		words = append(words, rule.Words...)
	}
	return cleanWords(words)
}
