package tglogger

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	logger *Logger
)

const (
	botTokenEnv  = "TG_LOGGER_TOKEN"
	botChatIDEnv = "TG_LOGGER_CHAT_ID"
)

type Logger struct {
	tgClient *TelegramClient
	fields   Fields
}

func init() {
	token, ex := os.LookupEnv(botTokenEnv)
	if !ex {
		log.Fatalf("telegram token is not set in env")
	}

	chatID, ex := os.LookupEnv(botChatIDEnv)
	if !ex {
		log.Fatalf("telegram chat id is not set in env")
	}

	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		log.Fatalf("telegram chat id shoud be int")
	}

	logger = register(token, chatIDInt)
}

func register(botToken string, chatID int64) *Logger {
	tgClient := newTelegramClient(botToken, chatID)

	return &Logger{tgClient: tgClient}
}

func entry() *Logger {
	return logger
}

func Errorf(format string, args ...interface{}) {
	entry().Error(levelError, WithStackErr(fmt.Errorf(format, args...)))
}

func Warnf(format string, args ...interface{}) {
	entry().Error(levelWarn, WithStackErr(fmt.Errorf(format, args...)))
}

func WithFields(fields Fields) *Logger {
	nl := &Logger{tgClient: entry().tgClient}
	nl.fields = make(Fields, len(fields))

	for k, v := range fields {
		nl.fields[k] = v
	}

	return nl
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(levelError, WithStackErr(fmt.Errorf(format, args...)))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Error(levelWarn, WithStackErr(fmt.Errorf(format, args...)))
}

func (l *Logger) Error(level string, err error) {
	l.sendError(level, err)
}

type Fields map[string]interface{}

func (l *Logger) WithFields(fields Fields) *Logger {
	nl := &Logger{tgClient: l.tgClient}
	nl.fields = make(Fields, len(fields))

	for k, v := range fields {
		nl.fields[k] = v
	}

	return nl
}

func (l *Logger) sendError(level string, err error) {
	d := time.Now()
	message := fmt.Sprintf(
		"<b>%s</b>\n"+
			"[  %s ] \n\n"+
			"%+v", level, d.Format(time.DateTime), err)

	if len(l.fields) > 0 {
		message += fmt.Sprintf(
			"\n\nExtra params:\n\n" + l.prettyPrintMap(l.fields))
	}

	l.tgClient.SendLog(message)
}

func (l *Logger) prettyPrintMap(m map[string]interface{}) string {
	res := "<code>"
	for k, v := range m {
		res += fmt.Sprint(k, ": ", v, "\n")
	}
	res += "</code>"

	return res
}
