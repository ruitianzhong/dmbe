/*
 * Copyright (c) 2023/12/2.  Ruitian Zhong
 */

package authentication

import "github.com/gorilla/sessions"

var store *sessions.CookieStore

func InitAuthentication(sessionKey string) {
	store = sessions.NewCookieStore([]byte(sessionKey))
}
