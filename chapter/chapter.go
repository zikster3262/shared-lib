package chapter

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
	table = "chapters"
)

type Chapter struct {
	PageID         int64  `json:"pageid"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	ChapterPattern string `json:"chapterpattern"`
	Append         bool   `json:"append"`
}

type SQL struct {
	PageID         int64        `db:"pageid"`
	Title          string       `db:"title"`
	URL            string       `db:"url"`
	ChapterPattern string       `db:"chapterpattern"`
	DateAdded      sql.NullTime `db:"dateadded"`
	Append         bool         `db:"append"`
}

func GetAllChapters(db *sqlx.DB) ([]SQL, error) {
	var chapters []SQL
	err := db.Select(&chapters, fmt.Sprintf("SELECT * FROM %s;", table))

	if err != nil {
		utils.FailOnError("db", err)
	}

	return chapters, errors.Unwrap(err)
}

func (ch Chapter) InsertChapter(db *sqlx.DB) error {
	mx.Lock()

	_, err := db.NamedExec("INSERT INTO "+table+"(pageid, title, url, chapterpattern, append) VALUES ((select id from db.pages WHERE id = :pageid), :title, :url, :chapterpattern, :append);", ch)

	if err != nil {
		utils.FailOnError("db", ErrDBInternalError)
	}

	mx.Unlock()

	return errors.Unwrap(err)
}

func GetChapter(db *sqlx.DB, p string) (SQL, bool, error) {
	var res SQL

	mx.Lock()
	err := db.Get(&res, fmt.Sprintf("SELECT * FROM "+table+" WHERE url = \"%v\"", p))
	mx.Unlock()

	if err != nil {
		return SQL{}, false, errors.Unwrap(err)
	}

	return res, true, errors.Unwrap(err)
}
