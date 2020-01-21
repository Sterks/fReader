package db

import (
	"fmt"
	"os"
	"time"

	"github.com/Sterks/FReader/logger"
	model "github.com/Sterks/FReader/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //....
)

//Database ...
type Database struct {
	database *gorm.DB
	logger   *logger.Logger
}

const (
	host     = "localhost"
	port     = 5432
	user     = "user_ro"
	password = "4r2w3e1q"
	dbname   = "freader"
)

// OpenDatabase ...
func (d *Database) OpenDatabase() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		d.logger.ErrorLog("Соединиться не удалось - %s", err)
	}
	if err2 := db.DB().Ping(); err2 != nil {
		d.logger.ErrorLog("База не отвечает", err2)
	}
	d.database = db
}

// CreateInfoFile ...
func (d *Database) CreateInfoFile(info os.FileInfo, region string, hash string, fullpath string) {
	// d.database.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&files)
	// filesTypes := d.database.Table("FileType")
	d.database.LogMode(true)

	checker := d.CheckExistFileDb(info, hash)
	if checker != 0 {
		var lf model.File
		d.database.Table("Files").Where("f_id = ?", checker).Find(&lf)
		lf.TDateLastCheck = time.Now()
		d.database.Save(&lf)
	}
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

// CheckExistFileDb ...
func (d *Database) CheckExistFileDb(file os.FileInfo, hash string) int {
	var ff model.File
	d.database.Table("Files").Where("f_hash = ? and f_size = ? and f_name = ?", hash, file.Size(), file.Name()).Find(&ff)
	return ff.TID
}

//CheckRegionsDb Проверка существует ли регион в базе данных
func (d *Database) CheckRegionsDb(region string) int {
	var reg model.SourceRegions
	d.database.Table("SourceRegions").Where("r_name = ?", region).First(&reg)
	return reg.RID
}

//ReaderRegionsDb Все регионы из базы
func (d *Database) ReaderRegionsDb() []model.SourceRegions {
	var regions []model.SourceRegions
	d.database.Table("SourceRegions").Find(&regions)
	return regions
}

//AddRegionsDb ...
func (d *Database) AddRegionsDb(region string) {
	var reg model.SourceRegions
	reg.RName = region
	reg.RDateCreate = time.Now()
	reg.RDateUpdate = time.Now()
	d.database.Table("SourceRegions").Create(&reg)
}
