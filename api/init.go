package api

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

type ResponseMsg struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

var DB *sql.DB
var cookieStore *sessions.CookieStore

func SqlInit(db *sql.DB) {
	DB = db
}

func InitCookieStore(store *sessions.CookieStore) {
	cookieStore = store
}
