package mainCMD

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Function to handle the "help_button" callback
func HandleHelpButton(username string, db *sql.DB, bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
    // Your logic for handling the "help_button" callback
    // Example: Send the help message
    msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Here are the available commands:")
    msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Admin", "admin_button"),
            tgbotapi.NewInlineKeyboardButtonData("Client", "client_button"),
            tgbotapi.NewInlineKeyboardButtonData("Info", "info_button"),
        ),
    )
    bot.Send(msg)
}