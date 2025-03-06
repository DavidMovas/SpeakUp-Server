package dbx

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func QueryBatch(b *pgx.Batch, sb squirrel.Sqlizer) error {
	query, args, err := sb.ToSql()
	if err != nil {
		return err
	}
	b.Queue(query, args...)
	return nil
}
