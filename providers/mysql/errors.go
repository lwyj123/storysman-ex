package mysql

import "errors"

var (
	DBNotFoundError   = errors.New("Record Not Found")
	DBNoRowBeAffected = errors.New("No Row Be Affected In Strick Update Mode")
)
