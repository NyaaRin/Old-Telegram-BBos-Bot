package adminCMD

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"TeleBot/master/modules"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

func Add(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	userName := message.From.UserName
	if !User.IsAdminAllowed(userName, db) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You do not have permission to use this command.")
		bot.Send(msg)
		return
	}


	args := message.CommandArguments()

	// Split the arguments
	arguments := strings.Fields(args)

	// Check if there are exactly 3 arguments
	if len(arguments) != 3 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide exactly 3 arguments: /add 'username' 'role' 'plan'")
		bot.Send(msg)
		return
	}

	username := arguments[0]
	role := arguments[1]
	plan := arguments[2]


	// Insert the user into the database
    _, err := db.Exec("INSERT INTO users (username, role, plan, expiry, concurrents, duration ) VALUES (?, ?, ?, 0, 0, 0)", username, role, plan)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to add user. Please try again.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Added user: Username - %s, Role - %s", username, role))
	bot.Send(msg)
	log.Printf("Added user: Username - %s, Role - %s", username, role)

}

