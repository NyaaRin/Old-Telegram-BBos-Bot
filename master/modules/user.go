// users.go

package User

import (
	"database/sql"
)

// DBInterface defines the database operations interface
type DBInterface interface {
	GetUser(username string) (*User, error)
	// Add other database operations as needed
}

type User struct {
	ID          int
	Plan        string
	Role        string
	Expiry      string
	Concurrents int
	Duration    int
}

func mapRoleToString(role string) string {
	switch role {
	case "1":
		return "Admin"
	case "2":
		return "Client"
	default:
		return "User"
	}
}

// GetUser retrieves user information from the database
func GetUser(username string, db *sql.DB) (*User, error) {
    row := db.QueryRow("SELECT id, plan, role, expiry, concurrents, duration FROM users WHERE username = ? LIMIT 1", username)
    var user User
    err := row.Scan(&user.ID, &user.Plan, &user.Role, &user.Expiry, &user.Concurrents, &user.Duration)
    if err != nil {
        if err == sql.ErrNoRows {
            // User not found, create a default user object
            user = User{
                ID:          0,        // Set default ID or any other default values
                Plan:        "None",   // Set default plan or any other default values
                Role:        "User",   // Set default role or any other default values
                Expiry:      "00-00-0000",       // Set default expiry or any other default values
                Concurrents: 0,        // Set default concurrents or any other default values
                Duration:    0,        // Set default duration or any other default values
            }
            return &user, nil
        }
        return nil, err // Return the actual error
    }

    // Map role values to human-readable strings
    user.Role = mapRoleToString(user.Role)

    return &user, nil
}