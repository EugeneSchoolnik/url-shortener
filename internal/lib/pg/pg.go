package pg

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func ParsePGError(err error) *pgconn.PgError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr
	}
	return nil
}
