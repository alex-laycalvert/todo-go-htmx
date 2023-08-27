package todos

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const DATABASE_URL = "db.sqlite"

type Todo struct {
	Id          uuid.UUID
	CreatedAt   time.Time
	Description string
}

func Init() error {
	return createTable()
}

func Todos() ([]Todo, error) {
	db, err := sql.Open("sqlite3", DATABASE_URL)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var todos []Todo
	rows, err := db.Query(`
	SELECT * FROM "Todo" ORDER BY "createdAt" DESC;
	`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var todo Todo
		err := rows.Scan(
			&todo.Id,
			&todo.CreatedAt,
			&todo.Description,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func Add(description string) error {
	if len(description) == 0 {
		return nil
	}
	db, err := sql.Open("sqlite3", DATABASE_URL)
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
	INSERT INTO "Todo" (
		"id",
		"description"
	) VALUES (?, ?);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	id := uuid.New().String()
	_, err = stmt.Exec(id, description)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func Remove(id string) error {
	if len(id) == 0 {
		return nil
	}
	db, err := sql.Open("sqlite3", DATABASE_URL)
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
	DELETE FROM "Todo" WHERE "id" = ?;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func createTable() error {
	db, err := sql.Open("sqlite3", DATABASE_URL)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "Todo" (
		"id" 			TEXT NOT NULL PRIMARY KEY,
		"createdAt" 	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		"description" 	TEXT NOT NULL
	);
	`)
	return err
}
