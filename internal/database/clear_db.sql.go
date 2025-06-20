// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: clear_db.sql

package database

import (
	"context"
)

const clearUserTable = `-- name: ClearUserTable :exec

DELETE FROM users
`

func (q *Queries) ClearUserTable(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, clearUserTable)
	return err
}
