/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package rdb

import "database/sql"

// Entity that is referenced by a unique token which may change over time.
type TokenReference struct {
	Token string `gorm:"unique;not null;size:128"`
}

// Entity that has a name and description.
type NamedEntity struct {
	Name        sql.NullString `gorm:"size:128"`
	Description sql.NullString `gorm:"size:1024"`
}

// Entity that has branding information.
type BrandedEntity struct {
	ImageUrl        sql.NullString `gorm:"size:512"`
	Icon            sql.NullString `gorm:"size:128"`
	BackgroundColor sql.NullString `gorm:"size:32"`
	ForegroundColor sql.NullString `gorm:"size:32"`
	BorderColor     sql.NullString `gorm:"size:32"`
}

// Creates a sql.NullString from a string constant.
func NullStrOf(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  true,
	}
}
