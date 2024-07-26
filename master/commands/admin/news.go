package adminCMD

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

// AddNews adds news to the news table
func AdminAddNews(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	if !hasAdminRole(message.From.UserName, db) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You do not have permission to use this command.")
		bot.Send(msg)
		return
	}

	args := message.CommandArguments()

	// Check if there is at least 1 argument
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide the News: /addNews 'news'")
		bot.Send(msg)
		return
	}

	news := args

	username := message.From.UserName
	// Get the current date in a string format
	currentDate := time.Now().Format("2006-01-02")

	_, err := db.Exec("INSERT INTO news (username, date, news) VALUES (?, ?, ?)", username, currentDate, news)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to add news. Please try again.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "News added successfully.")
	bot.Send(msg)
}

// RemoveNews removes news from the news table by specifying the username and news content
func AdminRemoveNews(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB, username, news string) {
	if !hasAdminRole(message.From.UserName, db) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You do not have permission to use this command.")
		bot.Send(msg)
		return
	}
	_, err := db.Exec("DELETE FROM news WHERE username = ? AND news = ?", username, news)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to remove news. Please try again.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "News removed successfully.")
	bot.Send(msg)
}

// DisplayRecentNews displays the most recent 3 news subjects
func News(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	rows, err := db.Query("SELECT date, username, news FROM news ORDER BY date ASC LIMIT 3")
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to fetch news data. Please try again.")
		bot.Send(msg)
		return
	}
	defer rows.Close()

	var newsList strings.Builder

	// Create a tabwriter with padding and alignment
	w := tabwriter.NewWriter(&newsList, 0, 0, 2, ' ', 0)


	// Iterate over the rows and add data to the list
	for rows.Next() {
		var date, username, news string

		err := rows.Scan(&date, &username, &news)
		if err != nil {
			log.Println(err)
			continue
		}

		// Print news data with formatting
		fmt.Fprintf(w, "Date: %s\t | User: %s\t \n  News: %s\t\n", date, username, news)
	}

	// Flush the tabwriter
	w.Flush()

	// Send the formatted news list as a message
	msg := tgbotapi.NewMessage(message.Chat.ID, "Most Recent News:\n"+newsList.String())
	bot.Send(msg)
}
