package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type Gender string

const (
	Female Gender = "female"
	Male   Gender = "male"
)

type Users []User

type User struct {
	ID        int    `json:"id"`
	FirstName string `form:"first_name" json:"first_name"`
	LastName  string `form:"last_name" json:"last_name"`
	Email     string `form:"email" json:"email"`
	Gender    Gender `form:"gender" json:"gender"`
	Age       int    `form:"age" json:"age"`
}

func ListUser(db *sql.DB, w string) (*Users, error) {
	smt := "select * from users"
	if w != "" {
		smt += fmt.Sprintf(" where %s", w)
	}
	fmt.Println(smt)
	rows, err := db.Query(smt)
	if err != nil {
		return nil, err
	}
	var users Users
	for rows.Next() {
		var tempUser User
		_ = rows.Scan(&tempUser.ID, &tempUser.FirstName, &tempUser.LastName, &tempUser.Email, &tempUser.Gender, &tempUser.Age)
		users = append(users, tempUser)
	}
	if err := rows.Close(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &users, nil
}

func (u *User) Detail(db *sql.DB, id string) error {
	row := db.QueryRow("select * from users where id = ?", id)
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Gender, &u.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Zero rows found")
		} else {
			return err
		}
	}
	defer db.Close()
	return nil
}

func (u *User) Create(db *sql.DB) error {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into users (first_name ,last_name ,email ,gender ,age) values (?,?,?,?,?)")
	row, err := stmt.Exec(
		u.FirstName,
		u.LastName,
		u.Email,
		u.Gender,
		u.Age,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	// row.Scan()
	id, err := row.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = int(id)
	defer db.Close()
	return nil
}

func (u *User) Update(db *sql.DB, id string) error {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("update users set first_name =?,last_name=? ,email=? ,gender=? ,age=? where id=?")
	_, err := stmt.Exec(
		u.FirstName,
		u.LastName,
		u.Email,
		u.Gender,
		u.Age,
		id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	defer db.Close()
	return nil
}

func (u *User) Delete(db *sql.DB, id string) error {
	row := db.QueryRow("select * from users where id = ?", id)
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Gender, &u.Age)
	if err != nil {
		return err
	}
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("delete from users where id=?")
	_, err = stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	defer db.Close()
	return nil
}
