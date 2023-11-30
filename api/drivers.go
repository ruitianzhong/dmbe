package api

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/schema"
	"net/http"
)

var decoder = schema.NewDecoder()

type AddDriversForm struct {
	DriverId string `schema:"driver_id,required"`
	Gender   string `schema:"gender,required"`
	FleetId  string `schema:"fleet_id,required"`
	Year     int    `schema:"year,required"`
	Name     string `schema:"name,required"`
}

// AddDrivers
/*
path: /api/drivers/add-drivers
params: driver_id,gender,fleet_id,year,name
return: code,msg etc
*/
func AddDrivers(w http.ResponseWriter, r *http.Request) {
	var adf AddDriversForm
	if DecodePostForm(&adf, r, w) {
		return
	}
	db, err := sql.Open("mysql", SqlConnectionPath)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	gender := 0
	if adf.Gender == "male" {
		gender = 1
	}
	query := `INSERT INTO driver (driver_id, name, year,sex,fleet_id,position,passwd) VALUES (?,?,?,?,?,?,?)`
	_, err = db.Exec(query, adf.DriverId, adf.Name, adf.Year, gender, adf.FleetId, 0, "123456")
	// just temporary
	m := ResponseMsg{}
	if err != nil {
		m.Code = "100"
		m.Msg = err.Error()
	} else {
		m.Code = "200"
	}
	WriteJson(w, m)
}

func validateAddDriverForm(form AddDriversForm) (bool, string) {
	return true, ""
}
