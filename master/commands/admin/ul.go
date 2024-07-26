package adminCMD

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

func AdminUL(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	rows, err := db.Query("SELECT role, username FROM users ORDER BY role ASC ")
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to fetch user data. Please try again.")
		bot.Send(msg)
		return
	}
	defer rows.Close()

	var userList strings.Builder

	// Iterate over the rows and add data to the list
	for rows.Next() {
		var role, username string

		err := rows.Scan(&role, &username)
		if err != nil {
			log.Println(err)
			continue
		}

		// Append the data to the list
		userList.WriteString(fmt.Sprintf("Role: %s\t | Usernames: %s\n", role, username))
	}

	// Send the user list as a message
	msg := tgbotapi.NewMessage(message.Chat.ID, "List of all users:\n"+userList.String())
	bot.Send(msg)
}

func AdminUserInfo(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	// Extract the username from the command arguments
	args := message.CommandArguments()
	username := args

	// Query the database for the user information
	var role, plan string
	var expiry, concurrents, duration int

	err := db.QueryRow("SELECT role, plan, expiry, concurrents, duration FROM users WHERE username = ?", username).
		Scan(&role, &plan, &expiry, &concurrents, &duration)

	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "User not found or an error occurred.")
		bot.Send(msg)
		return
	}

	// Create a message with the user information
	userInfo := fmt.Sprintf("User Information for: %s\n"+
		"Role: %s\n"+
		"Plan: %s\n"+
		"Expiry: %d\n"+
		"Concurrents: %d\n"+
		"Duration: %d",
		username, role, plan, expiry, concurrents, duration)

	// Send the user information as a message
	msg := tgbotapi.NewMessage(message.Chat.ID, userInfo)
	bot.Send(msg)
}
