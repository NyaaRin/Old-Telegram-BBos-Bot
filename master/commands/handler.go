package cmd

import (
	"database/sql"
	"TeleBot/master/commands/admin"
	"TeleBot/master/commands/client"
	"fmt"
	"log"
	"TeleBot/master/modules"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"TeleBot/master/commands/client/attacks"
)

// Function to handle commands
func HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	// Validate that the user has a plan to use client commands
	username := message.From.UserName

	// Use GetUser method from the user.Instance
	foundUser, err := User.GetUser(username, db)
	if err != nil {
		log.Println("User not found")
	}

	// Access the user ID field
	responseMsg := fmt.Sprintf("Welcome, %s!\nYour Plan: %s\nYour Role: %s\n", username, foundUser.Plan, foundUser.Role)

	switch message.Command() {
	//START MEE
	case "start":
		msg := tgbotapi.NewMessage(message.Chat.ID, responseMsg)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Click me", "help_button"),
			),
		)
		bot.Send(msg)
	
	//ADMIN SIDE
	case "add":
		adminCMD.Add(bot, message, db)
	case "news":
		adminCMD.News(bot, message, db)
	case "remove":
		adminCMD.AdminRemove(bot, message, db)
	case "users":
		adminCMD.AdminUL(bot, message, db)
	case "user":
		adminCMD.AdminUserInfo(bot, message, db)
	case "update":
		adminCMD.Update(bot, message, db)
	case "addAPI":
		adminCMD.AddAPI(bot, message, db)
	case "removeAPI":
		adminCMD.DeleteAPI(bot, message, db)
	case "APIs":
		adminCMD.ListAPIs(bot, message, db)
	case "addNews":
		adminCMD.AdminAddNews(bot, message, db)
	case "removeNews":
		adminCMD.AdminAddNews(bot, message, db)
		
	//CLIENT SIDE
	case "method":
		attacks.ClientMethods(bot, message, db)
	case "attack":
		clientCMD.ClientAttack(bot, message, db)

    // Handle unknown data
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "I don't know that command.")
		bot.Send(msg)
	}

	log.Printf("[USER] command: %s | From User: %s", message.Command(), username)
}

