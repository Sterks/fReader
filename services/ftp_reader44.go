package services

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

//FtpReader ...
type FtpReader44 struct {
	config *config.Config
	ftp    *goftp.Client
	Db     *db.Database
	Logger *logger.Logger
	amq    *amqp.ProducerMQ
	Data   *DatePeriod
}

func (f *FtpReader44) AddRegions() {
	panic("implement me")
}

type DatePeriod struct {
	From time.Time
	To   time.Time
}

// Установливаем дату на текущий момент
func (f *FtpReader44) AddTimeNow() {
	f.Data.To = time.Now()
	y, m, d := f.Data.To.Date()
	f.Data.From = time.Date(y, m, d, 0, 0, 0, 0, f.Data.To.Location())
}

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

func (f *FtpReader44) Connect(config *config.Config) *goftp.Client {
	ftpServ := goftp.Config{
		User:     config.FTPServer44.Username,
		Password: config.FTPServer44.Password,
		// Logger:   os.Stderr,
		Timeout: 3 * time.Minute,
	}
	c, err := goftp.DialConfig(ftpServ, config.FTPServer44.Url44)
	if err != nil {
		f.Logger.ErrorLog("Не могу подключится к FTP серверу - %v", err)
	}
	f.ftp = c
	return c
}

func (f *FtpReader44) Start(config *config.Config) { // *FtpReader44 {
	f.Logger.ConfigureLogger(config)
	f.Db.OpenDatabase()
	f.Logger.InfoLog("Сервис запускается ...", "")
	//return f
}

//func (f *FtpReader44) DownloaderFiles(listFiles []models.InformationFile){
//	for _, file := range listFiles {
//		ident := f.Db.LastID()
//		ident += ident
//		common.CreateFolder(f.config, ident)
//		pathLocal := common.PathLocalFile(f.config, ident)
//		createFile, err := os.Create(pathLocal)
//		if err != nil {
//			f.Logger.ErrorLog("Не могу создать файл - %v", err )
//		}
//		buf := new(bytes.Buffer)
//		if err := f.ftp.Retrieve(file.Fullpath, buf); err != nil {
//			f.Logger.ErrorLog("Не могу записать в буфер - %v", err )
//		}
//		var hasher = sha256.New()
//		_, err = io.Copy(hasher, buf)
//
//	}
//}

//Connect44 конфигурирование ftpClienta
func (f *FtpReader44) Connect44(user string, password string, hostname string) (*goftp.Client, error) {
	ftpServ := goftp.Config{
		User:     user,
		Password: password,
		// Logger:   os.Stderr,
		Timeout: 3 * time.Minute,
	}
	c, err := goftp.DialConfig(ftpServ, hostname)
	if err != nil {
		return nil, err
	}
	f.ftp = c
	return c, nil
}

// Start44 ...
func (f *FtpReader44) Start44(config *config.Config) *FtpReader44 {
	ftp, err := f.Connect44(
		config.FTPServer44.Username,
		config.FTPServer44.Password,
		config.FTPServer44.Url44)
	if err != nil {
		log.Printf("Проблемы с соединением - %v", err)
	}
	f.ftp = ftp
	f.Logger.ConfigureLogger(config)
	f.Db.OpenDatabase()
	f.Logger.InfoLog("Сервис запускается ...", "")
	return f
}

// GetFileInfo получение информации о файлеss
func (f *FtpReader44) GetFileInfo(regions []model.SourceRegions, typeFile string) []models.InformationFile {
	var informations []models.InformationFile
	for _, region := range regions {
		var tt string
		if typeFile == "notifications44" {
			tt = "notifications"
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
			f.Logger.InfoLog("Сейчас обрабатываем файл - ", fullPath)
			return nil
		}, f.Data.From, f.Data.To); err2 != nil {
			f.Logger.ErrorLog("Информация из Walk - %s", err2)
		}
	}
	return informations
}

func (f *FtpReader44) GetAllFolderRegionsDb() []model.SourceRegions {
	regions := f.Db.ReaderRegionsDb()
	return regions
}

// GetFileInfo ...
//func (f *FtpReader44) GetFileInfo(path string, from time.Time, to time.Time, region string, fileTT string) {
//	fmt.Println(path)
//	client := f.ftp
//	_ = common.Walk(client, path, func(fullPath string, info os.FileInfo, err error) error {
//		if err != nil {
//			// no permissions is okay, keep walking
//			if err.(goftp.Error).Code() == 550 {
//				return nil
//			}
//			return err
//		}
//
//		var hash string
//		res, hash := f.Db.CheckerExistFileDBNotHash(info)
//		if res == 0 {
//			id := f.Db.LastID()
//			var file []byte
//			hash, file = f.CheckDownloder(id, client, fullPath)
//			if fileTT == "notifications44" {
//				f.amq.PublishSend(f.config, info, "Notifications44", file, id, region, fullPath, fileTT)
//			} else {
//				f.amq.PublishSend(f.config, info, "Protocols44", file, id, region, fullPath, fileTT)
//			}
//		}
//		f.Db.CreateInfoFile(info, region, hash, fullPath, fileTT, fileTT)
//		return nil
//	}, from, to, region)
//}

