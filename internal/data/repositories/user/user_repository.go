package userrepository

import (
	"http-server/internal/data/database"
	"http-server/internal/models"
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

type userRepository struct {
	db database.DbRepository
}

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

func (r *userRepository) GetUserById(id int) (*models.User, error) {
	var userId int
	var username string
	var password string

	if err := getUserByIdStmt.QueryRow(id).Scan(&userId, &username, &password); err != nil {
		return nil, err
	}

	return &models.User{Id: userId, Username: username, Password: password}, nil
}

func (r *userRepository) CreateUser(username, password string) error {
	if _, err := createUserStmt.Exec(username, password); err != nil {
		return err
	}

	return nil
}
