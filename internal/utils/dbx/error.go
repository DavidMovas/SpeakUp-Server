package dbx

import (
	"errors"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func IsUniqueViolation(err error, name string) bool {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		return pgError.Code == pgerrcode.UniqueViolation && strings.Contains(pgError.ConstraintName, name)
	}

	return false
}

func IsForeignKeyViolation(err error, name string) bool {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		return pgError.Code == pgerrcode.ForeignKeyViolation && strings.Contains(pgError.ConstraintName, name)
	}

	return false
}

func IsNoRows(err error) bool {
	if err != nil && err.Error() == "no rows in result set" {
		return true
	}
	return errors.Is(err, pgx.ErrNoRows)
}

func NotValidEnumType(err error) bool {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		return pgError.Code == pgerrcode.InvalidTextRepresentation
	}
	return false
}
