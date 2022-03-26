/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package graphql

import (
	"database/sql"
	_ "embed"
	"time"
)

var (
	//go:embed common.gql
	CommonTypes string
)

// Converts a sql nullstring to a string pointer.
func NullStr(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}
	return nil
}

// Format time as a string.
func FormatTime(input time.Time) *string {
	if input.IsZero() {
		return nil
	}
	val := input.Format(time.RFC3339)
	return &val
}
