package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"tw-scanner/models"
)

const tokenFile = "auth_token.txt"

func getTokenPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return tokenFile
	}
	return filepath.Join(homeDir, ".tradingview", tokenFile)
}

func saveAuthToken(token string) error {
	tokenPath := getTokenPath()
	dir := filepath.Dir(tokenPath)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	return os.WriteFile(tokenPath, []byte(token), 0600)
}

func loadAuthToken() (string, error) {
	tokenPath := getTokenPath()
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func SignIn(creds models.Credentials) (string, error) {
	if token, err := loadAuthToken(); err == nil {
		return token, nil
	}

	token, err := signInWithRetry(creds, "", false)
	if err != nil {
		return "", err
	}

	if err := saveAuthToken(token); err != nil {
		log.Warn().Err(err).Msg("Failed to save auth token")
	}

	return token, nil
}

func signInWithRetry(creds models.Credentials, captchaResponse string, isRetry bool) (string, error) {
	url := "https://www.tradingview.com/accounts/signin/"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	_ = w.WriteField("username", creds.Username)
	_ = w.WriteField("password", creds.Password)
	_ = w.WriteField("remember", "true")
	if captchaResponse != "" {
		_ = w.WriteField("g-recaptcha-response-v2", captchaResponse)
	}
	err := w.Close()
	if err != nil {
		return "", fmt.Errorf("error closing multipart writer: %v", err)
	}

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

	var loginResp models.LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	if loginResp.Error != "" {
		log.Error().
			Str("error", loginResp.Error).
			Msg("TradingView login error occurred")

		if strings.Contains(strings.ToLower(loginResp.Error), "captcha") && !isRetry {
			captchaResponse, err := solveCaptcha()
			if err != nil {
				return "", fmt.Errorf("failed to solve captcha: %v", err)
			}
			return signInWithRetry(creds, captchaResponse, true)
		}
		return "", fmt.Errorf("login failed: %s", loginResp.Error)
	}

	return loginResp.User.AuthToken, nil
}
