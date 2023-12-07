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
	DB *sql.DB
)

func SqlInit(db *sql.DB) {

	DB = db
}

func InitCookieStore(s *sessions.CookieStore) {
	store = s
}
