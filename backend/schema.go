package main

import (
	"database/sql"
	"fmt"
	"strings"
)

// -----------------------------
// 対象 DB の全テーブルのスキーマを取得
// -----------------------------
func GetFullSchema(db *sql.DB, dbName string) (string, error) {

	_, err := db.Exec("USE " + dbName)
	if err != nil {
		return "", err
	}

	rows, err := db.Query("SHOW TABLES;")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tbl string
		rows.Scan(&tbl)
		tables = append(tables, tbl)
	}

	var schema strings.Builder

	for _, tbl := range tables {
		schema.WriteString(fmt.Sprintf("TABLE %s:\n", tbl))

		colRows, err := db.Query("DESCRIBE " + tbl)
		if err != nil {
			continue
		}

		for colRows.Next() {
			var field, colType, null, key, def, extra string
			colRows.Scan(&field, &colType, &null, &key, &def, &extra)
			schema.WriteString(fmt.Sprintf("  %s %s\n", field, colType))
		}
		colRows.Close()

		schema.WriteString("\n")
	}

	return schema.String(), nil
}
