package repository

import (
	"database/sql"

	"github.com/armzerpa/ABC-api-test/cmd/api/model"
	"github.com/sirupsen/logrus"
)

type Repo interface {
	GetAll() []model.User
	GetById(id string) *model.User
	GetByEmail(email string) *model.User
	DeleteById(id string) bool
	Create(book model.User) *model.User
	Update(book model.User) bool
}

type UserRepo struct {
	DbConnection *sql.DB
}

func NewRepository(db *sql.DB) Repo {
	return UserRepo{DbConnection: db}
}

func (b UserRepo) GetAll() []model.User {
	rows, err := b.DbConnection.Query("SELECT id, fullname, age, email, phone, date_created, date_updated FROM user")
	if err != nil {
		log.WithFields(logrus.Fields{}).Error("Error in the query to the database ", err)
		return nil
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User

		err := rows.Scan(&user.ID, &user.FullName, &user.Age, &user.Email, &user.Phone, &user.DateCreated, &user.DateUpdated)
		if err != nil {
			log.WithFields(logrus.Fields{}).Error("Some error scanning data from the database ", err)
			return nil
		}

		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.WithFields(logrus.Fields{}).Error("Some error selecting data from the database ", err)
		return nil
	}

	return users
}

func (b UserRepo) GetById(id string) *model.User {
	var user model.User
	err := b.DbConnection.QueryRow("SELECT id, fullname, age, email, phone, date_created, date_updated FROM user WHERE id = ?", id).Scan(&user.ID, &user.FullName, &user.Age, &user.Email, &user.Phone, &user.DateCreated, &user.DateUpdated)

	if err != nil {
		log.WithFields(logrus.Fields{}).Error("Error in SELECT to the database ", err)
		return nil
	}
	return &user
}

func (b UserRepo) GetByEmail(email string) *model.User {
	var user model.User
	err := b.DbConnection.QueryRow("SELECT id, fullname, age, email, phone, date_created, date_updated FROM user WHERE email = ?", email).Scan(&user.ID, &user.FullName, &user.Age, &user.Email, &user.Phone, &user.DateCreated, &user.DateUpdated)

	if err != nil {
		log.WithFields(logrus.Fields{}).Error("Error in SELECT to the database ", err)
		return nil
	}
	return &user
}

func (b UserRepo) DeleteById(id string) bool {
	sql := "DELETE FROM user WHERE ID = ?"
	_, err := b.DbConnection.Exec(sql, id)

	if err != nil {
		log.WithFields(logrus.Fields{}).Error("Error in the DELETE to the database ", err)
		return false
	}
	return true
}

func (b UserRepo) Create(user model.User) *model.User {
	sql := "INSERT INTO user (fullname, age, email, phone, date_created, date_updated) VALUES (?,?,?,?,?,?)"
	res, err := b.DbConnection.Exec(sql, user.FullName, user.Age, user.Email, user.Phone, user.DateCreated, user.DateCreated)
	lastId, err := res.LastInsertId()
	user.ID = int(lastId)

	if err != nil {
		log.WithFields(logrus.Fields{}).Error("Error in INSERT to the database ", err)
		return nil
	}
	return &user
}

func (b UserRepo) Update(user model.User) bool {
	sql := "UPDATE user SET fullname=?, age=?, phone=?, date_updated=? WHERE id=?;"
	_, err := b.DbConnection.Exec(sql, user.FullName, user.Age, user.Phone, *user.DateUpdated, user.ID)
	if err != nil {
		log.WithFields(logrus.Fields{}).Error("Error in UPDATE to the database ", err)
		return false
	}
	return true
}
