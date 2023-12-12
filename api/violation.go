package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type AllViolationTypes struct {
	ViolationTypes []string `json:"violation_types"`
	Code           string   `json:"code"`
}

// GetAllViolationTypes GetViolationType /api/violation/types
func GetAllViolationTypes(w http.ResponseWriter, _ *http.Request) {
	db := DB
	s := `SELECT violation_type_id from violation_type`
	rows, err := db.Query(s)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	var t string
	var allTypes AllViolationTypes
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		err := rows.Scan(&t)
		if err != nil {
			HandleError(err, w, http.StatusInternalServerError)
			return
		}
		allTypes.ViolationTypes = append(allTypes.ViolationTypes, t)
	}
	allTypes.Code = "200"
	WriteJson(w, allTypes)
}

type AddViolationForm struct {
	ViolationTypeId string `schema:"violation_type_id,required"`
	DriverId        string `schema:"driver_id,required"`
	StopId          string `schema:"stop_id,required"`
	Time            int64  `schema:"time,required"` // second
	BusId           string `schema:"bus_id,required"`
}

// AddViolation /api/violation/add-violation
func AddViolation(w http.ResponseWriter, r *http.Request) {
	var avf AddViolationForm
	if DecodePostForm(&avf, r, w) {
		return
	}
	t := time.Unix(avf.Time, 0)
	now, earliest := time.Now(), time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local)
	msg := ResponseMsg{Code: "100"}
	if t.After(now) || t.Before(earliest) {
		msg.Msg = "时间错误,请检查时间是否设置正确"
		WriteJson(w, msg)
		return
	}
	db := DB
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			_ = tx.Rollback()
		}
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	var lineId string
	s1 := `SELECT line_id from bus where bus_id=?`
	s2 := `SELECT driver_line.driver_id from driver_line 
    		inner join bus on bus.line_id=driver_line.line_id 
			where driver_line.line_id=? AND bus.bus_id=? AND driver_line.driver_id=?`
	s3 := `SELECT line_id from line_stop where stop_id=? AND line_id=?`
	s4 := `INSERT INTO violation_record (violation_type_id,time,driver_id,bus_id,fleet_id,stop_id,line_id) values (?,?,?,?,?,?,?)`
	s5 := `SELECT fleet_id from driver where driver_id=?`
	err = tx.QueryRow(s1, avf.BusId).Scan(&lineId)
	if err != nil {
		_ = tx.Rollback()
		msg.Msg = "线路和车牌号出现不一致"
		WriteJson(w, msg)
		return
	}
	var temp string
	err = tx.QueryRow(s2, lineId, avf.BusId, avf.DriverId).Scan(&temp)
	if err != nil {
		_ = tx.Rollback()
		msg.Msg = "数据不一致"
		WriteJson(w, msg)
		return
	}
	fmt.Println(lineId)
	err = tx.QueryRow(s3, avf.StopId, lineId).Scan(&temp)
	if err != nil {
		_ = tx.Rollback()
		msg.Msg = "站点和路线数据不一致 " + err.Error()
		WriteJson(w, msg)
		return
	}
	var fleetId string
	err = tx.QueryRow(s5, avf.DriverId).Scan(&fleetId)
	if err != nil {
		_ = tx.Rollback()
		msg.Msg = err.Error()
		WriteJson(w, msg)
		return
	}
	_, err = tx.Exec(s4, avf.ViolationTypeId, avf.Time, avf.DriverId, avf.BusId, fleetId, avf.StopId, lineId)
	if err != nil {
		_ = tx.Rollback()
		msg.Msg = err.Error()
		WriteJson(w, msg)
		return
	}
	err = tx.Commit()
	if err != nil {
		msg.Msg = "Commit Failed " + err.Error()
	} else {
		msg.Code = "200"
	}
	WriteJson(w, msg)
}

type ViolationInfo struct {
	ViolationTypeId string `json:"violation_type_id"`
	Time            int64  `json:"time"`
	BusId           string `json:"bus_id"`
	StopId          string `json:"stop_id"`
	FleetId         string `json:"fleet_id"`
	LineId          string `json:"line_id"`
}

type ViolationReply struct {
	Violations []ViolationInfo `json:"violations"`
	Code       string          `json:"code"`
}

type ViolationStatInfo struct {
	ViolationTypeId string `json:"violation_type_id"`
	Count           string `json:"count"`
}

type ViolationStatReply struct {
	Code          string              `json:"code"`
	ViolationStat []ViolationStatInfo `json:"violation_stat"`
}

// ViolationByTimeRangeAndDriverID /api/violation/violation-by-time-range-driver-id
func ViolationByTimeRangeAndDriverID(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if !query.Has("start") || !query.Has("end") || !query.Has("driver_id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	start, err1 := strconv.ParseInt(query.Get("start"), 10, 64)
	end, err2 := strconv.ParseInt(query.Get("end"), 10, 64)
	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	driverId := query.Get("driver_id")
	db := DB
	s := `SELECT time,violation_type_id,bus_id,stop_id,fleet_id,line_id from violation_record where time>=? AND time<? AND driver_id=?`
	rows, err := db.Query(s, start, end, driverId)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var reply ViolationReply
	var violation ViolationInfo
	for rows.Next() {
		err = rows.Scan(&violation.Time, &violation.ViolationTypeId, &violation.BusId, &violation.StopId, &violation.FleetId, &violation.LineId)
		if err != nil {
			HandleError(err, w, http.StatusInternalServerError)
			return
		}
		reply.Violations = append(reply.Violations, violation)
	}
	reply.Code = "200"
	WriteJson(w, reply)
}

// ViolationStatByTimeRange /api/violation/violation-stat-by-time-range-and-fleet-id
func ViolationStatByTimeRange(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if !query.Has("fleet_id") || !query.Has("start") || !query.Has("end") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	start, err1 := strconv.ParseInt(query.Get("start"), 10, 64)
	end, err2 := strconv.ParseInt(query.Get("end"), 10, 64)
	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fleetId := query.Get("fleet_id")
	db := DB
	s := `SELECT count(*),violation_type_id from violation_record where time>=? and time<? AND fleet_id=? group by violation_type_id`
	rows, err := db.Query(s, start, end, fleetId)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var stat ViolationStatInfo
	var reply ViolationStatReply
	for rows.Next() {
		if err = rows.Scan(&stat.Count, &stat.ViolationTypeId); err != nil {
			HandleError(err, w, http.StatusInternalServerError)
			return
		}
		reply.ViolationStat = append(reply.ViolationStat, stat)
	}
	reply.Code = "200"
	WriteJson(w, reply)
}
