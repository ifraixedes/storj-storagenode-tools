package sqlite

import (
	"os"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/pkg/errors"
)

// Open returns a connection to the SQLite3 database file path.
// The database file must exists.
func Open(dbfpath string) (*sqlite3.Conn, error) {
	fi, err := os.Stat(dbfpath)
	if err != nil {
		return nil, err
	}

	if (fi.Mode() & os.ModeType) != 0 {
		return nil, errors.New("invalid DB file")
	}

	conn, err := sqlite3.Open(dbfpath, sqlite3.OPEN_READWRITE)
	if err != nil {
		return nil, errors.Wrap(err, "open connection error")
	}

	if err = conn.Exec("PRAGMA journal_mode=wal"); err != nil {
		return nil, errors.Wrap(err, "error when enabling WAL mode")
	}

	return conn, nil
}
