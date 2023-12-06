package api

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type ResponseMsg struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

var DB *sql.DB

func SqlInit(db *sql.DB) {
	DB = db
}
