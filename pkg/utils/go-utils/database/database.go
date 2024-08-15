package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
	Err    error
)

func PostgreSQLConnect(username, password, host, databaseName, port, sslMode, timeZone string) {
	DBConn, Err = gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, username, password, databaseName, port, sslMode, timeZone)), &gorm.Config{})
}
