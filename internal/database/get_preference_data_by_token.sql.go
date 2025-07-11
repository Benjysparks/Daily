// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: get_preference_data_by_token.sql

package database

import (
	"context"
)

const getPreferencesByToken = `-- name: GetPreferencesByToken :one
SELECT up.user_id, up.preferences, up.preference_variables
FROM users u
JOIN user_preferences up ON u.id = up.user_id
WHERE u.jwt_token = $1
`

func (q *Queries) GetPreferencesByToken(ctx context.Context, jwtToken string) (UserPreference, error) {
	row := q.db.QueryRowContext(ctx, getPreferencesByToken, jwtToken)
	var i UserPreference
	err := row.Scan(&i.UserID, &i.Preferences, &i.PreferenceVariables)
	return i, err
}
