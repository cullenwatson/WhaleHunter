package core

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type TradingViewClient struct {
	Symbol    string
	Timeframe string
	Candles   int
	AuthToken string
	Ws        *websocket.Conn
}

func NewTradingViewClient(symbol, timeframe string, candles int, authToken string) *TradingViewClient {
	return &TradingViewClient{
		Symbol:    symbol,
		Timeframe: timeframe,
		Candles:   candles,
		AuthToken: authToken,
	}
}

func (tv *TradingViewClient) Connect() error {
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	headers := http.Header{}
	headers.Set("Origin", "https://www.tradingview.com")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Accept-Language", "en-US,en;q=0.9")
	headers.Set("Pragma", "no-cache")
	headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	ws, _, err := dialer.Dial("wss://prodata.tradingview.com/socket.io/websocket", headers)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}

	tv.Ws = ws
	return nil
}

// generateRandomString helper for session generation
func generateRandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// generateChartSession returns something like "cs_abcdef..."
func generateChartSession() string {
	return "cs_" + generateRandomString(12)
}

// createMessage replicates "~m~<len>~m~<json>" format
func createMessage(functionName string, paramList interface{}) string {
	payload := map[string]interface{}{
		"m": functionName,
		"p": paramList,
	}
	jsonBytes, _ := json.Marshal(payload)
	jsonStr := string(jsonBytes)
	return fmt.Sprintf("~m~%d~m~%s", len(jsonStr), jsonStr)
}

func (tv *TradingViewClient) SendMessage(functionName string, params interface{}) error {
	if tv.Ws == nil {
		return fmt.Errorf("websocket is not connected")
	}

	msg := createMessage(functionName, params)
	err := tv.Ws.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write websocket message: %w", err)
	}
	return nil
}

func (tv *TradingViewClient) Run() error {
	// 1. Generate a chart session
	chartSession := generateChartSession()

	// 2. Build a JSON-like string referencing the symbol
	symbolString := fmt.Sprintf("={\"adjustment\":\"splits\",\"currency-id\":\"USD\",\"session\":\"regular\",\"symbol\":\"%s\"}", tv.Symbol)

	// 3. Auth token
	if err := tv.SendMessage("set_auth_token", []string{tv.AuthToken}); err != nil {
		return err
	}

	// 4. Create a chart session
	if err := tv.SendMessage("chart_create_session", []string{chartSession, ""}); err != nil {
		return err
	}

	// 5. Resolve the symbol
	symbolLabel := "sds_sym_2"
	if err := tv.SendMessage("resolve_symbol", []interface{}{chartSession, symbolLabel, symbolString}); err != nil {
		return err
	}

	// 6. Attach symbol to the series
	seriesLabel := "sds_1"
	seriesVersion := "s1"
	if err := tv.SendMessage("create_series", []interface{}{
		chartSession,
		seriesLabel,
		seriesVersion,
		symbolLabel,
		tv.Timeframe,
		tv.Candles,
		"",
	}); err != nil {
		return err
	}

	return nil
}
