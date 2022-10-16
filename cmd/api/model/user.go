package model

type User struct {
	ID          int     `json:"id"`
	FullName    string  `json:"fullname" validate:"required"`
	Age         int     `json:"age" validate:"required"`
	Email       string  `json:"email" validate:"required"`
	Phone       string  `json:"phone" validate:"required"`
	DateCreated string  `json:"date_created" validate:"required"`
	DateUpdated *string `json:"date_updated"`
}
