package adminCMD

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

func AdminRemove(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	if !hasAdminRole(message.From.UserName, db) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You do not have permission to use this command.")
		bot.Send(msg)
		return
	}

	args := message.CommandArguments()

	// Split the arguments
	arguments := strings.Fields(args)

	// Check if there are exactly 1 argument (username to remove)
	if len(arguments) != 1 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide exactly 1 argument: /remove 'username'")
		bot.Send(msg)
		return
	}

	usernameToRemove := arguments[0]

	// Delete the user from the database based on username
	_, err := db.Exec("DELETE FROM users WHERE username = ?", usernameToRemove)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to remove user. Please try again.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Removed user with username: %s", usernameToRemove))
	bot.Send(msg)
	log.Printf("Removed user with username: %s", usernameToRemove)
}
