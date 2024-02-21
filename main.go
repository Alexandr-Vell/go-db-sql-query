package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Client struct {
	ID       int
	FIO      string
	Login    string
	Birthday string
	Email    string
}

func (c Client) String() string {
	return fmt.Sprintf("ID: %d FIO: %s Login: %s Birthday: %s Email: %s",
		c.ID, c.FIO, c.Login, c.Birthday, c.Email)
}

func main() {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	newClient := Client{
		FIO:      "Иван Иванович",
		Login:    "examplelogin",
		Birthday: "19860426",
		Email:    "example-email@mail.ru",
	}

	id, err := insertClient(db, newClient)
	if err != nil {
		fmt.Println(err)
		return
	}

	client, err := selectClient(db, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(client)

	newLogin := "newLoginExample"
	err = updateClientLogin(db, newLogin, id)
	if err != nil {
		fmt.Println(err)
		return
	}

	client, err = selectClient(db, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(client)

	err = deleteClient(db, id)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = selectClient(db, id)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func insertClient(db *sql.DB, client Client) (int64, error) {
	res, err := db.Exec("INSERT INTO clients (FIO, Login, Birthday, Email) VALUES (:FIO, :Login, :Birthday, :Email)",
		sql.Named("FIO", client.FIO),
		sql.Named("Login", client.Login),
		sql.Named("Birthday", client.Birthday),
		sql.Named("Email", client.Email))

	if err != nil {
		fmt.Println(err)
		return 0, nil
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}

	return id, nil
}

func updateClientLogin(db *sql.DB, login string, id int64) error {
	_, err := db.Exec("UPDATE clients SET login = :login WHERE id = :id",
		sql.Named("login", login),
		sql.Named("id", id))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func deleteClient(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM clients WHERE id = :id", sql.Named("id", id))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func selectClient(db *sql.DB, id int64) (Client, error) {
	client := Client{}

	row := db.QueryRow("SELECT id, fio, login, birthday, email FROM clients WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&client.ID, &client.FIO, &client.Login, &client.Birthday, &client.Email)

	return client, err
}
