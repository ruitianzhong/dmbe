package api

var SqlConnectionPath string

func SqlInit(address, port, dbName, username, password string) {

	SqlConnectionPath = username + ":" + password + "@(" + address + ":" + port + ")/" + dbName + "?parseTime=true"
}
