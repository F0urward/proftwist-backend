package entities

type ModerationResult struct {
	Allowed    bool     `json:"allowed"`
	Categories []string `json:"categories"`
}
