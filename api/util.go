package api

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, v any) (bool, []byte) {
	marshal, err := json.Marshal(v)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return true, nil
	}
	_, err = w.Write(marshal)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return true, nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return false, marshal
}

func DecodePostForm(dst interface{}, r *http.Request, w http.ResponseWriter) bool {
	err := r.ParseForm()
	if err != nil {
		HandleError(err, w, http.StatusBadRequest)
		return true
	}
	err = decoder.Decode(dst, r.PostForm)
	if err != nil {
		HandleError(err, w, http.StatusBadRequest)
		return true
	}
	return false
}

func initCookieStore() {

}
