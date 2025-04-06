package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
)

// UserSession - Session to track user's Q&A flow
type UserSession struct {
	CurrentQuestion int         // Index of the current question
	Answers         Transaction // Struct to store responses
}

// Questions array for the process
var questions = []string{
	"What is the name of the transaction?",
	"How much is the transaction?",
	"What currency is the transaction in?",
	"What is the date of transaction?",
	"Is it claimable?",
	"Is it payable for the family?",
}

// Map to track ongoing sessions (active users)
var userSessions = make(map[int64]*UserSession)

// SaveFilePath File to save responses
const SaveFilePath = "responses.txt"

// Available currencies
var currencies = []string{"USD", "EUR", "JPY", "SGD", "MYR"}

func main() {
	token, err := os.ReadFile("token.txt")
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(string(token))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Listen for incoming updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			err = processMessage(bot, update.Message)
			if err != nil {
				log.Println(err)
			}
		} else if update.CallbackQuery != nil {
			// Handle callback queries (button presses)
			err = processCallbackQuery(bot, update.CallbackQuery)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// Handle incoming messages
func processMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	chatID := message.Chat.ID

	// Handle start command to begin the Q&A process
	if message.Text == "/add" {
		err := startSession(bot, chatID)
		if err != nil {
			return err
		}
		return nil
	}

	// If the user is in the middle of a session, process their answer
	if session, exists := userSessions[chatID]; exists {
		err := processAnswer(bot, chatID, session, message.Text)
		if err != nil {
			return err
		}

		return nil
	}

	// If no session is active, guide the user
	msg := tgbotapi.NewMessage(chatID, "Send /add to begin!")
	_, err := bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// Start a new Q&A session
func startSession(bot *tgbotapi.BotAPI, chatID int64) error {
	// Create a session for the user
	userSessions[chatID] = &UserSession{
		CurrentQuestion: 0, // Start at the first question
	}

	// Ask the first question
	err := askCurrentQuestion(bot, chatID)
	if err != nil {
		return err
	}

	return nil
}

// Ask the current question in the session
func askCurrentQuestion(bot *tgbotapi.BotAPI, chatID int64) error {
	session := userSessions[chatID]
	question := questions[session.CurrentQuestion]

	// Send the current question to the user
	msg := tgbotapi.NewMessage(chatID, question)

	if session.CurrentQuestion == 4 || session.CurrentQuestion == 5 {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Yes", "yes"),
				tgbotapi.NewInlineKeyboardButtonData("No", "no"),
			),
		)
		msg.ReplyMarkup = keyboard
	}

	if session.CurrentQuestion == 2 {
		// Create currency buttons
		var currencyButtons []tgbotapi.InlineKeyboardButton
		for _, currency := range currencies {
			currencyButtons = append(currencyButtons, tgbotapi.NewInlineKeyboardButtonData(currency, currency))
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(currencyButtons...))
		msg.ReplyMarkup = keyboard
	}

	_, err := bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// Process the user's answer and move to the next question
func processAnswer(bot *tgbotapi.BotAPI, chatID int64, session *UserSession, answer string) error {
	var err error
	switch session.CurrentQuestion {
	case 0: // First question: name of transaction
		session.Answers.Name = answer
	case 1: // Second question: value of transaction (validate as float)
		session.Answers.Amount, err = parseFloat32(answer)
		if err != nil {
			return fmt.Errorf("invalid amount: %w", err)
		}
	case 2: // Third question: currency of transaction
		session.Answers.Currency = answer
	case 3: // Fourth question: date of transaction
		session.Answers.Date = answer
	case 4: // Fifth question: is it Claimable
		session.Answers.IsClaimable, err = parseBool(answer)
		if err != nil {
			return fmt.Errorf("invalid claimable value: %w", err)
		}
	case 5:
		session.Answers.PaidForFamily, err = parseBool(answer)
		if err != nil {
			return fmt.Errorf("invalid paid for family value: %w", err)
		}
	}

	// Move to the next question
	session.CurrentQuestion++

	// If all questions are answered, finish the session
	if session.CurrentQuestion >= len(questions) {
		err := finishSession(bot, chatID, session)
		if err != nil {
			return err
		}

		return nil
	}

	// Otherwise, ask the next question
	err = askCurrentQuestion(bot, chatID)
	if err != nil {
		return err
	}

	return nil
}

// Finish the Q&A session
func finishSession(bot *tgbotapi.BotAPI, chatID int64, session *UserSession) error {
	// Save the responses to a file
	saveResponseToFile(session.Answers)

	// Send a thank-you message and confirmation
	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("Thank you for your responses!\n\nHere are your answers:\nName: %s\nAmount: %f\nCurrency: %s\nDate: %s\nIs Claimable: %t\nPaid for Family: %t",
			session.Answers.Name, session.Answers.Amount, session.Answers.Currency, session.Answers.Date, session.Answers.IsClaimable, session.Answers.PaidForFamily))
	_, err := bot.Send(msg)
	if err != nil {
		return err
	}

	// Clean up the session
	delete(userSessions, chatID)
	return nil
}

// Save responses to a text file
func saveResponseToFile(response Transaction) {
	file, err := os.OpenFile(SaveFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(file)

	// Serialize the response as JSON and write it to the file
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}

	_, err = file.WriteString(fmt.Sprintf("%s\n", data))
	if err != nil {
		log.Printf("Error writing to file: %v", err)
	}
}

func parseFloat32(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

func parseBool(s string) (bool, error) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return b, nil
}

func processCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) error {
	chatID := callbackQuery.Message.Chat.ID
	session, exists := userSessions[chatID]
	if !exists {
		return fmt.Errorf("session not found for chat ID: %d", chatID)
	}

	// Update the answer based on the button pressed
	switch callbackQuery.Data {
	case "yes":
		if session.CurrentQuestion == 4 {
			session.Answers.IsClaimable = true
		} else if session.CurrentQuestion == 5 {
			session.Answers.PaidForFamily = true
		}
	case "no":
		if session.CurrentQuestion == 4 {
			session.Answers.IsClaimable = false
		} else if session.CurrentQuestion == 5 {
			session.Answers.PaidForFamily = false
		}
	default:
		// Handle currency selection
		if session.CurrentQuestion == 2 {
			session.Answers.Currency = callbackQuery.Data
		}
	}

	// Move to the next question
	session.CurrentQuestion++

	// If all questions are answered, finish the session
	if session.CurrentQuestion >= len(questions) {
		err := finishSession(bot, chatID, session)
		if err != nil {
			return err
		}
		return nil
	}

	// Otherwise, ask the next question
	err := askCurrentQuestion(bot, chatID)
	if err != nil {
		return err
	}

	// Respond to the callback query to remove the loading indicator
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	if _, err := bot.Request(callback); err != nil {
		return err
	}

	return nil
}
