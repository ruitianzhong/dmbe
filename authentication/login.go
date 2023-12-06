/*
 * Copyright (c) 2023.  Ruitian Zhong
 */

package authentication

import (
	"fmt"
	"github.com/gorilla/schema"
	"net/http"
)

var decoder = schema.NewDecoder()

// LoginForm insecure for now
type LoginForm struct {
	Password string `schema:"password,required"`
	Username string `schema:"username,required"`
}
type LoginResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

// Login /auth/login
func Login(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		HandleError(err, w, http.StatusBadRequest)
		return
	}
	var lf LoginForm
	err = decoder.Decode(&lf, r.PostForm)
	if err != nil {
		HandleError(err, w, http.StatusBadRequest)
		return
	}
	db := DB
	s := "SELECT passwd from driver where driver_id=?"
	lr := LoginResponse{}
	var passwd string
	if err = db.QueryRow(s, lf.Username).Scan(&passwd); err != nil || passwd != lf.Password {
		lr.Code = "100"
		WriteJson(w, lr)
		return
	}
	session, _ := store.Get(r, "dm-session")
	session.Values["authenticated"] = true
	session.Values["username"] = lf.Username
	session.Values["auth_level"] = 0
	lr.Code = "200"
	if err = session.Save(r, w); err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	WriteJson(w, lr)
}

// Logout /auth/logout
func Logout(w http.ResponseWriter, r *http.Request) {
	l := LoginResponse{}
	session, _ := store.Get(r, "dm-session")
	session.Values["authenticated"] = false
	l.Code = "200"
	err := session.Save(r, w)
	WriteJson(w, l)
	if err != nil {
		fmt.Println(err)
		return
	}
}
