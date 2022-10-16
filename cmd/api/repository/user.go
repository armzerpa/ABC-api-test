package repository

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/armzerpa/ABC-api-test/cmd/api/model"
)

type Repo interface {
	GetAll() []model.User
	GetById(id string) *model.User
	DeleteById(id string) bool
	Create(book model.User) *model.User
}

type UserRepo struct {
	DbConnection *sql.DB
}

func NewRepository(db *sql.DB) Repo {
	return UserRepo{DbConnection: db}
}

func (b UserRepo) GetAll() []model.User {
	rows, err := b.DbConnection.Query("SELECT id, fullname, email, phone, date_created, date_updated FROM user")
	if err != nil {
		log.Println("Error in the query to the database")
		return nil
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User

		err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.Phone, &user.DateCreated, &user.DateUpdated)
		if err != nil {
			log.Println("Some error scanning data from the database")
			return nil
		}

		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.Println("Some error selecting data from the database")
		return nil
	}

	return users
}

func (b UserRepo) GetById(id string) *model.User {
	var user model.User
	err := b.DbConnection.QueryRow("SELECT id, fullname, email, phone, date_created, date_updated FROM user WHERE ID = ?", id).Scan(&user.ID, &user.FullName, &user.Email, &user.Phone, &user.DateCreated, &user.DateUpdated)

	if err != nil {
		log.Println("Error in SELECT to the database ", err)
		return nil
	}
	return &user
}

func (b UserRepo) DeleteById(id string) bool {
	err := b.DbConnection.QueryRow("DELETE FROM user WHERE ID = ?", id)

	if err != nil {
		log.Println("Error in the DELETE to the database ", err)
		return false
	}
	return true
}

func (b UserRepo) Create(user model.User) *model.User {
	sql := "INSERT INTO user (fullname, email, phone, date_created, date_updated) VALUES (?,?,?,?,?)"
	res, err := b.DbConnection.Exec(sql, user.FullName, user.Email, user.Phone, user.DateCreated, *user.DateUpdated)
	lastId, err := res.LastInsertId()
	user.ID = strconv.FormatInt(int64(lastId), 10)

	if err != nil {
		log.Println("Error in INSERT to the database ", err)
		return nil
	}
	return &user
}
