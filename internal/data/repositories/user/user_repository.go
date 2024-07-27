package userrepository

import (
	"fmt"
	"http-server/internal/data/database"
	"http-server/internal/models"
	"log"
	"strings"
)

// Dao represents a data access object that interacts with the database.
type UserRepository interface {
	// GetUserById retrieves a user by its ID.
	// It returns an error if the provided ID does not exist.
	GetUserById(id int) (*models.User, error)

	// CreateUser stores a user with the specified username and password.
	// It returns an error in two scenarios:
	// 1. The username is already taken.
	// 2. The password is shorter than 6 characters.
	CreateUser(username, password string) error

	count() int

	deleteAll() error
}

type userRepository struct {
	db database.DbRepository
}

const (
	TABLE_NAME                         = "user"
	GET_USER_BY_ID_ERR                 = "No such ID exists."
	CREATE_USER_USERNAME_TAKEN_ERR     = "Username already exists."
	CREATE_USER_PASSWORD_TOO_SHORT_ERR = "Password must be 6 or more characters."
)

var (
	userRepositoryInstance *userRepository
)

func NewUserRepository() UserRepository {
	if userRepositoryInstance != nil {
		return userRepositoryInstance
	}

	userRepositoryInstance = &userRepository{
		db: database.NewDbRepository(),
	}

	return userRepositoryInstance
}

func (r *userRepository) count() int {
	count, err := r.db.Count(TABLE_NAME)

	if err != nil {
		log.Fatalf("Could not count %s: %s", TABLE_NAME, err)
	}

	return count
}

func (r *userRepository) deleteAll() error {
	return r.db.DeleteAll(TABLE_NAME)
}

func (r *userRepository) GetUserById(id int) (*models.User, error) {
	var userId int
	var username string
	var password string

	if err := getUserByIdStmt.QueryRow(id).Scan(&userId, &username, &password); err != nil {
		var msg string

		if err.Error() == "sql: no rows in result set" {
			msg = GET_USER_BY_ID_ERR
		} else {
			msg = "GetUserById unkown error: " + err.Error()
		}

		return nil, fmt.Errorf(msg)
	}

	return &models.User{Id: userId, Username: username, Password: password}, nil
}

func (r *userRepository) CreateUser(username, password string) error {
	if _, err := createUserStmt.Exec(username, password); err != nil {
		var msg string

		switch {
		case strings.HasPrefix(err.Error(), "UNIQUE"):
			msg = CREATE_USER_USERNAME_TAKEN_ERR
		case strings.HasPrefix(err.Error(), "CHECK"):
			msg = CREATE_USER_PASSWORD_TOO_SHORT_ERR
		default:
			msg = "CreateUser unknown error: " + err.Error()
		}

		return fmt.Errorf(msg)
	}

	return nil
}
