package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	logger *Logger
)

const botTokenEnv = "TG_TOKEN_LOGGER"

type Logger struct {
	tgClient *TelegramClient
}

func init() {
	token, ex := os.LookupEnv(botTokenEnv)
	if !ex {
		log.Fatalf("telegram token is not set in env")
	}

	logger = register(token)
}

func register(botToken string) *Logger {
	tgClient := newTelegramClient(botToken)

	return &Logger{tgClient: tgClient}
}

func entry() *Logger {
	return logger
}

func Errorf(format string, args ...interface{}) {
	entry().Errorf(format, args)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	errWithStack := fmt.Sprintf("%+v", fmt.Errorf(format, args...))
	l.tgClient.SendLog(errWithStack)
}
