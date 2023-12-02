/*
 * Copyright (c) 2023/12/2.  Ruitian Zhong
 */

package authentication

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
