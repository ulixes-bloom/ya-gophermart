package models

import (
	"fmt"
)

//easyjson:json
type User struct {
	ID       int64  `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *User) String() string {
	if u == nil {
		return "user is nil pointer"
	}

	return fmt.Sprintf("ID: %d, Login: %s", u.ID, u.Login)
}
