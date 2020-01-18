package model

type FileType struct {
	// gorm.Model
	FTID   int8   `gorm:"column:ft_id;primary_key"`
	FTName string `gorm:"column:ft_name"`
	FTExt  string `gorm:"column:ft_ext"`
}

func (FileType) TableName() string {
	return "FilesTypes"
}
