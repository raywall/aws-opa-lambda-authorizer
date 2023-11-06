package models

type User struct {
	ID       string `json:"userId"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	UserType string `json:"userType"`
}
