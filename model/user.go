package model

import (
	"errors"
)

// User is a user model
type User struct {
	ID   int
	Name string
	Role string
}

// Users is list of user
type Users []User

// Exist check if a user with the given id is exist in the list
func (u Users) Exist(id int) bool {
	for _, user := range u {
		if user.ID == id {
			return true
		}
	}
	return false
}

// FindByName returns the user with the given name. Return error if not found.
func (u Users) FindByName(name string) (User, error) {
	for _, user := range u {
		if user.Name == name {
			return user, nil
		}
	}
	return User{}, errors.New("user_not_found")
}
