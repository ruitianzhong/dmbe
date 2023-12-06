/*
 * Copyright (c) 2023/12/2.  Ruitian Zhong
 */

package authentication

import (
	"database/sql"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore
var (
	SqlConnectionPath string
	DriverName        string
	DB                *sql.DB
)

func InitAuthentication(sessionKey string) {
	store = sessions.NewCookieStore([]byte(sessionKey))
}

func SqlInit(db *sql.DB) {

	DB = db
}
