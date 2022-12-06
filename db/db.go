package db

import (
	"fmt"
	"os"
	"sync"

	"github.com/zikster3262/shared-lib/utils"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/jmoiron/sqlx"
)

var (
	sQLXConnection *sqlx.DB
	err            error
)

var onceSQLx sync.Once

func OpenSQLx() *sqlx.DB {
	for {
		onceSQLx.Do(func() {
			dsn := os.Getenv("DB_URL")

			utils.LogWithInfo("db", fmt.Sprintf("database connection: %v", dsn))

			sQLXConnection, err = sqlx.Connect("mysql", dsn)
			if err != nil {
				utils.FailOnError("db", err)
			}
		})

		err := sQLXConnection.Ping()
		if err == nil {
			break
		} else {
			continue
		}
	}

	utils.LogWithInfo("db", "connected to database")

	return sQLXConnection
}
