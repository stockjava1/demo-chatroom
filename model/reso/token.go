package reso

// GetTokenInfo GET "/token/info" response object
type GetTokenInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
