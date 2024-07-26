package mainCMD

import (
	"TeleBot/master/modules"
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// Function to check if a user has a valid plan
func ClientInfo(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, db *sql.DB) {
	// Validate that the user has a plan to use client commands
	username := callback.From.UserName
	chatID := callback.Message.Chat.ID

	// Use dependency injection to pass the database instance to GetUser
	foundUser, err := User.GetUser(username, db)
	if err != nil {
		log.Fatal(err)
	}

	// Create inline keyboard with a back button
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back to Help", "help_button"),
		),
	)

	userPlan := foundUser.Plan
	userRole := foundUser.Role
	userExpiry := foundUser.Expiry
	userCons := foundUser.Concurrents
	userMBT := foundUser.Duration

	// Construct a message with the user's information
	responseMsg := fmt.Sprintf(
		"Welcome, %s!\n"+
		"Your Plan: %s\n"+
		"Your Role: %s\n"+
		"Your Expiry: %s\n"+
		"Your Cons: %d\n"+
		"Your Max Time: %d\n",
		username, userPlan, userRole, userExpiry, userCons, userMBT)
	msg := tgbotapi.NewMessage(chatID, responseMsg)
	msg.ReplyMarkup = inlineKeyboard
	bot.Send(msg)
}
