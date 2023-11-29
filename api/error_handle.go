package api

import (
	"log"
	"net/http"
)

func HandleError(e error, w http.ResponseWriter, statusCode int) {
	log.Println(e.Error())
	w.WriteHeader(statusCode)
}
