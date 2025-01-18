package core

import (
	"crypto/tls"
	"fmt"
	"github.com/cullenwatson/WhaleHunter/indicator"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
)

type TradingViewClient struct {
	Symbol     string
	Timeframe  string
	Candles    int
	AuthToken  string
	Ws         *websocket.Conn
	Indicators []string
}

func NewTradingViewClient(symbol, timeframe string, candles int, authToken string, indicators []string) *TradingViewClient {
	return &TradingViewClient{
		Symbol:     symbol,
		Timeframe:  timeframe,
		Candles:    candles,
		AuthToken:  authToken,
		Indicators: indicators,
	}
}

func (tv *TradingViewClient) Connect() error {
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	headers := http.Header{}
	headers.Set("Origin", "https://www.tradingview.com")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Accept-Language", "en-US,en;q=0.9")
	headers.Set("Pragma", "no-cache")
	headers.Set("User-Agent", "Mozilla/5.0 ...")

	ws, _, err := dialer.Dial("wss://prodata.tradingview.com/socket.io/websocket?type=chart", headers)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}

	tv.Ws = ws
	return nil
}
func (tv *TradingViewClient) SendMessage(functionName string, params interface{}) error {
	if tv.Ws == nil {
		return fmt.Errorf("websocket is not connected")
	}
	msg := createMessage(functionName, params)
	return tv.Ws.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (tv *TradingViewClient) Run() error {
	chartSession := generateChartSession()
	log.Info().
		Str("symbol", tv.Symbol).
		Str("timeframe", tv.Timeframe).
		Int("candles", tv.Candles).
		Str("chart_session", chartSession).
		Strs("indicators", tv.Indicators).
		Msg("Starting TradingView client session")

	symbolString := fmt.Sprintf("={\"adjustment\":\"splits\",\"currency-id\":\"USD\",\"session\":\"regular\",\"symbol\":\"%s\"}", tv.Symbol)

	if err := tv.SendMessage("set_auth_token", []string{tv.AuthToken}); err != nil {
		return err
	}
	if err := tv.SendMessage("chart_create_session", []string{chartSession, ""}); err != nil {
		return err
	}
	symbolLabel := "sds_sym_1"
	if err := tv.SendMessage("resolve_symbol", []interface{}{chartSession, symbolLabel, symbolString}); err != nil {
		return err
	}

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

	// Loop over each indicator, adding each study
	for i, indName := range tv.Indicators {
		stLabel := fmt.Sprintf("st%d", i+1) // e.g. "st1","st2", etc.
		studyScript, err := indicator.GetIndicatorScript(indName)
		if err != nil {
			return err
		}
		// create_study
		if err := tv.SendMessage("create_study", []interface{}{
			chartSession,
			stLabel,
			"st1",
			seriesLabel,
			"Script@tv-scripting-101!",
			studyScript,
		}); err != nil {
			return err
		}
	}

	return nil
}
