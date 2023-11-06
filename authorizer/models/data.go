package models

type Data struct {
	Users     []User    `json:"users"`
	UserTypes []string  `json:"userTypes"`
	Services  []Service `json:"services"`
}
