package orderlimit

import (
	"context"
	"fmt"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/gogo/protobuf/proto"
	"storj.io/storj/pkg/pb"
	"storj.io/storj/pkg/storj"
)

// RemoveUnsentOrdersWithInvalidLimit removes all the invalid unsent orders
// which have an invalid order limit and it returns the number of deleted
// orders.
//
// An invalid order limit is one that cannot be deserialized.
//
// It panics if ctx or conn is nil.
//
// TODO: report the number of deleted orders by validation error type.
//
// TODO: have a force flag which cause to always delete any order which
// validation has failed despite that it could be caused by the unmarshaller
// mechanism.
func RemoveUnsentOrdersWithInvalidLimit(ctx context.Context, conn *sqlite3.Conn) (numDeleted uint64, _ error) {
	err := conn.WithTx(func() (err error) {
		delStmt, err := conn.Prepare(`DELETE FROM unsent_order WHERE satellite_id = ? AND serial_number = ?`)
		if err != nil {
			return err
		}

		defer func() {
			if err != nil {
				_ = delStmt.Close()
				return
			}

			err = delStmt.Close()
		}()

		selStmt, err := conn.Prepare(
			`SELECT satellite_id, serial_number, order_limit_serialized
			FROM unsent_order`,
		)
		if err != nil {
			return err
		}

		defer func() {
			if err != nil {
				_ = selStmt.Close()
				return
			}

			err = selStmt.Close()
		}()

		for {
			var hasRow bool
			hasRow, err = selStmt.Step()
			if err != nil {
				return err
			}

			if !hasRow {
				break
			}

			var (
				satID       []byte
				serialNum   []byte
				sOrderLimit []byte
			)
			err = selStmt.Scan(&satID, &serialNum, &sOrderLimit)
			if err != nil {
				return err
			}

			if err := validateSerializedOrderLimit(sOrderLimit); err == nil {
				continue
			} else {
				fmt.Printf("%+v\n", err)
			}

			if err = delStmt.Exec(satID, serialNum); err != nil {
				return err
			}

			if _, err = delStmt.Step(); err != nil {
				return err
			}

			numDeleted++
		}

		return nil
	})

	return numDeleted, err
}

func validateSerializedOrderLimit(sOrderLimit []byte) error {
	var orderLimit pb.OrderLimit
	err := proto.Unmarshal(sOrderLimit, &orderLimit)
	if err == nil {
		return nil
	}

	if storj.ErrSerialNumber.Has(err) ||
		storj.ErrNodeID.Has(err) ||
		storj.ErrPieceKey.Has(err) ||
		storj.ErrPieceID.Has(err) {
		return err
	}

	// This should be an error in the marshaller but we cannot ensure that
	// has been a temporary error so it's safer not to report a validation
	// error
	// TODO: Have a sentinel error for this case so the caller can decide
	// what to do
	return nil
}