// GetListFolderFtp ...
func (f *FtpReader44) GetListFolderFtp() []string {
	// rootPath := "/fcs_regions"
	rootPath := f.config.Directory.RootPath
	var listFolder []os.FileInfo
	listFolder, erro := f.ftp.ReadDir(rootPath)
	if erro != nil {
		log.Printf("Соединение - %v", erro)
	}
	var listPath []string
	for _, value := range listFolder {
		if value.IsDir() == true {
			listPath = append(listPath, value.Name())
		}
	}
	log.Printf("Получен список папок в /fcs_regions")

	return listPath
}

//GetListFolderDb ....
func (f *FtpReader44) GetListFolderDb() []string {
	var listFolder []string
	listRegDb := f.Db.ReaderRegionsDb()
	for _, value := range listRegDb {
		if value.RID != 0 && value.RFZLaw == 1 {
			listFolder = append(listFolder, value.RName)
		}
	}
	return listFolder
}

//GetListFolder ...
func (f *FtpReader44) GetListFolder() {
	rootPath := f.config.FTPServer44.RootPath
	listFolder, err := f.ftp.ReadDir(rootPath)
	if err != nil {
		log.Printf("Не удается подключиться к FTP серверу - ошибка %v", err)
	}
	//var regions []string
	for _, value := range listFolder {
		if value.IsDir() == true {
			if f.Db.CheckRegionsDb(value.Name()) == 0 {
				f.Db.AddRegionsDb(value.Name(), "44 ФЗ")
				fmt.Printf("Добавлен регион %v\n", value.Name())
			}
		}
	}
}

// TaskManager ...
//func (f *FtpReader44) TaskManager(typeFile string, config *config.Config) {
//	// str := "2020-07-04"
//	// from, _ := time.Parse(time.RFC3339, str)
//	// to := time.Now()
//	//now := time.Now()
//	//y, m, d := now.Date()
//	//from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
//	//to := time.Now()
//	f.FirstChecherRegions()
//
//	f.config = config
//	f.Logger.InfoLog("Запуск загрузки ", typeFile)
//	t1 := time.Now()
//
//	//Проверяем есть ли в базе записи о директориях
//	listRegions := f.GetListFolderDb()
//	//Получаем список регионов
//	for _, region := range listRegions {
//		rootPath := f.config.FTPServer44.RootPath
//		if typeFile == "notifications44" {
//			gg := "notifications"
//			pathServer := fmt.Sprintf("%s/%s/%s", rootPath, region, gg)
//			f.GetFileInfo(pathServer, region, typeFile)
//		} else if typeFile == "protocols44" {
//			gg := "protocols"
//			pathServer := fmt.Sprintf("%s/%s/%s", rootPath, region, gg)
//			f.GetFileInfo(pathServer, region, typeFile)
//		}
//	}
//	t2 := time.Now()
//	t3 := t2.Sub(t1)
//	f.Logger.InfoLog("Время работы загрузки ", t3.String())
//	f.Logger.InfoLog("Загрузка завершена \n", typeFile)
//}

//FirstChecherRegions ...
func (f *FtpReader44) FirstChecherRegions() {
	var checkVal string
	listFolder := f.GetListFolderDb()
	for _, value := range listFolder {
		if value != "" {
			continue
		} else {
			checkVal = "Есть"
		}
	}
	if checkVal == "" {
		f.GetListFolder()
	}
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
			log.Println(err)
		}
		var hasher = sha256.New()
		_, err = io.Copy(hasher, infoBuf)
		if err != nil {
			log.Println(err)
		}
		hash = hex.EncodeToString(hasher.Sum(nil))
		_, err = io.Copy(file, buf)
		if err != nil {
			log.Println(err)
		}
		fileRead, _ = ioutil.ReadFile(f.config.Directory.MainFolder + "/" + pathLocal + nameFile)
		fileList.Hash = hash
		fileList.Raw = fileRead
		rez := f.Db.CheckExistFileDb(fileList.Inform, fileList.Hash)
		if rez == 0 {
			f.Db.CreateInfoFile(fileList.Inform, fileList.Region.RName, fileList.Hash, fileList.Fullpath, fileList.TypeFile, fileList.TypeFile)
			f.amq.PublishSend(f.config, fileList.Inform, fileList.TypeFile, fileList.Raw, id, fileList.Region.RName, fileList.Fullpath, fileList.TypeFile)
		} else {
			fmt.Printf("Файл существует - %s\n", fileList.Inform.Name())
		}
	}
	//return hash, fileRead
	return listFiles
}
