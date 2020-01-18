package db

import (
	"fmt"
	"log"
	"os"
	"time"

	model "github.com/Sterks/FReader/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Database struct {
	database *gorm.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "user_ro"
	password = "4r2w3e1q"
	dbname   = "freader"
)

func (d *Database) OpenDatabase() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("Соединиться не удалось - %s", err)
	}
	if err2 := db.DB().Ping(); err2 != nil {
		log.Println(err2)
	}
	d.database = db

}

// CreateInfoFile ...
func (d *Database) CreateInfoFile(info os.FileInfo, region string, hash string, fullpath string) {
	// d.database.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&files)
	// filesTypes := d.database.Table("FileType")
	d.database.LogMode(true)

	checker := d.CheckExistFileDb(info, hash)
	if checker == 0 {

		var fileType model.FileType
		d.database.Table("FilesTypes").Where("ft_name = ?", "ZIP архив").Find(&fileType)

		d.database.Table("Files")
		d.database.Create(&model.File{
			TName:                 info.Name(),
			TArea:                 region,
			FileType:              fileType,
			TType:                 fileType.FTID,
			THash:                 hash,
			TSize:                 info.Size(),
			CreatedAt:             time.Now(),
			TDateCreateFromSource: info.ModTime(),
			TDateLastCheck:        time.Now(),
			TFullpath:             fullpath,
		})
	} else {
		fmt.Printf("Файл существует - %v\n", info.Name())
	}
}

func (d *Database) CheckExistFileDb(file os.FileInfo, hash string) int {
	var ff model.File
	d.database.Table("Files").Where("f_hash = ? and f_size = ? and f_name = ?", hash, file.Size(), file.Name()).Find(&ff)
	ff.TDateLastCheck = time.Now()
	d.database.Save(&ff)
	return ff.TID
}
