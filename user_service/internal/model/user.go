package model

type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserDB struct {
	Email        string
	Name         string
	PasswordHash string
}
