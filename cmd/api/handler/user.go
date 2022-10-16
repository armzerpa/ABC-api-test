package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/armzerpa/ABC-api-test/cmd/api/model"
	"github.com/armzerpa/ABC-api-test/cmd/api/repository"

	"github.com/gin-gonic/gin"
)

func NewBookHandler(db *sql.DB) *HandlerBook {
	return &HandlerBook{
		userRepository: repository.NewRepository(db),
	}
}

type HandlerBook struct {
	userRepository repository.Repo
}

func (h *HandlerBook) GetUsers(c *gin.Context) {
	books := h.userRepository.GetAll()
	c.IndentedJSON(http.StatusOK, books)
}

func (h *HandlerBook) GetUserById(c *gin.Context) {
	id := c.Param("id")
	book := h.userRepository.GetById(id)
	c.IndentedJSON(http.StatusOK, &book)
}

func (h *HandlerBook) DeleteUserById(c *gin.Context) {
	id := c.Param("id")
	result := h.userRepository.DeleteById(id)
	c.IndentedJSON(http.StatusOK, "Book deleted: "+fmt.Sprintf("%t", result))
}

func (h *HandlerBook) CreateUser(c *gin.Context) {
	var userToInsert model.User
	error := c.BindJSON(&userToInsert)
	if error != nil {
		c.IndentedJSON(http.StatusBadRequest, "Invalid body")
		return
	}

	user := h.userRepository.Create(userToInsert)
	c.IndentedJSON(http.StatusCreated, user)
}
