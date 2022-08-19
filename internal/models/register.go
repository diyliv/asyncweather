package models

type Credentials struct {
	Login    string `json:"user_login"`
	Password string `json:"user_password"`
}
