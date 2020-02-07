package db

import (
	"fmt"
	"os"
	"time"

	"github.com/Sterks/FReader/config"
	"github.com/Sterks/FReader/logger"
	model "github.com/Sterks/FReader/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //....
)

//Database ...
type Database struct {
	config   *config.Config
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
func (d *Database) OpenDatabase(config *config.Config, logger *logger.Logger) {
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
	d.config = config
	d.logger = logger
}

// CreateInfoFile ...
func (d *Database) CreateInfoFile(info os.FileInfo, region string, hash string, fullpath string) int {
	// d.database.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&files)
	// filesTypes := d.database.Table("FileType")
	d.database.LogMode(true)

	var gf model.SourceRegions
	d.database.Table("SourceRegions").Where("r_name = ?", region).Find(&gf)

	checker := d.CheckExistFileDb(info, hash)
	if checker != 0 {
		var lf model.File
		d.database.Table("Files").Where("f_id = ?", checker).Find(&lf)
		lf.TDateLastCheck = time.Now()
		d.database.Save(&lf)
		d.logger.InfoLog("Дата успешно обновлена", lf.TDateLastCheck.String())
	}
	if checker == 0 {

		var fileType model.FileType
		var lastID model.File
		d.database.Table("FilesTypes").Where("ft_name = ?", "ZIP архив").Find(&fileType)

		d.database.Table("Files")
		d.database.Create(&model.File{
			TName:                 info.Name(),
			TArea:                 gf.RID,
			FileType:              fileType,
			TType:                 fileType.FTID,
			THash:                 hash,
			TSize:                 info.Size(),
			CreatedAt:             time.Now(),
			TDateCreateFromSource: info.ModTime(),
			TDateLastCheck:        time.Now(),
			TFullpath:             fullpath,
		}).Scan(&lastID)
		d.logger.InfoLog("Файл успешно добавлен - ", lastID.TName)
		return lastID.TID
	} else {
		fmt.Printf("Файл существует - %v\n", info.Name())
		return 0
	}
}

//LastID ...
func (d *Database) LastID() int {
	var ff model.File
	d.database.Table("Files").Last(&ff)
	return ff.TID
}

// CheckerExistFileDBNotHash ...
func (d *Database) CheckerExistFileDBNotHash(file os.FileInfo) (int, string) {
	var ff model.File
	fmt.Printf("%v - %v - %v", file.Size(), file.Name(), file.ModTime())
	d.database.Table("Files").Where("f_size = ? and f_name = ? and f_date_create_from_source = ?", file.Size(), file.Name(), file.ModTime()).Find(&ff)
	return ff.TID, ff.THash
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

// FirstOrCreate Создать или получить
func (d *Database) FirstOrCreate(region string) model.SourceRegions {
	var reg model.SourceRegions
	reg.RName = region
	reg.RDateCreate = time.Now()
	reg.RDateUpdate = time.Now()
	d.database.Table("SourceRegions").Where("r_name = ?", region).FirstOrCreate(&reg)
	return reg
}
