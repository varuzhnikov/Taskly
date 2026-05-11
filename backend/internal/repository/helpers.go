package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// isUniqueViolation reports whether err is a PostgreSQL unique_violation (SQLSTATE 23505).
// Uses errors.As so it correctly unwraps layered errors from pgx.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
