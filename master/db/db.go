// db.go

package database

import (
	"database/sql"
	"log"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"TeleBot/master/modules"
	"fmt"
)

type Instance struct {
	Connected time.Time
	Conn      *sql.DB
}

var (
	Container *Instance = NewDB()
)

func NewDB() *Instance {
	db := &Instance{
		Connected: time.Now(),
	}

	SQL, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	db.Conn = SQL

	log.Printf("[NEW] DATABASE CONNECTED")

	// Create the users table if it doesn't exist
	_, err = db.Conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			username STRING,
			role INT,
			concurrents INT,
			plan STRING,
			expiry STRING,
			duration INT
		);
		
		CREATE TABLE IF NOT EXISTS api (
			id INTEGER PRIMARY KEY,
			apiLink VARCHAR(255) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS news (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username STRING,
			date STRING,
			news STRING
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	
	InitBot(db.Conn)
	return db
}

func (inst *Instance) GetUser(username string) (*User.User, error) { // Fully qualify User package
	// Implement the GetUser method
	// Example: return db.Instance.GetUser(username)
	return nil, fmt.Errorf("not implemented")
}