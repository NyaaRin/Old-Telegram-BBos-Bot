package adminCMD

import (
	"database/sql"
	"fmt"
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

func AddAPI(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	if !hasAdminRole(message.From.UserName, db) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You do not have permission to use this command.")
		bot.Send(msg)
		return
	}
	
	args := message.CommandArguments()

	// Check if there is at least 1 argument
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide the API link: /addapi 'apiLINK'")
		bot.Send(msg)
		return
	}

	apiLink := args

	// Insert the API into the database
	_, err := db.Exec("INSERT INTO api (apiLink) VALUES (?)", apiLink)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to add API. Please try again.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Added API: %s", apiLink))
	bot.Send(msg)
	log.Printf("Added API: %s", apiLink)
	
}

// DeleteAPI deletes an API based on its ID
func DeleteAPI(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	if !hasAdminRole(message.From.UserName, db) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You do not have permission to use this command.")
		bot.Send(msg)
		return
	}

	args := message.CommandArguments()

	// Check if there is at least 1 argument
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide the API ID: /deleteapi 'apiID'")
		bot.Send(msg)
		return
	}

	apiID := args

	// Delete the API from the database based on ID
	_, err := db.Exec("DELETE FROM api WHERE id = ?", apiID)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to delete API. Please try again.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Deleted API with ID: %s", apiID))
	bot.Send(msg)
	log.Printf("Deleted API with ID: %s", apiID)
}

// ListAPIs lists all APIs in the database
func ListAPIs(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	if !hasAdminRole(message.From.UserName, db) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You do not have permission to use this command.")
		bot.Send(msg)
		return
	}

	rows, err := db.Query("SELECT id, apiLink FROM api")
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to retrieve API list. Please try again.")
		bot.Send(msg)
		return
	}
	defer rows.Close()

	var apiList string
	for rows.Next() {
		var id sql.NullInt64
		var apiLink string
		err := rows.Scan(&id, &apiLink)
		if err != nil {
			log.Println(err)
			continue
		}

		idStr := "NULL"
		if id.Valid {
			idStr = fmt.Sprintf("%d", id.Int64)
		}

		apiList += fmt.Sprintf("ID: %s, API Link: %s\n", idStr, apiLink)
	}

	if apiList == "" {
		apiList = "No APIs found."
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, apiList)
	bot.Send(msg)
}