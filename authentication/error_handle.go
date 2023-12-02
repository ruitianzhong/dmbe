/*
 * Copyright (c) 2023/12/2.  Ruitian Zhong
 */

package authentication

import (
	"log"
	"net/http"
)

func HandleError(e error, w http.ResponseWriter, statusCode int) {
	log.Println(e.Error())
	w.WriteHeader(statusCode)
}
