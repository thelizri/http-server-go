package database

import (
	"http-server/models"
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
}

var (
	userRepository *repository
)

func NewUserRepository() UserRepository {
	if userRepository != nil {
		return userRepository
	}

	userRepository = &repository{
		db: db,
	}

	return userRepository
}

func (d *repository) GetUserById(id int) (*models.User, error) {
	var userId int
	var username string
	var password string

	if err := getUserByIdStmt.QueryRow(id).Scan(&userId, &username, &password); err != nil {
		return nil, err
	}

	return &models.User{Id: userId, Username: username, Password: password}, nil
}

func (d *repository) CreateUser(username, password string) error {
	if _, err := createUserStmt.Exec(username, password); err != nil {
		return err
	}

	return nil
}
