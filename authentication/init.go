/*
 * Copyright (c) 2023/12/2.  Ruitian Zhong
 */

package authentication

import "github.com/gorilla/sessions"

var store *sessions.CookieStore
var (
	SqlConnectionPath string
	DriverName        string
)

func InitAuthentication(sessionKey string) {
	store = sessions.NewCookieStore([]byte(sessionKey))
}

func SqlInit(address, port, dbName, username, password string) {

	SqlConnectionPath = username + ":" + password + "@(" + address + ":" + port + ")/" + dbName + "?parseTime=true"
	DriverName = "mysql"
}
