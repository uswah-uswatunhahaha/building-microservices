package db

import (
	"database/sql"
	"fmt"

	// go mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// CreateDatabase is a database setup
func CreateDatabase() (*sql.DB, error) {
	serverName := "localhost:3306"
	user := "root"
	password := "73niuzq@"
	dbName := "db_module"

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true", user, password, serverName, dbName)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
