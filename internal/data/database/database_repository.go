package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

// DbRepository represents a data access object that interacts with the database.
type DbRepository interface {
	Prepare(query string) (*sql.Stmt, error)

	// Health returns a map of health status information.
	// The keys and values in the map are dao-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// Counts all rows in a specified table
	Count(table string) (int, error)

	// Deletes all rows in a specified table
	DeleteAll(table string) error
}

type dbRepository struct {
	db *sql.DB
}

func NewDbRepository() DbRepository {
	if dbRepositoryInstance != nil {
		return dbRepositoryInstance
	}

	dbRepositoryInstance = &dbRepository{
		db: db,
	}

	return dbRepositoryInstance
}

func (r *dbRepository) Prepare(query string) (*sql.Stmt, error) {
	return r.db.Prepare(query)
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (r *dbRepository) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := r.db.PingContext(ctx)
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
	dbStats := r.db.Stats()
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
func (r *dbRepository) Close() error {
	log.Printf("Disconnected from database: %s", dburl)
	return r.db.Close()
}

func (r *dbRepository) Count(table string) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	var count int
	err := r.db.QueryRow(query).Scan(&count)

	return count, err
}

func (r *dbRepository) DeleteAll(table string) error {
	query := fmt.Sprintf("DELETE FROM %s", table)
	_, err := r.db.Exec(query)

	return err
}