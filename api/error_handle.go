package api

import (
	"log"
	"net/http"
)

func InternalServerError(e error, w http.ResponseWriter) {
	log.Println(e.Error())
	w.WriteHeader(http.StatusInternalServerError)
}
