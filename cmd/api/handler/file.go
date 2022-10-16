package handler

import (
	"database/sql"
	"net/http"

	"github.com/armzerpa/ABC-api-test/cmd/api/model"
	"github.com/armzerpa/ABC-api-test/cmd/api/repository"

	"github.com/gin-gonic/gin"
)

func NewFileHandler(db *sql.DB) *HandlerFile {
	return &HandlerFile{
		fileRepository: repository.NewFileRepository(db),
		userRepository: repository.NewRepository(db),
	}
}

type HandlerFile struct {
	fileRepository repository.File
	userRepository repository.Repo
}

func (h *HandlerFile) GetByUserId(c *gin.Context) {
	id := c.Param("id")
	user := h.userRepository.GetById(id)
	if user == nil {
		c.IndentedJSON(http.StatusNotFound, model.Message{Message: "user not found", Status: "not_found"})
		return
	}

	files := h.fileRepository.GetAll(id)
	if files == nil {
		c.IndentedJSON(http.StatusNotFound, model.Message{Message: "not files found", Status: "not_found"})
	} else {
		c.IndentedJSON(http.StatusOK, files)
	}
}

func (h *HandlerFile) DeleteByUserId(c *gin.Context) {
	id := c.Param("id")
	user := h.userRepository.GetById(id)
	if user == nil {
		c.IndentedJSON(http.StatusNotFound, model.Message{Message: "user not found", Status: "not_found"})
		return
	}

	result := h.fileRepository.DeleteAll(id)
	if result {
		c.IndentedJSON(http.StatusOK, model.Message{Message: "files deleted successfully", Status: "user remove file, id: " + id})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, model.Message{Message: "error deleting files", Status: "internal_error"})
	}
}

func (h *HandlerFile) CreateFile(c *gin.Context) {
	id := c.Param("id")
	user := h.userRepository.GetById(id)
	if user == nil {
		c.IndentedJSON(http.StatusNotFound, model.Message{Message: "user not found", Status: "not_found"})
		return
	}

	var fileToInsert model.File
	error := c.BindJSON(&fileToInsert)
	if error != nil {
		c.IndentedJSON(http.StatusBadRequest, model.Message{Message: "invalid body file", Status: "bad_request"})
		return
	}

	file := h.fileRepository.Create(fileToInsert, id)
	if file == nil {
		c.IndentedJSON(http.StatusInternalServerError, model.Message{Message: "error inserting file", Status: "internal_error"})
		return
	}
	c.IndentedJSON(http.StatusCreated, file)
}
