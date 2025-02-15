package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
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
	"Date of transaction?",
	"Indicate the date",
	"Is it claimable?",
	"Is it payable for the family?",
}

// Map to track ongoing sessions (active users)
var userSessions = make(map[int64]*UserSession)

// SaveFilePath File to save responses
const SaveFilePath = "responses.txt"

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
	_, err := bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// Process the user's answer and move to the next question
func processAnswer(bot *tgbotapi.BotAPI, chatID int64, session *UserSession, answer string) error {
	switch session.CurrentQuestion {
	case 0: // First question: name of transaction
		session.Answers.Name = answer
	case 1: // Second question: name of transaction (validate as float)
		var name string
		if _, err := fmt.Sscanf(answer, "%d", &name); err != nil {
			// If the input is not valid, ask again
			msg := tgbotapi.NewMessage(chatID, "Please enter a valid number for your amount!")
			_, err = bot.Send(msg)
			if err != nil {
				return err
			}

			return nil
		}
		session.Answers.Name = name
	case 2: // Third question: isClaimable
		session.Answers.IsClaimable = answer
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

	return nil
}

// Finish the Q&A session
func finishSession(bot *tgbotapi.BotAPI, chatID int64, session *UserSession) error {
	// Save the responses to a file
	saveResponseToFile(session.Answers)

	// Send a thank-you message and confirmation
	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("Thank you for your responses, %v!\n\nHere are your answers:\nName: %s\nAge: %d\nCity: %s",
			session.Answers.Amount, session.Answers.Currency, session.Answers.Category, session.Answers.IsClaimable))
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
