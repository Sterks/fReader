package services

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	model "github.com/Sterks/Pp.Common.Db/models"
	"github.com/Sterks/fReader/pkg/models"

	"github.com/Sterks/Pp.Common.Db/db"

	"github.com/Sterks/fReader/amqp"
	"github.com/Sterks/fReader/common"
	"github.com/Sterks/fReader/config"
	"github.com/Sterks/fReader/logger"
	"github.com/secsy/goftp"
)

//FtpReader44 Структура для 44ФЗ
type FtpReader44 struct {
	config *config.Config
	ftp    *goftp.Client
	Db     *db.Database
	Logger *logger.Logger
	amq    *amqp.ProducerMQ
	Data   *DatePeriod
}

//DatePeriod Структура для даты
type DatePeriod struct {
	From time.Time
	To   time.Time
}

// AddTimeNow Установливаем дату на текущий момент
func (f *FtpReader44) AddTimeNow() {
	f.Data.To = time.Now()
	y, m, d := f.Data.To.Date()
	f.Data.From = time.Date(y, m, d, 0, 0, 0, 0, f.Data.To.Location())
}

// NewFtpReader44 Новая структура по 44ФЗ
func NewFtpReader44(conf *config.Config) *FtpReader44 {
	return &FtpReader44{
		config: conf,
		Db:     &db.Database{},
		ftp:    &goftp.Client{},
		Logger: logger.NewLogger(),
		amq:    &amqp.ProducerMQ{},
		Data:   &DatePeriod{},
	}
}

// Connect Подключаемся к FTP серверу
func (f *FtpReader44) Connect(config *config.Config) *goftp.Client {
	ftpServ := goftp.Config{
		User:     config.FTPServer44.Username,
		Password: config.FTPServer44.Password,
		Timeout:  3 * time.Minute,
	}
	c, err := goftp.DialConfig(ftpServ, config.FTPServer44.Url44)
	if err != nil {
		f.Logger.ErrorLog("Не могу подключится к FTP серверу", err)
	}
	f.ftp = c
	return c
}

// Start Запускаем сервис и указываем настройки
func (f *FtpReader44) Start(config *config.Config) { // *FtpReader44 {
	f.Logger.ConfigureLogger(config)
	f.amq.Logger = f.Logger
	f.Db.OpenDatabase(f.config.Postgres.Host, f.config.Postgres.Port, f.config.Postgres.User, f.config.Postgres.Password, f.config.Postgres.DBName)
	f.Logger.InfoLog("Сервис запускается ...")
}

//GetListFolder ...
func (f *FtpReader44) GetListFolder() {
	rootPath := f.config.FTPServer44.RootPath
	listFolder, err := f.ftp.ReadDir(rootPath)
	if err != nil {
		f.Logger.InfoLog("Не удается подключиться к FTP серверу", err)
	}
	//var regions []string
	for _, value := range listFolder {
		if value.IsDir() == true {
			if f.Db.CheckRegionsDb(value.Name()) == 0 {
				f.Db.AddRegionsDb(value.Name(), "44 ФЗ")
				f.Logger.InfoLog("Добавлен регион %v\n", value.Name())
			}
		}
	}
}

//GetAllFolderRegionsDb Все регионы из базы
func (f *FtpReader44) GetAllFolderRegionsDb() []model.SourceRegions {
	regions := f.Db.ReaderRegionsDb()
	return regions
}

// GetFileInfo получение информации о файлеss
func (f *FtpReader44) GetFileInfo(regions []model.SourceRegions, typeFile string) []models.InformationFile {
	var informations []models.InformationFile
	for _, region := range regions {
		var tt string
		if typeFile == "notifications44" {
			tt = "notifications"
		} else {
			tt = "protocols"
		}
		ff := fmt.Sprintf("/fcs_regions/%s/%s", region.RName, tt)
		if err2 := common.Walk(f.ftp, ff, func(fullPath string, info os.FileInfo, err error) error {
			if err != nil {
				// no permissions is okay, keep walking
				if err.(goftp.Error).Code() == 550 {
					return nil
				}
				return err
			}
			var informationFile models.InformationFile
			informationFile.Inform = info
			informationFile.Fullpath = fullPath
			informationFile.Region = region
			informationFile.TypeFile = typeFile
			informations = append(informations, informationFile)
			f.Logger.InfoLog("Сейчас обрабатываем файл -", fullPath)
			return nil
		}, f.Data.From, f.Data.To); err2 != nil {
			f.Logger.ErrorLog("Информация из Walk", err2)
		}
	}
	return informations
}

// CheckDownloder ...
func (f *FtpReader44) CheckDownloder(listFiles []models.InformationFile) []models.InformationFile {
	var hash string
	var fileRead []byte
	for _, fileList := range listFiles {
		id := f.Db.LastID()
		pathLocal := common.CreateFolder(f.config, id)
		nameFile := common.GenerateID(id)
		buf := new(bytes.Buffer)
		file, _ := os.Create(f.config.Directory.MainFolder + "/" + pathLocal + nameFile)
		defer file.Close()
		infoBuf := io.TeeReader(buf, file)
		err := f.ftp.Retrieve(fileList.Fullpath, buf)
		if err != nil {
			f.Logger.ErrorLog("Не могу скачать файл", err)
		}
		var hasher = sha256.New()
		_, err = io.Copy(hasher, infoBuf)
		if err != nil {
			f.Logger.ErrorLog("Не могу определить hash", err)
		}
		hash = hex.EncodeToString(hasher.Sum(nil))
		_, err = io.Copy(file, buf)
		if err != nil {
			f.Logger.ErrorLog("Не могу определить скопировать буфер для загрузки на диск", err)
		}
		fileRead, _ = ioutil.ReadFile(f.config.Directory.MainFolder + "/" + pathLocal + nameFile)
		fileList.Hash = hash
		fileList.Raw = fileRead
		rez := f.Db.CheckExistFileDb(fileList.Inform, fileList.Hash)
		if rez == 0 {
			f.Db.CreateInfoFile(fileList.Inform, fileList.Region.RName, fileList.Hash, fileList.Fullpath, fileList.TypeFile, fileList.TypeFile)
			f.amq.PublishSend(f.config, fileList.Inform, fileList.TypeFile, fileList.Raw, id, fileList.Region.RName, fileList.Fullpath, fileList.TypeFile)
		} else {
			f.Logger.InfoLog("Файл существует", fileList.Inform.Name())
		}
	}
	//return hash, fileRead
	return listFiles
}
