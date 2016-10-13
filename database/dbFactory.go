package database

import (
	"database/sql"
	"fmt"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

type DbFactory struct {
	Host string
	Port int
	Username string
	Password string
	DatabaseName string
}

func (db *DbFactory)NewConnect() (*sql.DB){
	dbScheme := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db.Username, db.Password, db.Host, db.Port, db.DatabaseName)
	connection, err := sql.Open("mysql", dbScheme)
	if err != nil {
		return nil
	}
	return connection
}

func (db *DbFactory) Check() error {
	dbScheme := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db.Username, db.Password, db.Host, db.Port, db.DatabaseName)
		sqlDb, err := sql.Open("mysql",dbScheme)
		if err != nil {
			sqlDb.Close()
		}
		return err
}
