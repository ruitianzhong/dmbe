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
	db := DB
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

type SetFleetCaptainForm struct {
	FleetId  string `schema:"fleet_id"`
	DriverId string `schema:"driver_id"`
}
type FleetMember struct {
	Name     string `json:"name"`
	DriverId string `json:"driver_id"`
}

type FleetMemberReply struct {
	FleetMembers []FleetMember `json:"fleet_members"`
	Captain      FleetMember   `json:"captain"`
	HasCaptain   bool          `json:"has_captain"`
}

// SetFleetCaptain /api/line/set-fleet-captain
func SetFleetCaptain(w http.ResponseWriter, r *http.Request) {
	var scf SetFleetCaptainForm
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
	var actualFleetId string
	s1 := `UPDATE driver set position=0 where position=1 AND fleet_id=?`
	s2 := `UPDATE driver set position=1 where fleet_id=? AND driver_id=?`
	s3 := `SELECT fleet_id from driver where driver_id=?`
	msg := ResponseMsg{}
	err = tx.QueryRow(s3, scf.DriverId).Scan(&actualFleetId)
	if err != nil || actualFleetId != scf.FleetId {
		_ = tx.Rollback()
		msg.Code = "100"
		msg.Msg = "司机与车队不匹配"
		WriteJson(w, msg)
		return
	}
	_, err = tx.Exec(s1, scf.FleetId)
	if err != nil {
		msg.Code = "100"
		msg.Msg = err.Error()
		_ = tx.Rollback()
		WriteJson(w, msg)
		return
	}
	_, err = tx.Exec(s2, scf.FleetId, scf.DriverId)
	if err != nil {
		msg.Code = "100"
		msg.Msg = err.Error()
		_ = tx.Rollback()
		WriteJson(w, msg)
		return
	}
	err = tx.Commit()
	msg.Code = "200"
	if err != nil {
		msg.Code = "100"
		msg.Msg = "Failed to commit"
	}
	WriteJson(w, msg)
}

// GetFleetLineMembersByFleetId /api/fleet/get-fleet-members
func GetFleetLineMembersByFleetId(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("fleet_id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fleetId := r.URL.Query().Get("fleet_id")
	db := DB
	var (
		member   FleetMember
		reply    FleetMemberReply
		position int
	)

	s := `SELECT driver_id,name,position from driver where fleet_id=?`
	rows, err := db.Query(s, fleetId)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		err = rows.Scan(&member.DriverId, &member.Name, &position)
		if err != nil {
			HandleError(err, w, http.StatusInternalServerError)
			return
		}
		reply.FleetMembers = append(reply.FleetMembers, member)
		if position == 1 {
			reply.HasCaptain = true
			reply.Captain = member
		}
	}
	WriteJson(w, reply)
}
