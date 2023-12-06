package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

type AllStopsId struct {
	Stops []string `json:"stop_ids"`
}

// GetAllStops /api/line/get-all-stops
func GetAllStops(w http.ResponseWriter, _ *http.Request) {

	db := DB
	query := `SELECT stop_id from stop`
	rows, err := db.Query(query)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			HandleError(err, w, http.StatusInternalServerError)
			return
		}
	}(rows)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	stopIds := AllStopsId{}
	var stopId string
	for rows.Next() {
		err := rows.Scan(&stopId)
		if err != nil {
			HandleError(err, w, http.StatusInternalServerError)
			return
		}
		stopIds.Stops = append(stopIds.Stops, stopId)
	}
	WriteJson(w, stopIds)
}

type AddStopForm struct {
	StopId string `schema:"stop_id,required"`
}

// AddStop /api/line/add-stop
func AddStop(w http.ResponseWriter, r *http.Request) {
	var asf AddStopForm
	if DecodePostForm(&asf, r, w) {
		return
	}
	db := DB
	s := `insert into stop (stop_id) values(?)`
	var msg ResponseMsg
	_, err := db.Exec(s, asf.StopId)
	if err != nil {
		msg.Code = "100"
		msg.Msg = err.Error()
	} else {
		msg.Code = "200"
	}
	WriteJson(w, msg)
}

type AllLinesInfo struct {
	AllInfo []LineInfo `json:"all_info"`
}

type LineInfo struct {
	LineId      string `json:"line_id"`
	LineCaptain string `json:"driver_id"`
	LineFleetId string `json:"fleet_id"`
}

// GetAllLineInfo /api/line/get-all-line-info
func GetAllLineInfo(w http.ResponseWriter, _ *http.Request) {
	s := `SELECT line.line_id,line.fleet_id,captain.driver_id FROM line left join (SELECT line_id,driver_id from driver_line where position=1) as captain on captain.line_id=line.line_id`
	db := DB
	rows, err := db.Query(s)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			HandleError(err, w, http.StatusInternalServerError)
			return
		}
	}(rows)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	var info LineInfo
	var all AllLinesInfo
	for rows.Next() {
		ns := sql.NullString{}
		err := rows.Scan(&info.LineId, &info.LineFleetId, &ns)
		if ns.Valid {
			info.LineCaptain = ns.String
		} else {
			info.LineCaptain = ""
		}
		if err != nil {
			HandleError(err, w, http.StatusInternalServerError)
			return
		}
		all.AllInfo = append(all.AllInfo, info)
	}
	WriteJson(w, all)
}

type StopByLineId struct {
	StopIds []string `json:"stop_ids"`
	BusIds  []string `json:"bus_ids"`
	Code    string   `json:"code"`
}

// GetStopsAndBusByLineId /api/line/get-stops-by-line-id
func GetStopsAndBusByLineId(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("line_id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	lineId := r.URL.Query().Get("line_id")
	db := DB
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			_ = tx.Rollback()
		}
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	s1 := `SELECT stop_id from line_stop where line_id=? ORDER BY  stop_order ASC `
	s2 := `SELECT bus_id from bus where line_id=?`
	rows, err := tx.Query(s1, lineId)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		_ = tx.Rollback()
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			HandleError(err, w, http.StatusInternalServerError)
			return
		}
	}(rows)
	var stopId string
	var stopIds StopByLineId
	for rows.Next() {
		err = rows.Scan(&stopId)
		if err != nil {
			HandleError(err, w, http.StatusInternalServerError)
			_ = tx.Rollback()
			return
		}
		stopIds.StopIds = append(stopIds.StopIds, stopId)
	}
	if err = rows.Err(); err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		_ = tx.Rollback()
		return
	}
	rows2, err := tx.Query(s2, lineId)
	if err != nil {
		_ = tx.Rollback()
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	defer func(rows2 *sql.Rows) {
		_ = rows2.Close()
	}(rows2)
	var busId string
	for rows2.Next() {
		if err = rows2.Scan(&busId); err != nil {
			_ = tx.Rollback()
			HandleError(err, w, http.StatusInternalServerError)
			return
		}
		stopIds.BusIds = append(stopIds.BusIds, busId)
	}
	err = tx.Commit()

	if len(stopIds.StopIds) == 0 || err != nil {
		stopIds.Code = "100"
	} else {
		stopIds.Code = "200"
	}
	WriteJson(w, stopIds)

}

type AddNewLineForm struct {
	LineId  string   `json:"line_id,required"`
	FleetId string   `json:"fleet_id,required"`
	StopIds []string `json:"stop_ids,required"`
}

