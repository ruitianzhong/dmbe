package main

import (
	"database/sql"
	"dmbe/api"
	"dmbe/authentication"
	"dmbe/config"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

var GlobalConfig *config.Config

func main() {
	GlobalConfig = config.InitConfig()
	initDb()
	initCookieStore(GlobalConfig.Auth.SessionKey)
	r := mux.NewRouter()
	r.HandleFunc("/api/fleets/get-all-fleets", api.GetAllFleets).Methods(http.MethodGet)
	r.HandleFunc("/api/drivers/add-drivers", api.AddDrivers).Methods(http.MethodPost)
	r.HandleFunc("/api/driver/get-all-driver-info", api.GetAllDriverInfo).Methods(http.MethodGet)
	r.HandleFunc("/api/line/get-all-stops", api.GetAllStops).Methods(http.MethodGet)
	r.HandleFunc("/api/line/add-stop", api.AddStop).Methods(http.MethodPost)
	r.HandleFunc("/api/line/get-stops-by-line-id", api.GetStopsAndBusByLineId).Methods(http.MethodGet)
	r.HandleFunc("/api/line/get-all-line-info", api.GetAllLineInfo).Methods(http.MethodGet)
	r.HandleFunc("/api/line/add-new-line", api.AddNewLine).Methods(http.MethodPost)
	r.HandleFunc("/api/violation/types", api.GetAllViolationTypes).Methods(http.MethodGet)
	r.HandleFunc("/api/bus/get-all-bus", api.GetAllBus).Methods(http.MethodGet)
	r.HandleFunc("/api/bus/add-one-bus", api.AddOneBus).Methods(http.MethodPost)
	r.HandleFunc("/api/driver/get-fleet-captain-by-driver-id", api.GetFleetCaptainByDriverId).Methods(http.MethodGet)
	r.HandleFunc("/api/driver/get-line-captain-by-driver-id", api.GetLineCaptainByDriverId).Methods(http.MethodGet)
	r.HandleFunc("/api/driver/modify-driver-info", api.ModifyDriverInfo).Methods(http.MethodPost)
	r.HandleFunc("/auth/login", authentication.Login).Methods(http.MethodPost)
	r.HandleFunc("/auth/logout", authentication.Logout).Methods(http.MethodPost)
	r.HandleFunc("/api/violation/add-violation", api.AddViolation).Methods(http.MethodPost)
	r.HandleFunc("/api/violation/violation-by-time-range-driver-id", api.ViolationByTimeRangeAndDriverID).Methods(http.MethodGet)
	r.HandleFunc("/api/violation/violation-stat-by-time-range-and-fleet-id", api.ViolationStatByTimeRange).Methods(http.MethodGet)
	r.HandleFunc("/api/line/get-line-captain-by-line-id", api.GetLineMembersByLineId).Methods(http.MethodGet)
	r.HandleFunc("/api/line/set-line-captain", api.SetLineCaptain).Methods(http.MethodPost)
	r.HandleFunc("/api/line/set-fleet-captain", api.SetFleetCaptain).Methods(http.MethodPost)
	r.HandleFunc("/api/fleet/get-fleet-members", api.GetFleetLineMembersByFleetId).Methods(http.MethodGet)
	r.HandleFunc("/api/line/get-line-by-fleet-id", api.GetLineByFleetId).Methods(http.MethodGet)
	r.HandleFunc("/api/fleet/get-all-fleet-detailed-info", api.GetAllFleetDetailedInfo).Methods(http.MethodGet)
	r.HandleFunc("/api/user/info", api.GetUserInfoByCookie).Methods(http.MethodGet)
	r.HandleFunc("/api/usr/update-user-password", api.UpdateUserPassword).Methods(http.MethodPost)
	r.HandleFunc("/api/user/check-if-login", api.CheckIfLogin).Methods(http.MethodGet)
	r.Use(authentication.AuthMiddleware)
	err := http.ListenAndServe(":"+GlobalConfig.App.Port, r)
	if err != nil {
		log.Fatal(err.Error())
	}
}
func initDb() {

	driverName, path := ConnectionDriverAndPath(GlobalConfig.Db.Address, GlobalConfig.Db.Port,
		GlobalConfig.Db.DbName, GlobalConfig.Db.Username,
		GlobalConfig.Db.Password)

	DB, err := sql.Open(driverName, path)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	api.SqlInit(DB)
	authentication.SqlInit(DB)
}

func ConnectionDriverAndPath(address, port, dbName, username, password string) (driverName string, connectionPath string) {
	sqlConnectionPath := username + ":" + password + "@(" + address + ":" + port + ")/" + dbName + "?parseTime=true"
	driverName = "mysql"
	return driverName, sqlConnectionPath

}

func initCookieStore(sessionKey string) {
	store := sessions.NewCookieStore([]byte(sessionKey))
	api.InitCookieStore(store)
	authentication.InitCookieStore(store)
}
