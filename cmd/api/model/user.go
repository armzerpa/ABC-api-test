package model

type User struct {
	ID          string  `json:"id"`
	FullName    string  `json:"fullname"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	DateCreated string  `json:"date_created"`
	DateUpdated *string `json:"date_updated"`
}
