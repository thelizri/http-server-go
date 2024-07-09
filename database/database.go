package database

import (
	"context"
	"database/sql"
	"fmt"
	"http-server/models"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

// Dao represents a data access object that interacts with the database.
type Dao interface {
	// GetUserById retrieves a user by its ID.
	// It returns an error if the provided ID does not exist.
	GetUserById(id int) (*models.User, error)

	// CreateUser stores a user with the specified username and password.
	// It returns an error in two scenarios:
	// 1. The username is already taken.
	// 2. The password is shorter than 6 characters.
	CreateUser(username, password string) error

	// Health returns a map of health status information.
	// The keys and values in the map are dao-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type dao struct {
	db *sql.DB
}

var (
	dburl       = os.Getenv("DB_URL")
	daoInstance *dao
)

var (
	createUserStmt  *sql.Stmt
	getUserByIdStmt *sql.Stmt
)

func GetDao() Dao {
	// Reuse Connection
	if daoInstance != nil {
		return daoInstance
	}

	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	prepareStatements(db)

	daoInstance = &dao{
		db: db,
	}
	return daoInstance
}

func prepareStatements(db *sql.DB) {
	prepareCreateUserStmt(db)
	prepareGetUserById(db)
}

func prepareCreateUserStmt(db *sql.DB) {
	query := "INSERT INTO user (username, password) VALUES (?, ?)"

	if stmt, err := db.Prepare(query); err != nil {
		log.Fatal("Could not prepare Create User statement")
	} else {
		createUserStmt = stmt
	}
}

func prepareGetUserById(db *sql.DB) {
	query := "SELECT * FROM user WHERE ID = ?"

	if stmt, err := db.Prepare(query); err != nil {
		log.Fatal("Could not prepare Get User By Id statement")
	} else {
		getUserByIdStmt = stmt
	}
}

func (d *dao) GetUserById(id int) (*models.User, error) {
	var userId int
	var username string
	var password string

	if err := getUserByIdStmt.QueryRow(id).Scan(&userId, &username, &password); err != nil {
		return nil, err
	}

	return &models.User{Id: userId, Username: username, Password: password}, nil
}

func (d *dao) CreateUser(username, password string) error {
	if _, err := createUserStmt.Exec(username, password); err != nil {
		return err
	}

	return nil
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (d *dao) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := d.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := d.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (d *dao) Close() error {
	log.Printf("Disconnected from database: %s", dburl)
	return d.db.Close()
}
