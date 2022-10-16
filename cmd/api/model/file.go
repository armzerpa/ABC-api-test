package model

type File struct {
	ID   int    `json:"id"`
	Name string `json:"fullname" validate:"required"`
	Path int    `json:"age" validate:"required"`
}
