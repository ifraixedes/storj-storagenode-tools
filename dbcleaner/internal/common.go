package internal

import (
	"log"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
)

func closeDBConn(conn *sqlite3.Conn, debug bool) {
	err := conn.Close()
	if debug && err != nil {
		log.Printf("Error closing DB connection. %v", err)
	}
}
