package models

type SourceResources struct {
	SRID int `gorm:"column:sr_id;primary_key"`
	SRNAME string `gorm:"column:sr_name"`
	SRFULLNAME string `gorm:"column:sr_fullname"`
}