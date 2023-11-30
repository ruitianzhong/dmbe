package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type AllStopsId struct {
	Stops []string `json:"stop_ids"`
}

// GetAllStops /api/line/get-all-stops
func GetAllStops(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", SqlConnectionPath)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	body, err := json.Marshal(stopIds)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
}
