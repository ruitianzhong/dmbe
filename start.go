package main

import (
	"dmbe/api"
	"dmbe/authentication"
	"dmbe/config"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var GlobalConfig *config.Config

func main() {
	GlobalConfig = config.InitConfig()
	api.SqlInit(GlobalConfig.Db.Address, GlobalConfig.Db.Port,
		GlobalConfig.Db.DbName, GlobalConfig.Db.Username,
		GlobalConfig.Db.Password)
	r := mux.NewRouter()
	r.HandleFunc("/", api.AddDrivers)
	r.HandleFunc("/api/fleets/get-all-fleets", api.GetAllFleets).Methods("get")
	r.HandleFunc("/api/drivers/add-drivers", api.AddDrivers).Methods("post")
	r.HandleFunc("/api/line/get-all-stops", api.GetAllStops).Methods("get")
	r.HandleFunc("/api/line/add-stop", api.AddStop).Methods("post")
	r.HandleFunc("/api/line/get-stops-by-line-id", api.GetStopsByLineId).Methods("get")
	r.HandleFunc("/api/line/get-all-line-info", api.GetAllLineInfo).Methods("get")
	r.HandleFunc("/api/line/add-new-line", api.AddNewLine).Methods("post")
	r.HandleFunc("/api/violation/types", api.GetAllViolationTypes).Methods("get")
	r.Use(authentication.AuthMiddleware)
	err := http.ListenAndServe(":"+GlobalConfig.App.Port, r)
	if err != nil {
		log.Fatal(err.Error())
	}
}
