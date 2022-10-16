package repository

import (
	"database/sql"
	"log"

	"github.com/armzerpa/ABC-api-test/cmd/api/model"
)

type File interface {
	GetAll(userId string) []model.File
	DeleteAll(userId string) bool
	Create(file model.File, userId string) *model.File
}

type FileRepo struct {
	DbConnection *sql.DB
}

func NewFileRepository(db *sql.DB) File {
	return FileRepo{DbConnection: db}
}

func (b FileRepo) GetAll(userId string) []model.File {
	rows, err := b.DbConnection.Query("SELECT id, name, path FROM file WHERE userid = ?", userId)
	if err != nil {
		log.Println("Error in the query to the database")
		return nil
	}
	defer rows.Close()

	var Files []model.File

	for rows.Next() {
		var File model.File

		err := rows.Scan(&File.ID, &File.Name, &File.Path)
		if err != nil {
			log.Println("Some error scanning data from the database")
			return nil
		}

		Files = append(Files, File)
	}
	if err = rows.Err(); err != nil {
		log.Println("Some error selecting files from the database")
		return nil
	}

	return Files
}

func (b FileRepo) DeleteAll(userid string) bool {
	sql := "DELETE FROM file WHERE userid = ?"
	_, err := b.DbConnection.Exec(sql, userid)

	if err != nil {
		log.Println("Error in the DELETE to the database ", err)
		return false
	}
	return true
}

func (b FileRepo) Create(file model.File, userId string) *model.File {
	sql := "INSERT INTO file (name, path, userid) VALUES (?,?,?)"
	res, err := b.DbConnection.Exec(sql, file.Name, file.Path, userId)
	lastId, err := res.LastInsertId()
	file.ID = int(lastId)

	if err != nil {
		log.Println("Error in INSERT file to the database ", err)
		return nil
	}
	return &file
}
