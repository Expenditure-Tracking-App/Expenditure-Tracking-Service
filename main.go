package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

// Session to track user's Q&A flow
type UserSession struct {
	CurrentQuestion int         // Index of the current question
	Answers         Transaction // Struct to store responses
}

// Questions array for the process
var questions = []string{
	"What is your name?",
	"How old are you? (please enter a number)",
	"What city do you live in?",
}

// Map to track ongoing sessions (active users)
var userSessions = make(map[int64]*UserSession)

// File to save responses
const SaveFilePath = "responses.txt"

func main() {
	token, err := os.ReadFile("token.txt")

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
			processMessage(bot, update.Message)
		}
	}
}

// Handle incoming messages
func processMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID

	// Handle start command to begin the Q&A process
	if message.Text == "/start" {
		startSession(bot, chatID, message.From.UserName)
		return
	}

	// If the user is in the middle of a session, process their answer
	if session, exists := userSessions[chatID]; exists {
		processAnswer(bot, chatID, session, message.Text)
		return
	}

	// If no session is active, guide the user
	msg := tgbotapi.NewMessage(chatID, "Send /start to begin!")
	bot.Send(msg)
}

// Start a new Q&A session
func startSession(bot *tgbotapi.BotAPI, chatID int64, username string) {
	// Create a session for the user
	userSessions[chatID] = &UserSession{
		CurrentQuestion: 0, // Start at the first question
		Answers: Transaction{
			Amount: username,
		},
	}

	// Ask the first question
	askCurrentQuestion(bot, chatID)
}

// Ask the current question in the session
func askCurrentQuestion(bot *tgbotapi.BotAPI, chatID int64) {
	session := userSessions[chatID]
	question := questions[session.CurrentQuestion]

	// Send the current question to the user
	msg := tgbotapi.NewMessage(chatID, question)
	bot.Send(msg)
}

// Process the user's answer and move to the next question
func processAnswer(bot *tgbotapi.BotAPI, chatID int64, session *UserSession, answer string) {
	switch session.CurrentQuestion {
	case 0: // First question: currency
		session.Answers.Currency = answer
	case 1: // Second question: category (validate as integer)
		var age int
		if _, err := fmt.Sscanf(answer, "%d", &age); err != nil {
			// If the input is not valid, ask again
			msg := tgbotapi.NewMessage(chatID, "Please enter a valid number for your age!")
			bot.Send(msg)
			return
		}
		session.Answers.Category = age
	case 2: // Third question: isClaimable
		session.Answers.IsClaimable = answer
	}

	// Move to the next question
	session.CurrentQuestion++

	// If all questions are answered, finish the session
	if session.CurrentQuestion >= len(questions) {
		finishSession(bot, chatID, session)
		return
	}

	// Otherwise, ask the next question
	askCurrentQuestion(bot, chatID)
}

// Finish the Q&A session
func finishSession(bot *tgbotapi.BotAPI, chatID int64, session *UserSession) {
	// Save the responses to a file
	saveResponseToFile(session.Answers)

	// Send a thank-you message and confirmation
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Thank you for your responses, %s!\n\nHere are your answers:\nName: %s\nAge: %d\nCity: %s",
		session.Answers.Amount, session.Answers.Currency, session.Answers.Category, session.Answers.IsClaimable))
	bot.Send(msg)

	// Clean up the session
	delete(userSessions, chatID)
}

// Save responses to a text file
func saveResponseToFile(response Transaction) {
	file, err := os.OpenFile(SaveFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close()

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
