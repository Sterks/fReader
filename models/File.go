package model

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//File - структура пользователя
type File struct {
	// gorm.Model
	TID                   int       `gorm:"column:f_id;primary_key"`
	TParent               int       `gorm:"column:f_parent"`
	TName                 string    `gorm:"column:f_name"`
	TArea                 string    `gorm:"column:f_area"`
	FileType              FileType  `gorm:"foregignkey:TType;association_foreignkey:ft_id"`
	TType                 int8      `gorm:"column:f_type"`
	THash                 string    `gorm:"column:f_hash"`
	TSize                 int64     `gorm:"column:f_size"`
	CreatedAt             time.Time `gorm:"column:f_date_create"`
	TDateCreateFromSource time.Time `gorm:"column:f_date_create_from_source"`
	TFullpath             string    `gorm:"column:f_fullpath"`
	TDateLastCheck        time.Time `gorm:"column:f_date_last_check"`
	TFileIsDir            string    `gorm:"column:f_file_is_dir"`
}

// TableName ...
func (File) TableName() string {
	return "Files"
}
