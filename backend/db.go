package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// -----------------------------
// MySQL 接続
// -----------------------------
func ConnectMySQL(user, pass, host string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", user, pass, host)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// -----------------------------
// DB一覧取得
// -----------------------------
func GetDatabaseList(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW DATABASES;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		list = append(list, name)
	}
	return list, nil
}
