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

type DriverInfo struct {
	DriverId string `json:"driver_id"`
	Name     string `json:"name"`
	Year     int    `json:"year"`
	FleetId  string `json:"fleet_id"`
	LineId   string `json:"line_id"`
	Gender   string `json:"gender"`
	Position string `json:"position"`
}

type AllDriverInfo struct {
	Code       string       `json:"code"`
	DriverInfo []DriverInfo `json:"driver_info"`
}

// GetAllDriverInfo /api/driver/get-all-driver-info
func GetAllDriverInfo(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(DriverName, SqlConnectionPath)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	s1 := `SELECT driver.driver_id,driver.name,driver.sex,driver.fleet_id,driver.position,driver.year,driver_line.line_id FROM driver left join driver_line on driver_line.driver_id=driver.driver_id`
	rows, err := db.Query(s1)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	info := DriverInfo{}
	all := AllDriverInfo{}
	for rows.Next() {
		var gender, position int
		var n_line_id sql.NullString
		err = rows.Scan(&info.DriverId, &info.Name, &gender, &info.FleetId, &position, &info.Year, &n_line_id)
		if err != nil {
			HandleError(err, w, http.StatusInternalServerError)
			return
		}
		if gender == 1 {
			info.Gender = "male"
		} else {
			info.Gender = "female"
		}
		switch position {
		case 0:
			info.Position = "普通司机"
			break
		case 1:
			info.Position = "路线队长"
			break
		}
		if n_line_id.Valid {
			info.LineId = n_line_id.String
		}
		all.DriverInfo = append(all.DriverInfo, info)
	}
	all.Code = "200"
	WriteJson(w, all)
}

type CaptainInfo struct {
	DriverId string `json:"driver_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
}

// GetFleetCaptainByDriverId /api/driver/get-fleet-captain-by-driver-id
// TODO: testing is needed
func GetFleetCaptainByDriverId(w http.ResponseWriter, r *http.Request) {
	ok, driverId := getDriverIdFromURL(w, r)
	if !ok {
		return
	}
	db, err := sql.Open(DriverName, SqlConnectionPath)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	info := CaptainInfo{}
	s1 := `SELECT driver_id,name from driver where fleet_id=(SELECT fleet_id from driver where driver_id=?) AND position=1`
	err = db.QueryRow(s1, driverId).Scan(&info.DriverId, &info.Name)
	if err != nil {
		HandleError(err, w, http.StatusOK)
		info.Code = "100"
	} else {
		info.Code = "200"
	}
	WriteJson(w, info)

}

// GetLineCaptainByDriverId  /api/driver/get-line-captain-by-driver-id
// TODO:Testing is needed
func GetLineCaptainByDriverId(w http.ResponseWriter, r *http.Request) {
	ok, driverId := getDriverIdFromURL(w, r)
	if !ok {
		return
	}
	db, err := sql.Open(DriverName, SqlConnectionPath)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	info := CaptainInfo{}
	s1 := `SELECT driver.driver_id,driver.name from driver 
        where driver.driver_id=(SELECT driver_line.driver_id from driver_line where driver_line.position=1 
        AND line_id=(SELECT line_id from driver_line where driver_id=?))`
	err = db.QueryRow(s1, driverId).Scan(&info.DriverId, &info.Name)
	if err != nil {
		HandleError(err, w, http.StatusOK)
		info.Code = "100"
	} else {
		info.Code = "200"
	}
	WriteJson(w, info)
}

func getDriverIdFromURL(w http.ResponseWriter, r *http.Request) (bool, string) {
	var driverId string
	if r.URL.Query().Has("driver_id") {
		driverId = r.URL.Query().Get("driver_id")
		return true, driverId
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return false, ""
	}
}

type ModifyDriverForm struct {
	Year     int    `schema:"year,required"`
	FleetId  string `schema:"fleet_id,required"`
	LineId   string `schema:"line_id,required"`
	Gender   string `schema:"gender,required"`
	DriverId string `schema:"driver_id,required"`
	Name     string `schema:"name,required"`
}

// ModifyDriverInfo /api/driver/modify-driver-info
func ModifyDriverInfo(w http.ResponseWriter, r *http.Request) {
	var mdf ModifyDriverForm
	if DecodePostForm(&mdf, r, w) {
		return
	}
	db, err := sql.Open(DriverName, SqlConnectionPath)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			_ = tx.Rollback()
		}
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	var gender int
	if mdf.Gender == "male" {
		gender = 1
	}
	msg := ResponseMsg{}
	s1 := `UPDATE driver set year=?,fleet_id=?,sex=?,name=? where driver_id=?`
	_, err = tx.Exec(s1, mdf.Year, mdf.FleetId, gender, mdf.Name, mdf.DriverId)
	if err != nil {
		msg.Code = "100"
		msg.Msg = err.Error()
		WriteJson(w, msg)
		_ = tx.Rollback()
		return
	}

	s2 := `SELECT fleet_id from line where line_id=? AND fleet_id=?`
	rows, err := tx.Query(s2, mdf.LineId, mdf.FleetId)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		_ = tx.Rollback()
		return
	} else if !rows.Next() {
		msg.Code = "100"
		msg.Msg = "路线和车队不一致"
		_ = rows.Close()
		_ = tx.Rollback()
		WriteJson(w, msg)
		return
	}
	_ = rows.Close()

	s5 := `SELECT driver_id from driver_line where driver_id=?`
	var temp string
	err = tx.QueryRow(s5, mdf.DriverId).Scan(&temp)

	if err == nil {
		s3 := `UPDATE driver_line set line_id=? where driver_id=?`
		_, err := tx.Exec(s3, mdf.LineId, mdf.DriverId)
		if err != nil {
			_ = tx.Rollback()
			HandleError(err, w, http.StatusInternalServerError)
			return
		}

		err = tx.Commit()
		if err != nil {
			msg.Code = "100"
			msg.Msg = "提交失败"
		} else {
			msg.Code = "200"

		}
		WriteJson(w, msg)
		return
	}
	s4 := `INSERT INTO driver_line (driver_id,line_id,position) values (?,?,?)`
	_, err = tx.Exec(s4, mdf.DriverId, mdf.LineId, 0)
	if err != nil {
		_ = tx.Rollback()
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		msg.Code = "100"
		msg.Msg = "Failed to commit " + err.Error()
	} else {
		msg.Code = "200"
	}
	WriteJson(w, msg)

}
