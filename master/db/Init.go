package database

import (
	"database/sql"
	"log"
	"TeleBot/master/commands"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)


func InitBot(db *sql.DB) {
	// Set your bot token here
	bot, err := tgbotapi.NewBotAPI("7046216198:AAEgKjf5UbUTNf-R3MChuF-cweqBsjCo_I0")
	if err != nil {
		log.Fatal(err)
	}
	// Enable logging
	log.Printf("Authorized on account %s", bot.Self.UserName)

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		
		updates, err := bot.GetUpdatesChan(u)
		if err != nil {
			log.Fatal(err)
		}
		
		// Handle incoming updates
		for update := range updates {
			if update.Message != nil {
				// Handle regular messages
				cmd.HandleCommand(bot, update.Message, db)
			} else if update.CallbackQuery != nil {
				// Handle callback queries
				cmd.HandleCallback(bot, update.CallbackQuery, db)
			}
		}
}