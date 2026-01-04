package main

import "database/sql"

// -----------------------------
// 構造体
// -----------------------------
type ConnectRequest struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Host string `json:"host"`
}

type NLQueryRequest struct {
	DBName string `json:"db"`
	Query  string `json:"query"`
}

// -----------------------------
// グローバル（現在接続中のDB）
// -----------------------------
var currentDB *sql.DB
