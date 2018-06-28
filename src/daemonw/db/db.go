package db

import (
	"daemonw/conf"
	"fmt"
	"daemonw/log"
	"github.com/jmoiron/sqlx"
)

const (
	DB_CONN_WITHOUT_PASSWORD = "postgres://%s@%s:%d/%s?sslmode=%s"
	DB_CONN_WITH_PASSWORD    = "postgres://%s:%s@%s:%d/%s?sslmode=%s"
)

var (
	db *sqlx.DB
)

func init() {
	var err error
	c := &conf.Config.Database
	//connStr := "postgres://postgres:a123456@localhost:5432/mydb?sslmode=disable"
	var connParams string
	if c.Password == "" {
		connParams = fmt.Sprintf(DB_CONN_WITHOUT_PASSWORD, c.User, c.Host, c.Port, c.Name, c.SSLMode)
	} else {
		connParams = fmt.Sprintf(DB_CONN_WITH_PASSWORD, c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
	}
	db, err = sqlx.Connect("postgres", connParams)
	if err != nil {
		log.Fatal().Err(err).Msg("connect to database failed")
	}

	//create table which not exist
	err = initTables()
	fatalIfErr(err)
}

func initTables() error {
	tx, err := db.Beginx()
	fatalIfErr(err)
	defer tx.Rollback()

	//init user table
	if !existTable("users") {
		_, err = tx.Exec(SCHEMA_CREATE_USER_TABLE)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}

func GetDB() *sqlx.DB {
	return db
}

func existTable(name string) bool {
	rowNum := 0
	err := db.Get(&rowNum, `select count(*) from pg_class where relname = $1;`, name)
	fatalIfErr(err)
	return rowNum > 0
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
