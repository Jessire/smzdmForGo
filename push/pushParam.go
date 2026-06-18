package push

type TelegramMessageParam struct {
	ChatID                string `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview,omitempty"`
}

type TelegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
}
