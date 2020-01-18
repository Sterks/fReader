package model

import "time"

type SourceRegions struct {
	RID         int       `gorm:"column:r_id;primary_key"`
	RName       string    `gorm:"column:r_name"`
	RDateCreate time.Time `gorm:"column:r_date_create"`
	RDateUpdate time.Time `gorm:"column:r_date_update"`
}
