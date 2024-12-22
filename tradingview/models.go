package tradingview

type Credentials struct {
	Username string
	Password string
}

type LoginResponse struct {
	Error string `json:"error"`
	User  struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		AuthToken string `json:"auth_token"`
	} `json:"user"`
}
