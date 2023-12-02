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
	initDb()
	authentication.InitAuthentication(GlobalConfig.Auth.SessionKey)
	r := mux.NewRouter()
	r.HandleFunc("/", api.AddDrivers)
	r.HandleFunc("/api/fleets/get-all-fleets", api.GetAllFleets).Methods("get")
	r.HandleFunc("/api/drivers/add-drivers", api.AddDrivers).Methods("post")
	r.HandleFunc("/api/driver/get-all-driver-info", api.GetAllDriverInfo).Methods("get")
	r.HandleFunc("/api/line/get-all-stops", api.GetAllStops).Methods("get")
	r.HandleFunc("/api/line/add-stop", api.AddStop).Methods("post")
	r.HandleFunc("/api/line/get-stops-by-line-id", api.GetStopsByLineId).Methods("get")
	r.HandleFunc("/api/line/get-all-line-info", api.GetAllLineInfo).Methods("get")
	r.HandleFunc("/api/line/add-new-line", api.AddNewLine).Methods("post")
	r.HandleFunc("/api/violation/types", api.GetAllViolationTypes).Methods("get")
	r.HandleFunc("/api/bus/get-all-bus", api.GetAllBus).Methods("get")
	r.HandleFunc("/api/bus/add-one-bus", api.AddOneBus).Methods("post")
	r.HandleFunc("/api/driver/get-fleet-captain-by-driver-id", api.GetFleetCaptainByDriverId).Methods("get")
	r.HandleFunc("/api/driver/get-line-captain-by-driver-id", api.GetLineCaptainByDriverId).Methods("get")
	r.HandleFunc("/api/driver/modify-driver-info", api.ModifyDriverInfo).Methods("post")
	r.HandleFunc("/auth/login", authentication.Login).Methods("post")
	r.HandleFunc("/auth/logout", authentication.Logout).Methods("post")
	r.HandleFunc("/api/violation/add-violation", api.AddViolation).Methods("post")
	r.HandleFunc("/api/violation/violation-by-time-range-driver-id", api.ViolationByTimeRangeAndDriverID).Methods("get")
	r.HandleFunc("/api/violation/violation-stat-by-time-range-and-fleet-id", api.ViolationStatByTimeRange).Methods("get")
	r.HandleFunc("/api/line/get-line-captain-by-line-id", api.GetLineMembersByLineId).Methods("get")
	r.HandleFunc("/api/line/set-line-captain", api.SetLineCaptain).Methods("post")
	r.HandleFunc("/api/line/set-fleet-captain", api.SetFleetCaptain).Methods("post")
	r.HandleFunc("/api/fleet/get-fleet-members", api.GetFleetLineMembersByFleetId).Methods("get")
	r.Use(authentication.AuthMiddleware)
	err := http.ListenAndServe(":"+GlobalConfig.App.Port, r)
	if err != nil {
		log.Fatal(err.Error())
	}
}
func initDb() {
	api.SqlInit(GlobalConfig.Db.Address, GlobalConfig.Db.Port,
		GlobalConfig.Db.DbName, GlobalConfig.Db.Username,
		GlobalConfig.Db.Password)
	authentication.SqlInit(GlobalConfig.Db.Address, GlobalConfig.Db.Port,
		GlobalConfig.Db.DbName, GlobalConfig.Db.Username,
		GlobalConfig.Db.Password)
}
