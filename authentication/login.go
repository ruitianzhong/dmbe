/*
 * Copyright (c) 2023.  Ruitian Zhong
 */

package authentication

import "net/http"

// LoginForm insecure for now
type LoginForm struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
type LoginResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func Login(w http.ResponseWriter, r *http.Request) {

}
