package tradingview

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

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

// SignIn attempts to login to TradingView and returns the auth token
func SignIn(creds Credentials) (string, error) {
	url := "https://www.tradingview.com/accounts/signin/"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	_ = w.WriteField("username", creds.Username)
	_ = w.WriteField("password", creds.Password)
	_ = w.WriteField("remember", "true")
	w.Close()

	// Create request
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.tradingview.com")
	req.Header.Set("Referer", "https://www.tradingview.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("X-Language", "en")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	if loginResp.Error != "" {
		return "", fmt.Errorf("login failed: %s", loginResp.Error)
	}

	return loginResp.User.AuthToken, nil
}
