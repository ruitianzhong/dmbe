package api

import (
	"database/sql"
	"net/http"
)

type BusInfo struct {
	BusId   string `json:"bus_id" schema:"bus_id,required"`
	LineId  string `json:"line_id" schema:"line_id,required"`
	FleetId string `json:"fleet_id" schema:"fleet_id"`
}

type AllBusInfo struct {
	BusInfo []BusInfo `json:"bus_info"`
	Code    string    `json:"code"`
}

// GetAllBus /api/bus/get-all-bus
func GetAllBus(w http.ResponseWriter, r *http.Request) {
	db := DB
	s := `SELECT bus.bus_id,bus.line_id,line.fleet_id from bus inner join line on line.line_id=bus.line_id `
	rows, err := db.Query(s)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var busInfo AllBusInfo
	var info BusInfo
	for rows.Next() {
		err = rows.Scan(&info.BusId, &info.LineId, &info.FleetId)
		if err != nil {
			HandleError(err, w, http.StatusInternalServerError)
			return
		}
		busInfo.BusInfo = append(busInfo.BusInfo, info)
	}
	WriteJson(w, busInfo)
}

// AddOneBus /api/bus/add-one-bus
func AddOneBus(w http.ResponseWriter, r *http.Request) {
	var info BusInfo
	if DecodePostForm(&info, r, w) {
		return
	}
	db := DB
	s := `INSERT INTO bus (bus_id, line_id) VALUES (?,?)`
	_, err := db.Exec(s, info.BusId, info.LineId)
	if err != nil {
		return
	}
	m := ResponseMsg{}
	if err != nil {
		m.Msg = "插入失败,请检查车辆是否已经存在"
		m.Code = "100"
	} else {
		m.Code = "200"
	}
	WriteJson(w, m)
}
