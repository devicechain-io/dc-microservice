/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package rdb

// Entity that is referenced by a unique token which may change over time.
type TokenReference struct {
	Token string `gorm:"unique;not null;size:128"`
}

// Entity that has a name and description.
type NamedEntity struct {
	Name        string `gorm:"size:128"`
	Description string `gorm:"size:1024"`
}

// Entity that has branding information.
type BrandedEntity struct {
	ImageUrl        string `gorm:"size:512"`
	Icon            string `gorm:"size:128"`
	BackgroundColor string `gorm:"size:32"`
	ForegroundColor string `gorm:"size:32"`
	BorderColor     string `gorm:"size:32"`
}
