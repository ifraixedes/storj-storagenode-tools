package internal

import (
	"context"
	"fmt"

	"go.fraixed.es/storj/sntools/dbcleaner/orderlimit"
	"go.fraixed.es/storj/sntools/internal/sqlite"
)

// UnsentInvalidLimits peforms the operations for the unsent-invalid-limits
// command.
func UnsentInvalidLimits(dbfile string, debug bool) (err error) {
	conn, err := sqlite.Open(dbfile)
	if err != nil {
		return err
	}
	defer closeDBConn(conn, debug)

	// TODO: use a context with a timeout
	ctx := context.Background()
	numDeleted, err := orderlimit.RemoveUnsentOrdersWithInvalidLimit(ctx, conn)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted %d unsert orders\n", numDeleted)
	return nil
}
