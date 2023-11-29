package api

var SqlConnectionPath string

type ResponseMsg struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func SqlInit(address, port, dbName, username, password string) {

	SqlConnectionPath = username + ":" + password + "@(" + address + ":" + port + ")/" + dbName + "?parseTime=true"
}
