package model

type File struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required"`
	Path string `json:"path" validate:"required"`
}
