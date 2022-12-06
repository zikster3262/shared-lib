package source

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
	table = "sources"
)

type Source struct {
	MangaURL       string `json:"mangaurl"`
	HomePattern    string `json:"homepattern"`
	PagePattern    string `json:"pagepattern"`
	ChapterPattern string `json:"chapterpattern"`
	Append         bool   `json:"append"`
}

type SQL struct {
	ID             int64        `db:"id"`
	MangaURL       string       `db:"manga_url"`
	HomePattern    string       `db:"homepattern"`
	PagePattern    string       `db:"pagepattern"`
	ChapterPattern string       `json:"chapterpattern"`
	DateAdded      sql.NullTime `db:"dateadded"`
	Append         bool         `db:"append"`
}

// GetAllSources func return all rows in SourceSQL array from page table.
func GetAllSources(db *sqlx.DB) ([]SQL, error) {
	var sources []SQL
	err := db.Select(&sources, fmt.Sprintf("SELECT * FROM %s;", table))

	if err != nil {
		utils.FailOnError("db", err)
	}

	return sources, errors.Unwrap(err)
}

// Selects Source ID based on ID param.
func GetSourceID(db *sqlx.DB, id int64) SQL {
	var result SQL

	mx.Lock()
	err := db.Get(&result, fmt.Sprintf("SELECT * FROM "+table+" WHERE id = %d;", id))

	if err != nil {
		utils.LogWithInfo("db", "record does not exists in the database")
	}

	mx.Unlock()

	return result
}

var InsertSourceQuery = "INSERT INTO " + table + "(manga_url, home_pattern, page_pattern, append, chapter_pattern) VALUES (:manga_url, :home_pattern, :page_pattern, :append, :chapter_pattern);"

// InsertSource inserts interface input into source database table with sqlx DB struct
// Returns internal DB error on err
func InsertSource(db *sqlx.DB, m interface{}) error {
	_, err := db.NamedExec(InsertSourceQuery, m)

	if err != nil {
		utils.FailOnError("db", ErrDBInternalError)
	}

	return errors.Unwrap(err)
}

// GetSourcePage function takes sqlx DB struct and parameter string and returns SourceSQL.
func GetSourcePage(db *sqlx.DB, p string) SQL {
	var res SQL
	err := db.Get(&res, fmt.Sprintf("SELECT * FROM "+table+" WHERE manga_url = \"%v\"", p))

	if err != nil {
		utils.LogWithInfo("db", "record does not exists in the database")
	}

	return res
}
