package db

import "fmt"

const (
	mysqldb = "mysql"
	mongo   = "mongo"
)

func getSupportedDbTypes() []string {
	return []string{mysqldb}
}

func isSupportedType(dbType string) bool {
	for _, d := range getSupportedDbTypes() {
		if dbType == d {
			return true
		}
	}
	return false
}

// ErrNoRecords is the error returned when no rows are found in DB
var ErrNoRecords = fmt.Errorf("No records found")
