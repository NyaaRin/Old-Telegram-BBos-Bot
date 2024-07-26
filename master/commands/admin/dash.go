package adminCMD

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// AdminDash handles the admin dashboard commands
func AdminDash(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, db *sql.DB) {
	username := callback.From.UserName
	chatID := callback.Message.Chat.ID

	if !hasAdminRole(username, db) {
		msg := tgbotapi.NewMessage(chatID, "You do not have permission to use this command.")
		bot.Send(msg)
		return
	}

	// Create inline keyboard with a back button
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back to Help", "help_button"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, "+==========+\n+Admin Dash+\n+==========+\n/add\n/remove\n/users\n/user [username]\n/update\n/addAPI")
	msg.ReplyMarkup = inlineKeyboard
	bot.Send(msg)
}

// hasAdminRole checks if the user has admin role
func hasAdminRole(username string, db *sql.DB) bool {
	var role int
	err := db.QueryRow("SELECT role FROM users WHERE username = ? LIMIT 1", username).Scan(&role)
	if err != nil {
		return false
	}
	return role == 1
}
