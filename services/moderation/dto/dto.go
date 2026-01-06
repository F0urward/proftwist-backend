package dto

type ModerationResult struct {
	Allowed    bool     `json:"allowed"`
	Categories []string `json:"categories"`
}

type ModerateContentRequest struct {
	Content string `json:"content"`
}

type ModerateContentResponse struct {
	Result ModerationResult `json:"result"`
}
