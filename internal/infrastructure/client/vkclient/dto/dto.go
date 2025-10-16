package dto

type VKTokenExchangeResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	UserID       int64  `json:"user_id"`
	State        string `json:"state,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

type VKUserInfoResponse struct {
	User struct {
		UserID    string `json:"user_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Avatar    string `json:"avatar,omitempty"`
		Email     string `json:"email"`
	} `json:"user"`
}
