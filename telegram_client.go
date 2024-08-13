package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func newTelegramClient(botToken string) *TelegramClient {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("failed to connect to telegram bot: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	var chatID int64
	for update := range bot.GetUpdatesChan(u) {
		if update.Message != nil && update.Message.Chat != nil {
			chatID = update.Message.Chat.ID
			break
		}
	}

	return &TelegramClient{
		chatID: chatID,
		client: cleanhttp.DefaultPooledClient(),
	}
}

type Message struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
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
		ChatID: t.chatID,
		Text:   text,
	})
	if err != nil {
		log.Printf("failed to marshal log message")

		return
	}

	method := fmt.Sprintf(tgApiSendMessage, t.token)

	resp, err := t.client.Post(method, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Printf("failed to send log")

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("failed to send log - resp is not ok but %d", resp.StatusCode)
	}
}
