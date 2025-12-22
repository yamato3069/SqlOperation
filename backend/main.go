package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

// /connect ハンドラ
func handleConnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	var req ConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON decode error: "+err.Error(), http.StatusBadRequest)
		return
	}

	db, err := ConnectMySQL(req.User, req.Pass, req.Host)
	if err != nil {
		http.Error(w, "DB接続エラー:"+err.Error(), http.StatusInternalServerError)
		return
	}

	// グローバルに保持
	currentDB = db

	dbList, err := GetDatabaseList(db)
	if err != nil {
		http.Error(w, "DB一覧取得エラー:"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"databases": dbList,
	})
}

// /nl_query ハンドラ
func handleNLQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	if currentDB == nil {
		http.Error(w, "DB未接続です", http.StatusBadRequest)
		return
	}

	var req NLQueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON decode error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.DBName == "" {
		http.Error(w, "dbName が指定されていません", http.StatusBadRequest)
		return
	}

	// フロントから受け取った DB 名を使用
	schema, err := GetFullSchema(currentDB, req.DBName)
	if err != nil {
		http.Error(w, "スキーマ取得エラー:"+err.Error(), http.StatusInternalServerError)
		return
	}

	prompt := fmt.Sprintf(`
You are an SQL generator.
Below is the schema:

%s

User request:
%s

Write only ONE SQL SELECT statement.
No explanation.
`, schema, req.Query)

	sqlQuery, err := generateSQLWithPrompt(prompt)
	if err != nil {
		http.Error(w, "SQL生成エラー:"+err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := currentDB.Query(sqlQuery)
	if err != nil {
		http.Error(w, "SQL実行エラー:"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	var results []map[string]interface{}

	for rows.Next() {
		columnPointers := make([]interface{}, len(cols))
		columns := make([]interface{}, len(cols))

		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			http.Error(w, "データ読み込みエラー:"+err.Error(), http.StatusInternalServerError)
			return
		}

		row := map[string]interface{}{}
		for i, col := range cols {
			v := columnPointers[i].(*interface{})

			// ★ ここで []byte → string に変換することで base64 を防ぐ
			switch val := (*v).(type) {
			case []byte:
				row[col] = string(val)
			default:
				row[col] = val
			}
		}

		results = append(results, row)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sql":  sqlQuery,
		"data": results,
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/connect", handleConnect)
	mux.HandleFunc("/nl_query", handleNLQuery)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	fmt.Println("API Server started : http://localhost:8080")
	http.ListenAndServe(":8080", c.Handler(mux))
}
