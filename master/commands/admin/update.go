package adminCMD

import (
	"database/sql"
	"fmt"
	"strings"
	"TeleBot/master/modules"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"

)

// Function to handle the /updateuser command
func Update(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	userName := message.From.UserName
	if !User.IsAdminAllowed(userName, db) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You do not have permission to use this command.")
		bot.Send(msg)
		return
	}
	
	args := message.CommandArguments()

	// Split the arguments
	arguments := strings.Fields(args)

	// Check if there are enough arguments
	if len(arguments) < 3 {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide at least 3 arguments: /update [Username] [Table] [Value] [additional_fields]")
		bot.Send(msg)
	example := tgbotapi.NewMessage(message.Chat.ID, "/update FileGone rank 1 || /update FileGone rank 1 plan test\n Tables: duration, plan, role, expiry, concurrents")
		bot.Send(example)
		return
	}

	username := arguments[0]
	field := arguments[1]
	value := arguments[2]

	// Update user's information in the database
	err := updateUserInfo(username, field, value, db, arguments[3:]...)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Error updating user's information: %v", err))
		bot.Send(msg)
		return
	}

	// Send a success message
	successMsg := fmt.Sprintf("User '%s' information updated - %s: '%s'", username, field, value)
	msg := tgbotapi.NewMessage(message.Chat.ID, successMsg)
	bot.Send(msg)
}

// Function to update user's information in the database
func updateUserInfo(username, field, value string, db *sql.DB, additionalFields ...string) error {
	// Build SQL query based on the field to be updated
	var query string
	switch field {
	case "role", "plan":
		query = fmt.Sprintf("UPDATE users SET %s = ? WHERE username = ?", field)
	default:
		// Handle additional fields
		query = "UPDATE users SET " + field + " = ?"
		for _, additionalField := range additionalFields {
			query += ", " + additionalField + " = ?"
		}
		query += " WHERE username = ?"
	}

	// Prepare SQL statement
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Build parameters for the SQL statement
	var params []interface{}
	params = append(params, value)
	for _, additionalField := range additionalFields {
		params = append(params, additionalField)
	}
	params = append(params, username)

	// Execute the update statement
	_, err = stmt.Exec(params...)
	return err
}
