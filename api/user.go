/*
 * Copyright (c) 2023/12/5.  Ruitian Zhong
 */

package api

import (
	"log"
	"net/http"
)

type UserInfo struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// GetUserInfoByCookie /api/user/info
func GetUserInfoByCookie(w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStore.Get(r, "dm-session")
	auth, ok1 := session.Values["authenticated"].(bool)
	username, ok2 := session.Values["username"].(string)
	if !ok1 || !ok2 || !auth {
		log.Println("reject get user info")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	db := DB
	info := UserInfo{}
	s := `SELECT name,driver_id from driver where driver_id=?`
	err := db.QueryRow(s, username).Scan(&info.Name, &info.Id)
	if err != nil {
		HandleError(err, w, http.StatusBadRequest)
		return
	}
	WriteJson(w, info)
}

type UpdatePasswordForm struct {
	NewPassword       string `schema:"new_password,required"`
	ConfirmedPassword string `schema:"confirmed_password,required"`
	OriginalPassword  string `schema:"original_password,required"`
}

// UpdateUserPassword /api/usr/update-user-password
func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	upf := UpdatePasswordForm{}
	if DecodePostForm(&upf, r, w) {
		return
	}
	var msg ResponseMsg
	if upf.NewPassword != upf.ConfirmedPassword {
		msg.Code = "100"
		msg.Msg = "确认密码和新密码不一致"
		WriteJson(w, msg)
		return
	} else if len(upf.NewPassword) < 5 {
		msg.Code = "100"
		msg.Msg = "密码必须不少于5位"
		WriteJson(w, msg)
		return
	}

	session, _ := cookieStore.Get(r, "dm-session")
	auth, ok1 := session.Values["authenticated"].(bool)
	username, ok2 := session.Values["username"].(string)
	if !ok1 || !ok2 || !auth {
		log.Println("reject get user info")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	db := DB
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			_ = tx.Rollback()
		}
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	s1 := `SELECT passwd from driver where driver_id=?`
	s2 := `UPDATE driver set passwd=? where driver_id=?`
	var passwd string

	if err = tx.QueryRow(s1, username).Scan(&passwd); err != nil {
		HandleError(err, w, http.StatusUnauthorized)
		_ = tx.Rollback()
		return
	}
	if passwd != upf.OriginalPassword {
		msg.Code = "100"
		msg.Msg = "原密码错误"
		_ = tx.Rollback()
		WriteJson(w, msg)
		return
	}

	_, err = tx.Exec(s2, upf.NewPassword, username)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		_ = tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		msg.Code = "100"
		msg.Msg = "Failed to commit"
	} else {
		msg.Code = "200"
	}
	WriteJson(w, msg)

}

// CheckIfLogin /api/user/check-if-login
func CheckIfLogin(w http.ResponseWriter, r *http.Request) {
	msg := ResponseMsg{}
	session, _ := cookieStore.Get(r, "dm-session")
	auth, ok := session.Values["authenticated"].(bool)
	if ok && auth {
		msg.Code = "200"
	} else {
		msg.Code = "100"
	}

	WriteJson(w, msg)
}
