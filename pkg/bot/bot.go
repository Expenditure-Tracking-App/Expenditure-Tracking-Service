package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"main/pkg/session"
	"main/pkg/storage"
	"main/pkg/transaction"
)

// Map to track ongoing sessions (active users)
var userSessions = make(map[int64]*session.UserSession)

// Bot represents the Telegram bot.
type Bot struct {
	api *tgbotapi.BotAPI
}

// NewBot creates a new bot instance.
func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}
	api.Debug = true
	log.Printf("Authorized on account %s", api.Self.UserName)
	return &Bot{api: api}, nil
}

// StartListening starts listening for updates.
func (b *Bot) StartListening(userSessions map[int64]*session.UserSession) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			err := b.handleTextMessage(update.Message, userSessions)
			if err != nil {
				log.Println(err)
			}
		} else if update.CallbackQuery != nil {
			err := b.handleCallbackQuery(update.CallbackQuery, userSessions)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// handleTextMessage handles incoming text messages.
func (b *Bot) handleTextMessage(message *tgbotapi.Message, userSessions map[int64]*session.UserSession) error {
	chatID := message.Chat.ID

	if message.Text == "/add" {
		return b.startSession(chatID, userSessions)
	}

	if _, exists := userSessions[chatID]; exists {
		return b.handleAnswer(chatID, userSessions, message.Text)
	}

	msg := tgbotapi.NewMessage(chatID, "Send /add to begin!")
	_, err := b.api.Send(msg)
	return err
}

// startSession starts a new session for the user.
func (b *Bot) startSession(chatID int64, userSessions map[int64]*session.UserSession) error {
	userSessions[chatID] = session.NewUserSession()
	return b.askCurrentQuestion(chatID, userSessions)
}

// askCurrentQuestion sends the current question to the user.
func (b *Bot) askCurrentQuestion(chatID int64, userSessions map[int64]*session.UserSession) error {
	userSession := userSessions[chatID]
	question := session.Questions[userSession.CurrentQuestion]

	msg := tgbotapi.NewMessage(chatID, question)

	if userSession.CurrentQuestion == session.QuestionName {
		var quickInputButtons []tgbotapi.InlineKeyboardButton
		for _, input := range session.QuickInput {
			quickInputButtons = append(quickInputButtons, tgbotapi.NewInlineKeyboardButtonData(input, input))
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(quickInputButtons...))
		msg.ReplyMarkup = keyboard
	}

	if userSession.CurrentQuestion == session.QuestionIsClaimable || userSession.CurrentQuestion == session.QuestionPaidForFamily {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Yes", "yes"),
				tgbotapi.NewInlineKeyboardButtonData("No", "no"),
			),
		)
		msg.ReplyMarkup = keyboard
	}

	if userSession.CurrentQuestion == session.QuestionCurrency {
		var currencyButtons []tgbotapi.InlineKeyboardButton
		for _, currency := range session.Currencies {
			currencyButtons = append(currencyButtons, tgbotapi.NewInlineKeyboardButtonData(currency, currency))
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(currencyButtons...))
		msg.ReplyMarkup = keyboard
	}

	if userSession.CurrentQuestion == session.QuestionCategory {
		var categoryButtons []tgbotapi.InlineKeyboardButton
		for _, category := range session.TransactionCategory {
			categoryButtons = append(categoryButtons, tgbotapi.NewInlineKeyboardButtonData(category, category))
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(categoryButtons...))
		msg.ReplyMarkup = keyboard
	}

	_, err := b.api.Send(msg)
	return err
}

// handleAnswer processes the user's answer.
func (b *Bot) handleAnswer(chatID int64, individualSession map[int64]*session.UserSession, answer string) error {
	err := individualSession[chatID].HandleAnswer(answer)
	if err != nil {
		return err
	}

	if individualSession[chatID].IsSessionComplete() {
		return b.completeSession(chatID, individualSession[chatID])
	}

	return b.askCurrentQuestion(chatID, individualSession)
}

// completeSession finishes the session.
func (b *Bot) completeSession(chatID int64, session *session.UserSession) error {
	if storage.UseDBToSave {
		// Save the responses to the database
		err := storage.SaveTransactionToDB(transaction.TransactionV2(session.Answers))
		if err != nil {
			// Inform the user if saving failed
			errMsg := tgbotapi.NewMessage(chatID, "Sorry, there was an error saving your transaction. Please try again later.")
			_, sendErr := b.api.Send(errMsg)
			if sendErr != nil {
				log.Printf("Error sending save error message: %v", sendErr)
			}
			// Also return the original save error
			return fmt.Errorf("failed to save transaction to DB: %w", err)
		}
	} else {
		// Save the responses to file
		storage.SaveResponseToFile(session.Answers)
	}

	// Send a thank-you message and confirmation
	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("Thank you for your responses!\n\nHere are your answers:\nName: %s\nAmount: %f\nCurrency: %s\nDate: %s\nIs Claimable: %t\nPaid for Family: %t\nCategory: %s",
			session.Answers.Name, session.Answers.Amount, session.Answers.Currency, session.Answers.Date, session.Answers.IsClaimable, session.Answers.PaidForFamily, session.Answers.Category))

	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Error sending confirmation message: %v", err)
		return err
	}

	_, err = b.api.Send(tgbotapi.NewMessage(chatID, "Send /add to begin!"))
	if err != nil {
		return err
	}

	delete(userSessions, chatID)
	return nil
}

// handleCallbackQuery handles callback queries.
func (b *Bot) handleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery, userSessions map[int64]*session.UserSession) error {
	chatID := callbackQuery.Message.Chat.ID
	userSession, exists := userSessions[chatID]
	if !exists {
		return fmt.Errorf("session not found for chat ID: %d", chatID)
	}

	switch callbackQuery.Data {
	case "yes":
		if userSession.CurrentQuestion == session.QuestionIsClaimable {
			userSession.Answers.IsClaimable = true
		} else if userSession.CurrentQuestion == session.QuestionPaidForFamily {
			userSession.Answers.PaidForFamily = true
		}
	case "no":
		if userSession.CurrentQuestion == session.QuestionIsClaimable {
			userSession.Answers.IsClaimable = false
		} else if userSession.CurrentQuestion == session.QuestionPaidForFamily {
			userSession.Answers.PaidForFamily = false
		}
	default:
		if userSession.CurrentQuestion == session.QuestionName {
			userSession.Answers.Name = callbackQuery.Data
		} else if userSession.CurrentQuestion == session.QuestionCurrency {
			userSession.Answers.Currency = callbackQuery.Data
		} else if userSession.CurrentQuestion == session.QuestionCategory {
			userSession.Answers.Category = callbackQuery.Data
		}
	}

	if userSession.CurrentQuestion == session.QuestionIsClaimable && session.DefaultPaidForFamily(userSession.Answers.Name) {
		userSession.Answers.PaidForFamily = session.DefaultPaidForFamily(userSession.Answers.Name)
		userSession.CurrentQuestion++
	}

	if userSession.CurrentQuestion == session.QuestionPaidForFamily && len(session.DefaultCategory(userSession.Answers.Name)) > 0 {
		userSession.Answers.Category = session.DefaultCategory(userSession.Answers.Name)
		userSession.CurrentQuestion++
	}

	userSession.CurrentQuestion++

	if userSession.IsSessionComplete() {
		err := b.completeSession(chatID, userSession)
		if err != nil {
			return err
		}
		return nil
	}

	err := b.askCurrentQuestion(chatID, userSessions)
	if err != nil {
		return err
	}

	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	if _, err := b.api.Request(callback); err != nil {
		return err
	}

	return nil
}
