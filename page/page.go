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
	table              = "pages"
	mx                 sync.Mutex
	ErrDBInternalError = errors.New("record was not created due to internal error")
)

type Page struct {
	Title   string `json:"title"`
	Url     string `json:"url"`
	Page_Id int64  `json:"page_id"`
	Append  bool   `json:"append"`
}

type PageSQL struct {
	Id         int64        `db:"id"`
	Title      string       `db:"title"`
	Url        string       `db:"url"`
	Page_Id    int          `db:"page_id"`
	Date_Added sql.NullTime `db:"date_added"`
	Append     bool         `db:"append"`
}

var selectAllPagesQuery = "SELECT * FROM " + table + ";"

// GetAllPage func return all rows in PageSQL array from page table
func GetAllPages(db *sqlx.DB) (p []PageSQL, err error) {

	err = db.Select(&p, selectAllPagesQuery)
	if err != nil {
		utils.FailOnError("db", err)
	}
	return p, err

}

var selectPageQuery = "SELECT * FROM " + table + " WHERE id = %v;"

// Selects Page ID based on ID param
func GetPageID(db *sqlx.DB, id int64) (result PageSQL) {
	mx.Lock()
	err := db.Get(&result, fmt.Sprintf(selectPageQuery, id))
	if err != nil {
		utils.LogWithInfo("db", "record does not exists in the database")
	}

	mx.Unlock()
	return result
}

var InsertPageQuery = "INSERT INTO " + table + "(title, url, page_id, append) VALUES (:title, :url, (select id from db.sources WHERE id = :page_id), :append);"

// InsertPage inserts interface input into Page database table with sqlx DB struct
// Returns internal DB error on err
func InsertPage(db *sqlx.DB, m interface{}) error {
	mx.Lock()
	_, err := db.NamedExec(InsertPageQuery, m)
	if err != nil {
		utils.FailOnError("db", ErrDBInternalError)
	}
	mx.Unlock()
	return err
}

var GetPageQuery = "SELECT * FROM " + table + " WHERE title = \"%v\""

// GetPage function takes sqlx DB struct and parameter string and returns PageSQL
func GetPage(db *sqlx.DB, p string) (PageSQL, bool, error) {
	mx.Lock()
	var res PageSQL
	err := db.Get(&res, fmt.Sprintf(GetPageQuery, p))
	if err != nil {
		utils.LogWithInfo("db", "record does not exists in the database")
		return PageSQL{}, false, err
	}
	mx.Unlock()
	return res, true, err
}
