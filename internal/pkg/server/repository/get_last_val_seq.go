package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// GetLastValSeq returns id of last record not at all, because sequence auto increment
func GetLastValSeq(ctx context.Context, tx pgx.Tx, nameTable pgx.Identifier) (uint64, error) {
	sanitizedNameTable := nameTable.Sanitize()
	SQLGetLastValSeq := fmt.Sprintf(`SELECT last_value FROM %s;`, sanitizedNameTable)
	seqRow := tx.QueryRow(ctx, SQLGetLastValSeq)

	var count uint64

	if err := seqRow.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
