package cmd

import (
	"database/sql"
	"TeleBot/master/commands/admin"
	"TeleBot/master/commands/client"
	"TeleBot/master/commands/main"
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Function to handle callback queries
func HandleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, db *sql.DB) {
	
	username := callback.From.UserName
	data := callback.Data

	switch data {
	case "help_button": 
		mainCMD.HandleHelpButton(username, db, bot, callback)
	case "info_button":
		mainCMD.ClientInfo(bot, callback, db)

	// Call the AdminDash function from adminCMD package
	case "admin_button":
		adminCMD.AdminDash(bot, callback, db)
		

	// Call the ClientDash function from clientCMD package
	case "client_button":
		clientCMD.ClientDash(bot, callback, db) 

	
	default:
		// Handle unknown callback data
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Unknown callback data")
		bot.Send(msg)
	}
	log.Printf("[CALLBACK] data: %s | From User: %s", data, username)
	log.Printf("Switch data: %s", data)
}
