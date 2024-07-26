package clientCMD

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

// Function to check if a user has a valid plan
func ClientDash(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, db *sql.DB) {
	username := callback.From.UserName
	chatID := callback.Message.Chat.ID // Use callback.Message.Chat.ID to get the chat ID

	if !hasPlan(username, db) {
		msg := tgbotapi.NewMessage(chatID, "User does not have a valid plan to use client commands.")
		bot.Send(msg)
		return
	}

	// Create inline keyboard with a back button
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back to Help", "help_button"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, "+==========+\n+Client Dash!+\n+==========+\n/attack\n/methods\n/info\n/start")
	msg.ReplyMarkup = inlineKeyboard
	bot.Send(msg)
}

func hasPlan(username string, db *sql.DB) bool {
	var plan string
	err := db.QueryRow("SELECT plan FROM users WHERE username = ? LIMIT 1", username).Scan(&plan)
	if err != nil {
		return false
	}
	return plan != "none"
}
