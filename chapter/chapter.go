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
	table              = "chapters"
	mx                 sync.Mutex
	ErrDBInternalError = errors.New("record was not created due to internal error")
)

type Chapter struct {
	Page_id int64  `json:"page_id"`
	Title   string `json:"title"`
	Url     string `json:"url"`
	Append  bool   `json:"append"`
}

type ChapterSQL struct {
	Page_id    int64        `db:"page_id"`
	Title      string       `db:"title"`
	Url        string       `db:"url"`
	Date_Added sql.NullTime `db:"date_added"`
	Append     bool         `db:"append"`
}

func CreateNewChapter(page_id int64, title, url string, append bool) Chapter {
	return Chapter{
		Page_id: page_id,
		Title:   title,
		Url:     url,
		Append:  append,
	}
}

var selectAllQuery = fmt.Sprintf("SELECT * FROM %s;", table)

func GetAllChapters(db *sqlx.DB) (p []ChapterSQL, err error) {
	err = db.Select(&p, selectAllQuery)
	if err != nil {
		utils.FailOnError("db", err)
	}
	return p, err

}

func (ch Chapter) InsertChapter(db *sqlx.DB) error {
	mx.Lock()
	_, err := db.NamedExec("INSERT INTO "+table+"(page_id, title, url, append) VALUES ((select id from db.pages WHERE id = :page_id), :title, :url, :append);", ch)
	if err != nil {
		utils.FailOnError("db", ErrDBInternalError)
	}

	mx.Unlock()
	return err

}

func GetChapter(db *sqlx.DB, p string) (ChapterSQL, bool, error) {

	var res ChapterSQL
	mx.Lock()
	err := db.Get(&res, fmt.Sprintf("SELECT * FROM "+table+" WHERE title = \"%v\"", p))
	mx.Unlock()

	if err != nil {
		return ChapterSQL{}, false, err
	}

	return res, true, err
}
