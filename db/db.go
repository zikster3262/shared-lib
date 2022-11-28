package db

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/zikster3262/shared-lib/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	SQLXConnection *sqlx.DB
	err            error
)

type SQLxProvider struct{}

var onceSQLx sync.Once

func OpenSQLx() *sqlx.DB {

	for {

		onceSQLx.Do(func() {
			dsn := os.Getenv("DB_URL")

			utils.LogWithInfo("db", fmt.Sprintf("database connection: %v", dsn))

			SQLXConnection, err = sqlx.Connect("mysql", dsn)
			if err != nil {
				utils.FailOnError("db", err)
			}
		})

		err := SQLXConnection.Ping()
		if err == nil {
			break
		} else {
			time.Sleep(time.Second * 1)
			continue
		}

	}

	utils.LogWithInfo("db", "connected to database")
	return SQLXConnection
}
