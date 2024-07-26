package attacks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

// MethodsResponse represents the structure of the JSON response
type MethodsResponse struct {
	METHODS map[string]Method `json:"METHODS"`
}

// Method represents a method in the system
type Method struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Subnet      int    `json:"subnet"`
	MType       int    `json:"mtype"`
}

// ByID is a slice of Method to allow sorting by ID
type ByID []Method

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

// Function to list methods and their descriptions from a JSON file
func listMethodsFromFile(bot *tgbotapi.BotAPI, chatID int64, filePath string) {
	// Read JSON data from file
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("Error reading JSON file:", err)
		return
	}

	// Parse JSON data into MethodsResponse struct
	var response MethodsResponse
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		log.Println("Error parsing JSON:", err)
		return
	}

	// Extract methods into a slice for sorting
	var methodSlice []Method
	for _, method := range response.METHODS {
		methodSlice = append(methodSlice, method)
	}

	// Sort methods by ID
	sort.Sort(ByID(methodSlice))

	// Construct a message with the list of methods and descriptions
	var methodList []string
	for _, method := range methodSlice {
		methodList = append(methodList, fmt.Sprintf("%d: %s - %s", method.ID, method.Name, method.Description))
	}

	// Join the method list into a string
	methodsMessage := strings.Join(methodList, "\n")

	// Send the message back to the chat
	msg := tgbotapi.NewMessage(chatID, methodsMessage)
	bot.Send(msg)
}

// Function to handle the /methods command
func ClientMethods(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB) {

	// Define the file path for the JSON file
	filePath := "master/methods.json"

	// Check if the file exists
	_, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading JSON file %s: %v\n", filePath, err)
		return
	}

	// Call the function to list methods and their descriptions from the JSON file
	listMethodsFromFile(bot, message.Chat.ID, filePath)
}