// AddNewLine /api/line/add-new-line
/*
@method post
@param line_id,fleet_id,stop_ids
*/
func AddNewLine(w http.ResponseWriter, r *http.Request) {
	l := r.ContentLength
	body := make([]byte, l)
	_, _ = r.Body.Read(body)
	fmt.Println(string(body))
	var anf AddNewLineForm
	err := json.Unmarshal(body, &anf)
	if err != nil {
		HandleError(err, w, http.StatusBadRequest)
		return
	}
	msg := ResponseMsg{}
	if len(anf.StopIds) < 2 {
		msg.Code = "100"
		msg.Msg = "站点数量过少"
		WriteJson(w, msg)
		return
	}

	db := DB
	s1 := `INSERT INTO line (line_id,fleet_id) values (?,?)`
	s2 := `INSERT INTO line_stop (stop_id,line_id,stop_order) values(?,?,?)`
	tx, err := db.Begin()

	if err != nil {
		if tx != nil {
			_ = tx.Rollback()
		}
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	_, err = tx.Exec(s1, anf.LineId, anf.FleetId)
	if err != nil {
		msg.Code = "200"
		msg.Msg = "该路线已经存在，请保证路线名的唯一性"
		WriteJson(w, msg)
		_ = tx.Rollback()
		return
	}

	for i := 0; i < len(anf.StopIds); i++ {
		_, err = tx.Exec(s2, anf.StopIds[i], anf.LineId, i)
		if err != nil {
			msg.Code = "100"
			msg.Msg = "更新时发生错误:" + err.Error()
			_ = tx.Rollback()
			WriteJson(w, msg)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		msg.Code = "100"
		msg.Msg = "提交时发生错误，请重试:" + err.Error()
	} else {
		msg.Code = "200"
	}
	WriteJson(w, msg)
}

type LineMember struct {
	DriverId string `json:"driver_id"`
	Name     string `json:"name"`
}

type LineMemberReply struct {
	Members    []LineMember `json:"members"`
	Captain    LineMember   `json:"captain"`
	HasCaptain bool         `json:"has_captain"`
}

// GetLineMembersByLineId  /api/line/get-line-captain-by-line-id
func GetLineMembersByLineId(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if !query.Has("line_id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db := DB

	lineId := query.Get("line_id")
	s := "SELECT driver.name,driver_line.driver_id,driver_line.position from driver_line inner join driver on driver.driver_id=driver_line.driver_id where driver_line.line_id=?"
	rows, err := db.Query(s, lineId)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	member := LineMember{}
	reply := LineMemberReply{}
	var position int
	for rows.Next() {
		err = rows.Scan(&member.Name, &member.DriverId, &position)
		if err != nil {
			HandleError(err, w, http.StatusInternalServerError)
			return
		}
		reply.Members = append(reply.Members, member)
		if position == 1 {
			reply.Captain = member
			reply.HasCaptain = true
		}
	}
	WriteJson(w, reply)
}

type SetLineCaptainForm struct {
	LineId   string `schema:"line_id,required"`
	DriverId string `schema:"driver_id,required"`
}

// SetLineCaptain /api/line/set-line-captain
func SetLineCaptain(w http.ResponseWriter, r *http.Request) {
	scf := SetLineCaptainForm{}
	if DecodePostForm(&scf, r, w) {
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
	msg := ResponseMsg{Code: "200"}
	s1 := `UPDATE driver_line set position=0 where line_id=? AND position=1`
	s2 := `UPDATE driver_line set position=1 where line_id=? AND driver_id=?`
	s3 := `SELECT line_id from driver_line where driver_id=?`
	var actualLineId string
	err = tx.QueryRow(s3, scf.DriverId).Scan(&actualLineId)
	if err != nil || actualLineId != scf.LineId {
		msg.Code = "100"
		msg.Msg = "司机不在所选路线当中"
		_ = tx.Rollback()
		WriteJson(w, msg)
		return
	}

	_, err = tx.Exec(s1, scf.LineId)
	if err != nil {
		_ = tx.Rollback()
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	_, err = tx.Exec(s2, scf.LineId, scf.DriverId)
	if err != nil {
		_ = tx.Rollback()
		HandleError(err, w, http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		msg.Code = "200"
		msg.Msg = err.Error()
	}
	WriteJson(w, msg)
}

type QueriedLineInfo struct {
	LineId []string `json:"line_id"`
}

// GetLineByFleetId /api/line/get-line-by-fleet-id
func GetLineByFleetId(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("fleet_id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fleetId := r.URL.Query().Get("fleet_id")
	s := `SELECT line_id from line where fleet_id=?`
	db := DB
	rows, err := db.Query(s, fleetId)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	var lineId string
	var info QueriedLineInfo
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		err = rows.Scan(&lineId)
		if err != nil {
			HandleError(err, w, http.StatusInternalServerError)
			return
		}
		info.LineId = append(info.LineId, lineId)
	}

	WriteJson(w, info)

}
