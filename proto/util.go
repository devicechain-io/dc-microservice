/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package proto

// Creates a uint64* from uint*.
func NullUint64Of(value *uint) *uint64 {
	if value != nil {
		conv := uint64(*value)
		return &conv
	}
	return nil
}

// Creates a uint* from uint64*.
func NullUintOf(value *uint64) *uint {
	if value != nil {
		conv := uint(*value)
		return &conv
	}
	return nil
}
