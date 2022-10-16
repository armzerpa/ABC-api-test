package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/armzerpa/ABC-api-test/cmd/api/model"
	"github.com/armzerpa/ABC-api-test/cmd/api/repository"
	"github.com/go-playground/validator"

	"github.com/gin-gonic/gin"
)

func NewUserHandler(db *sql.DB) *HandlerUser {
	return &HandlerUser{
		userRepository: repository.NewRepository(db),
	}
}

type HandlerUser struct {
	userRepository repository.Repo
}

func (h *HandlerUser) GetUsers(c *gin.Context) {
	users := h.userRepository.GetAll()
	if users == nil {
		c.IndentedJSON(http.StatusNotFound, model.Message{Message: "not users found", Status: "not_found"})
	} else {
		c.IndentedJSON(http.StatusOK, users)
	}
}

func (h *HandlerUser) GetUserById(c *gin.Context) {
	id := c.Param("id")
	user := h.userRepository.GetById(id)
	if user == nil {
		c.IndentedJSON(http.StatusNotFound, model.Message{Message: "user not found", Status: "not_found"})
	} else {
		c.IndentedJSON(http.StatusOK, &user)
	}
}

func (h *HandlerUser) DeleteUserById(c *gin.Context) {
	id := c.Param("id")
	user := h.userRepository.GetById(id)
	if user == nil {
		c.IndentedJSON(http.StatusNotFound, model.Message{Message: "user not found", Status: "not_found"})
		return
	}

	result := h.userRepository.DeleteById(id)
	if result {
		c.IndentedJSON(http.StatusOK, model.Message{Message: "user deleted successfully", Status: "user remove, id: " + id})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, model.Message{Message: "error deleting user", Status: "internal_error"})
	}
}

func (h *HandlerUser) CreateUser(c *gin.Context) {
	var userToInsert model.User
	error := c.BindJSON(&userToInsert)
	if error != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "invalid body", Status: "bad_request"})
		return
	}
	if !h.validateUserInput(c, userToInsert, true) {
		return
	}

	user := h.userRepository.Create(userToInsert)
	if user == nil {
		c.IndentedJSON(http.StatusInternalServerError, model.Message{Message: "error inserting user", Status: "internal_error"})
		return
	}
	c.IndentedJSON(http.StatusCreated, user)
}

func (h *HandlerUser) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	user := h.userRepository.GetById(id)
	if user == nil {
		c.IndentedJSON(http.StatusNotFound, model.Message{Message: "user not found", Status: "not_found"})
		return
	}

	var userToUpdate model.User
	error := c.BindJSON(&userToUpdate)
	if error != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "invalid body", Status: "bad_request"})
		return
	}
	if !h.validateUserInput(c, userToUpdate, false) {
		return
	}

	userToUpdate.ID, _ = strconv.Atoi(id)
	currentTime := time.Now()
	t := currentTime.Format("2006-01-02")
	userToUpdate.DateUpdated = &t
	result := h.userRepository.Update(userToUpdate)
	if result {
		c.IndentedJSON(http.StatusOK, userToUpdate)
	} else {
		c.IndentedJSON(http.StatusInternalServerError, model.Message{Message: "error updating user", Status: "internal_error"})
	}
}

func (h *HandlerUser) validateUserInput(c *gin.Context, userToInsert model.User, isInsert bool) bool {
	if isInsert {
		validate := validator.New()
		err := validate.Struct(userToInsert)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "missing required field", Status: "bad_request"})
			return false
		}
		user := h.userRepository.GetByEmail(userToInsert.Email)
		if user != nil {
			c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "user email must be unique", Status: "bad_request"})
			return false
		}
	} else {
		if userToInsert.Email != "" {
			c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "cannot update email", Status: "bad_request"})
			return false
		}
	}
	if userToInsert.Age < 18 {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "user cannot be less than 18 years old", Status: "bad_request"})
		return false
	}
	return true
}
