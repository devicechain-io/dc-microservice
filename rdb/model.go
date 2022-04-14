/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package rdb

import (
	"database/sql"
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

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

// Entity that has extra attached metadata.
type MetadataEntity struct {
	Metadata *datatypes.JSON
}

// Create JSON value from string input.
func MetadataStrOf(value *string) *datatypes.JSON {
	if value != nil {
		result := json.RawMessage{}
		err := result.UnmarshalJSON([]byte(*value))
		if err != nil {
			return nil
		}
		conv := datatypes.JSON(result)
		return &conv
	}
	return nil
}

// Creates a sql.NullString from a string constant.
func NullStrOf(value *string) sql.NullString {
	if value != nil {
		return sql.NullString{
			String: *value,
			Valid:  true,
		}
	} else {
		return sql.NullString{
			Valid: false,
		}
	}
}

// Information for paged result sets
type Pagination struct {
	PageNumber int32
	PageSize   int32
}

// Scope function used to implement pagination.
func Paginate(pag Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (pag.PageNumber - 1) * pag.PageSize
		return db.Offset(int(offset)).Limit(int(pag.PageSize))
	}
}
