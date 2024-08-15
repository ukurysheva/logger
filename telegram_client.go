package tglogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"log"
	"net/http"
)

const (
	tgApiSendMessage = "https://api.telegram.org/bot%s/sendMessage"

	levelWarn  = "WARNING"
	levelError = "ERROR"
)

type TelegramClient struct {
	token  string
	chatID int64
	client *http.Client
}

func newTelegramClient(botToken string, chatID int64) *TelegramClient {
	return &TelegramClient{
		token:  botToken,
		chatID: chatID,
		client: cleanhttp.DefaultPooledClient(),
	}
}

type Message struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// APIResponse is a response from the Telegram API with the result
// stored raw.
type APIResponse struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int             `json:"error_code"`
	Description string          `json:"description"`
}

func (t *TelegramClient) SendLog(text string) {
	payload, err := json.Marshal(Message{
		ChatID:    t.chatID,
		Text:      text,
		ParseMode: "HTML",
	})
	if err != nil {
		log.Printf("logger - failed to marshal log messag\ne")

		return
	}

	method := fmt.Sprintf(tgApiSendMessage, t.token)

	resp, err := t.client.Post(method, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Printf("logger - failed to send log\n\n")

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("logger - failed to send log - resp is not ok but %d\n", resp.StatusCode)
	}
}
