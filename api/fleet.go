package api

import (
	"net/http"
)
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type AllFleetId struct {
	FleetIds []string `json:"fleet_ids"`
}

func GetAllFleets(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", SqlConnectionPath)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	query := `SELECT fleet_id from fleet`
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
	fleetIds := AllFleetId{}
	var fleetId string
	for rows.Next() {
		err := rows.Scan(&fleetId)
		if err != nil {
			HandleError(err, w, http.StatusInternalServerError)
			return
		}
		fleetIds.FleetIds = append(fleetIds.FleetIds, fleetId)
	}
	WriteJson(w, fleetIds)
}
