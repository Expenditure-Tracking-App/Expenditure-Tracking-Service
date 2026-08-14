package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"main/pkg/session"
)

// TelegramSender abstracts Telegram message sending for mocking
type TelegramSender interface {
	Send(msg tgbotapi.MessageConfig) (tgbotapi.Message, error)
	Request(chattable tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
}

// SessionStore manages user sessions
type SessionStore interface {
	Get(chatID int64) (*session.UserSession, bool)
	Set(chatID int64, s *session.UserSession)
	Delete(chatID int64)
}
