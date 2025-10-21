package models

import (
	"io"
	"time"

	"fiber-boilerplate/internal/pkg/database"
	logging "fiber-boilerplate/internal/pkg/logging"
)

// EntityBlock :
type EntityBlock struct {
	Appuser   *AppuserBlock
	ArrayTest *ArrayTestBlock
}

// SQL :
var SQL *database.SQL

// Setup :
func Setup() {
	SQL = new(database.SQL)
	var lastErr error
LOOP:
	for trial := 0; trial < 5; trial++ {
		err := SQL.Connect(database.DriverPostgres)
		switch err {
		case nil:
			break LOOP
		case io.EOF:
			lastErr = err
			if trial == 4 {
				// Last attempt failed
				panic(err)
			}
			logging.Warn(err, "DB connection attempt %d/5 failed (io.EOF)", trial+1)
			time.Sleep(time.Second)
		default:
			panic(err)
		}
	}

	// if err exists, panic
	if lastErr != nil {
		panic(lastErr)
	}
}
