package page

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/zikster3262/shared-lib/utils"
)

var (
	mx                 sync.Mutex
	ErrDBInternalError = errors.New("record was not created due to internal error")
)

const (
	table = "pages"
)

type Page struct {
	Title          string `json:"title"`
	URL            string `json:"url"`
	SourceID       int64  `json:"sourceid"`
	PagePattern    string `json:"pagepattern"`
	ChapterPattern string `json:"chapterpattern"`
	Append         bool   `json:"append"`
}

type SQL struct {
	ID             int64        `db:"id"`
	Title          string       `db:"title"`
	URL            string       `db:"url"`
	SourceID       int          `db:"sourceid"`
	PagePattern    string       `db:"pagepattern"`
	ChapterPattern string       `json:"chapterpattern"`
	DateAdded      sql.NullTime `db:"dateadded"`
	Append         bool         `db:"append"`
}

// GetAllPage func return all rows in SQL array from page table.
func GetAllPages(db *sqlx.DB) ([]SQL, error) {
	var pages []SQL

	err := db.Select(&pages, fmt.Sprintf("SELECT * FROM %s;", table))

	if err != nil {
		utils.FailOnError("db", err)
	}

	return pages, errors.Unwrap(err)
}

// Selects Page ID based on ID param.
func GetPageID(db *sqlx.DB, id int64) SQL {
	var result SQL

	mx.Lock()
	err := db.Get(&result, fmt.Sprintf("SELECT * FROM "+table+" WHERE id = %v;", id))

	if err != nil {
		utils.LogWithInfo("db", "record does not exists in the database")
	}

	mx.Unlock()

	return result
}

// InsertPage inserts interface input into Page database table with sqlx DB struct
// Returns internal DB error on err.
func (p Page) InsertPage(db *sqlx.DB) error {
	mx.Lock()

	_, err := db.NamedExec("INSERT INTO "+table+"(title, url, source_id, append, chapter_pattern) VALUES (:title, :url, (select id from db.sources WHERE id = :source_id), :append, :chapter_pattern);", p)
	if err != nil {
		utils.FailOnError("db", ErrDBInternalError)
	}
	mx.Unlock()

	return errors.Unwrap(err)
}

// GetPage function takes sqlx DB struct and parameter string and returns SQL
func GetPage(db *sqlx.DB, p string) (SQL, bool, error) {
	var res SQL

	mx.Lock()
	err := db.Get(&res, fmt.Sprintf("SELECT * FROM "+table+" WHERE title = \"%v\"", p))
	mx.Unlock()

	if err != nil {
		return SQL{}, false, errors.Unwrap(err)
	}

	return res, true, errors.Unwrap(err)
}
