package api

import (
	"database/sql"
	"net/http"
)

type AllViolationTypes struct {
	ViolationTypes []string `json:"violation_types"`
	Code           string   `json:"code"`
}

// GetAllViolationTypes GetViolationType /api/violation/types
func GetAllViolationTypes(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(DriverName, SqlConnectionPath)
	if err != nil {
		HandleError(err, w, http.StatusInternalServerError)
		return
	}
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
