package models

import "fmt"

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) String() string {
	return fmt.Sprintf("User(Id: %d, Username: %s, Password: %s)", u.Id, u.Username, u.Password)
}
