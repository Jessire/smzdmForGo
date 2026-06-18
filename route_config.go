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
	Words         []string `json:"words"`
	FilterWords   []string `json:"filterWords"`
	LowCommentNum int      `json:"lowCommentNum"`
	LowWorthyNum  int      `json:"lowWorthyNum"`
	MinPrice      float64  `json:"minPrice"`
	MaxPrice      float64  `json:"maxPrice"`
}

func productConfigFromConfig(conf file.Config) productConfigRequest {
	rules := make([]keywordRuleConfigRequest, 0, len(conf.KeywordRules))
	for _, rule := range conf.KeywordRules {
		item := keywordRuleConfigRequest{
			Words:       cleanWords(rule.Words),
			FilterWords: cleanWords(rule.FilterWords),
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
	conf.KeyWords = cleanWords(req.KeyWords)
	conf.FilterWords = cleanWords(req.FilterWords)
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

	rules := make([]file.KeywordRule, 0, len(req.KeywordRules))
	for _, ruleReq := range req.KeywordRules {
		words := cleanWords(ruleReq.Words)
		if len(words) == 0 {
			continue
		}
		rule := file.KeywordRule{
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
