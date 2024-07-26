package attacks

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"encoding/json"
	"io/ioutil"
)

// Function to check if the attack time is within the allowed duration and equal to the stored duration
func IsAttackTimeValid(username, time string, db *sql.DB) bool {
	// Query the database to get the duration for the user
	durationQuery := "SELECT `duration` FROM `users` WHERE `username` = ?"
	var durationInt int
	err := db.QueryRow(durationQuery, username).Scan(&durationInt)
	if err != nil {
		log.Printf("[VALIDATE] Error retrieving duration for user %s: %s", username, err)
		return false
	}

	// Convert the provided time to an integer
	timeInt, err := strconv.Atoi(time)
	if err != nil {
		log.Printf("[VALIDATE] Error converting provided time to integer: %s", err)
		return false
	}

	// Check if the provided time is within the allowed duration and less than or equal to the stored duration
	return timeInt <= durationInt
}

// Function to validate if the provided method is in the JSON file
func IsValidMethod(methodName string, filePath string) bool {
	// Read JSON data from file
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("Error reading JSON file:", err)
		return false
	}

	// Parse JSON data into MethodsResponse struct
	var response MethodsResponse
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		log.Println("Error parsing JSON:", err)
		return false
	}

	// Check if the provided method exists in the JSON file
	_, exists := response.METHODS[methodName]
	return exists
}

