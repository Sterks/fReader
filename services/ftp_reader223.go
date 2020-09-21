package services

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/Sterks/Pp.Common.Db/db"
	model "github.com/Sterks/Pp.Common.Db/models"
	"github.com/Sterks/fReader/amqp"
	"github.com/Sterks/fReader/common"
	"github.com/Sterks/fReader/config"
	"github.com/Sterks/fReader/logger"
	"github.com/Sterks/fReader/pkg/models"
	"github.com/secsy/goftp"
)

// FtpReader223 ...
type FtpReader223 struct {
	config *config.Config
	ftp    *goftp.Client
	Db     *db.Database
	Logger *logger.Logger
	amq    *amqp.ProducerMQ
	Data   *DatePeriod223
}

//DatePeriod223 ...
type DatePeriod223 struct {
	From time.Time
	To   time.Time
}

//AddTimeNow Установливаем дату на текущий момент
func (f *FtpReader223) AddTimeNow() {
	f.Data.To = time.Now()
	y, m, d := f.Data.To.Date()
	f.Data.From = time.Date(y, m, d, 0, 0, 0, 0, f.Data.To.Location())
	// str := "2020-09-14"
	// from, _ := time.Parse("2006-01-02", str)
	// f.Data.From = from
}

// NewFtpReader223 Новая структура по 223
func NewFtpReader223(conf *config.Config) *FtpReader223 {
	return &FtpReader223{
		config: conf,
		Db:     &db.Database{},
		ftp:    &goftp.Client{},
		Logger: logger.NewLogger(),
		amq:    &amqp.ProducerMQ{},
		Data:   &DatePeriod223{},
	}
}

//Connect конфигурирование ftpClienta
func (f *FtpReader223) Connect(config *config.Config) *goftp.Client {
	ftpServ := goftp.Config{
		User:     config.FTPServer223.Username,
		Password: config.FTPServer223.Password,
		// Logger:   os.Stderr,
		Timeout: 3 * time.Minute,
	}
	c, err := goftp.DialConfig(ftpServ, config.FTPServer223.Url223)
	if err != nil {
		f.Logger.ErrorLog("Не могу подключится к FTP серверу ", err)
	}
	f.ftp = c
	return c
}

// Start ...
func (f *FtpReader223) Start(config *config.Config) {
	f.Logger.ConfigureLogger(config)
	f.Db.OpenDatabase(f.config.Postgres.Host, f.config.Postgres.Port, f.config.Postgres.User, f.config.Postgres.Password, f.config.Postgres.DBName)
	f.amq.Logger = f.Logger
	f.Logger.InfoLog("Сервис запускается ...")
}

//GetListFolder ...
func (f *FtpReader223) GetListFolder() {
	rootPath := f.config.FTPServer223.RootPath
	listFolder, err := f.ftp.ReadDir(rootPath)
	if err != nil {
		f.Logger.ErrorLog("Не удается подключиться к FTP серверу - ошибка", err)
	}
	for _, value := range listFolder {
		if value.IsDir() == true {
			if f.Db.CheckRegionsDb(value.Name()) == 0 {
				f.Db.AddRegionsDb(value.Name(), "223 ФЗ")
				f.Logger.InfoLog("Добавлен регион ", value.Name())
			}
		}
	}
	fmt.Println("Наличие новых регионов проверено")
}

// GetAllFolderRegionsDb ...
func (f *FtpReader223) GetAllFolderRegionsDb() []model.SourceRegions {
	regions := f.Db.ReaderRegionsDb()
	return regions
}

// GetFileInfo получение информации о файлеss
func (f *FtpReader223) GetFileInfo(regions []model.SourceRegions, typeFile string) []models.InformationFile {
	var informations []models.InformationFile
	for _, region := range regions {
		if region.RFZLaw == 2 {
			rootFolderInRegion := fmt.Sprintf("%s/%s", f.config.FTPServer223.RootPath, region.RName)
			listFolder, _ := f.ftp.ReadDir(rootFolderInRegion)
			pattern := ".Notice"
			pattern2 := ".+Protocol"
			for _, value := range listFolder {
				if typeFile == "notifications223" {
					matched, err := regexp.MatchString(pattern, value.Name())
					if err != nil {
						f.Logger.ErrorLog("Не выполнен regexp", err)
					}
					if matched {
						ff := fmt.Sprintf("/out/published/%s/%s", region.RName, value.Name())
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

							//Можно ускорить если передавать в канал
							informations = append(informations, informationFile)
							// f.Logger.InfoLog("Сейчас обрабатываем файл - ", fullPath)
							return nil
						}, f.Data.From, f.Data.To); err2 != nil {
							f.Logger.ErrorLog("Информация из Walk", err2)
						}
					}
				} else {
					matched2, err := regexp.MatchString(pattern2, value.Name())
					if err != nil {
						f.Logger.ErrorLog("Не работает regexp 2", err)
					}
					if matched2 {
						ff := fmt.Sprintf("/out/published/%s/%s", region.RName, value.Name())
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
							// f.Logger.InfoLog("Сейчас обрабатываем файл - ", fullPath)
							return nil
						}, f.Data.From, f.Data.To); err2 != nil {
							f.Logger.ErrorLog("Информация из Walk", err2)
						}
					}
				}
			}
		}
	}
	return informations
}

//GetListFolderDb ....
func (f *FtpReader223) GetListFolderDb() []string {
	var listFolder []string
	listRegDb := f.Db.ReaderRegionsDb()
	for _, value := range listRegDb {
		if value.RID != 0 && value.RFZLaw == 2 {
			listFolder = append(listFolder, value.RName)
		}
	}
	return listFolder
}

//CheckDownloder Загрузка
func (f *FtpReader223) CheckDownloder(listFiles []models.InformationFile) []models.InformationFile {
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
			f.Logger.ErrorLog("Не могу скачать файл с FTP", err)
		}
		var hasher = sha256.New()
		_, err = io.Copy(hasher, infoBuf)
		if err != nil {
			f.Logger.ErrorLog("Не могу посчитать хеш -", err)
		}
		hash = hex.EncodeToString(hasher.Sum(nil))
		_, err = io.Copy(file, buf)
		if err != nil {
			f.Logger.ErrorLog("Не могу посчитать хеш -", err)
		}
		fileRead, _ = ioutil.ReadFile(f.config.Directory.MainFolder + "/" + pathLocal + nameFile)
		fileList.Hash = hash
		fileList.Raw = fileRead
		rez := f.Db.CheckExistFileDb(fileList.Inform, fileList.Hash)
		if rez == 0 {
			f.Db.CreateInfoFile(fileList.Inform, fileList.Region.RName, fileList.Hash, fileList.Fullpath, fileList.TypeFile, fileList.TypeFile)
			f.amq.PublishSend(f.config, fileList.Inform, fileList.TypeFile, fileList.Raw, id, fileList.Region.RName, fileList.Fullpath, fileList.TypeFile)
		} else {
			f.Logger.InfoLog("Файл существует -", fileList.Inform.Name())
		}
	}
	//return hash, fileRead
	return listFiles
}
