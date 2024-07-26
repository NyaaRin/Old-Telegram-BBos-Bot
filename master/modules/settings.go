package User

import (

	"encoding/json"
	"database/sql"
	"io/ioutil"
	"log"
	"fmt"
)

// Config struct to hold the settings
type Config struct {
	Live bool `json:"live"`
}

// LoadConfig reads the settings from the settings.json file
func LoadConfig() (*Config, error) {
	file, err := ioutil.ReadFile("master/settings.json")
	if err != nil {
		log.Println("Error reading settings file:", err)
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Println("Error unmarshalling settings:", err)
		return nil, err
	}

	return &config, nil
}

// isAdminAllowed checks if the user has admin role or is allowed in non-live mode
func IsAdminAllowed(username string, db *sql.DB) bool {
	config, err := LoadConfig()
	if err != nil {
		// Handle error (you might want to return false or handle it differently)
		fmt.Println("Error loading config:", err)
		return false
	}

	if config.Live {
		// Check admin role only if "live" is true
		result := hasAdminRole(username, db)
		fmt.Printf("hasAdminRole result: %v\n", result)
		return result
	}

	// Allow the user in non-live mode without checking admin role
	fmt.Println("Non-live mode, allowing access.")
	return true
}


// hasAdminRole checks if the user has admin role
func hasAdminRole(username string, db *sql.DB) bool {
	var role int
	err := db.QueryRow("SELECT role FROM users WHERE username = ? LIMIT 1", username).Scan(&role)
	if err != nil {
		return false
	}
	return role == 1
}