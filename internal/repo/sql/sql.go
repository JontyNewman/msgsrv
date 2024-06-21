package sql

import "database/sql"

type SqlMessageRepository struct {
	db              *sql.DB
	selectQuery     string
	insertStatement string
}

func InitSqliteMessageRepository(db *sql.DB) (*SqlMessageRepository, error) {

	initStatement := `
	CREATE TABLE IF NOT EXISTS message (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		body TEXT
	);`

	return InitSqlMessageRepository(db, initStatement)
}

func InitSqlMessageRepository(
	db *sql.DB,
	initStatement string,
	args ...any,
) (*SqlMessageRepository, error) {

	_, err := db.Exec(initStatement, args...)

	if err != nil {
		return nil, err
	}

	repo := NewSqlMessageRepository(
		db,
		"SELECT body FROM message WHERE id = ?;",
		"INSERT INTO message (body) VALUES (?);",
	)

	return &repo, nil
}

func NewSqlMessageRepository(db *sql.DB, selectQuery string, insertStatement string) SqlMessageRepository {
	return SqlMessageRepository{
		db:              db,
		selectQuery:     selectQuery,
		insertStatement: insertStatement,
	}
}

func (repo *SqlMessageRepository) Add(message string) (uint, error) {
	res, err := repo.db.Exec(repo.insertStatement, message)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return uint(id), err
}

func (repo *SqlMessageRepository) Fetch(id uint) (string, bool, error) {

	var message string
	err := repo.db.QueryRow(repo.selectQuery, id).Scan(&message)

	if err == sql.ErrNoRows {
		return "", false, nil
	}

	if err != nil {
		return "", false, err
	}

	return message, true, nil
}
