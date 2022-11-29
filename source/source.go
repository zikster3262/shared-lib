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
	table              = "sources"
	mx                 sync.Mutex
	ErrDBInternalError = errors.New("record was not created due to internal error")
)

type Source struct {
	Manga_URL       string `json:"manga_url"`
	Home_Pattern    string `json:"home_pattern"`
	Page_Pattern    string `json:"page_pattern"`
	Chapter_Pattern string `json:"chapter_pattern"`
	Append          bool   `json:"append"`
}

type SourceSQL struct {
	Id              int64        `db:"id"`
	Manga_URL       string       `db:"manga_url"`
	Home_Pattern    string       `db:"home_pattern"`
	Page_Pattern    string       `db:"page_pattern"`
	Chapter_Pattern string       `json:"chapter_pattern"`
	Date_Added      sql.NullTime `db:"date_added"`
	Append          bool         `db:"append"`
}

var selectAllSourcesQuery = fmt.Sprintf("SELECT * FROM %s;", table)

// GetAllSources func return all rows in SourceSQL array from page table
func GetAllSources(db *sqlx.DB) (p []SourceSQL, err error) {
	err = db.Select(&p, selectAllSourcesQuery)
	if err != nil {
		utils.FailOnError("db", err)
	}
	return p, err

}

// Selects Source ID based on ID param
func GetSourceID(db *sqlx.DB, id int64) (result SourceSQL) {
	mx.Lock()
	err := db.Get(&result, fmt.Sprintf("SELECT * FROM "+table+" WHERE id = %d;", id))
	if err != nil {
		utils.LogWithInfo("db", "record does not exists in the database")
	}

	mx.Unlock()
	return result
}

var InsertSourceQuery = "INSERT INTO " + table + "(manga_url, home_pattern, page_pattern, append) VALUES (:manga_url, :home_pattern, :page_pattern, :append);"

// InsertSource inserts interface input into source database table with sqlx DB struct
// Returns internal DB error on err
func InsertSource(db *sqlx.DB, m interface{}) error {
	_, err := db.NamedExec(InsertSourceQuery, m)
	if err != nil {
		utils.FailOnError("db", ErrDBInternalError)
	}
	return err
}

var GetSourcePageQuery = "SELECT * FROM " + table + " WHERE manga_url = \"%v\""

// GetSourcePage function takes sqlx DB struct and parameter string and returns SourceSQL
func GetSourcePage(db *sqlx.DB, p string) (res SourceSQL) {
	err := db.Get(&res, fmt.Sprintf(GetSourcePageQuery, p))
	if err != nil {
		utils.LogWithInfo("db", "record does not exists in the database")
	}

	return res
}
