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

type FtpReader223 struct {
	config *config.Config
	ftp    *goftp.Client
	Db     *db.Database
	Logger *logger.Logger
	amq    *amqp.ProducerMQ
	Data   *DatePeriod223
}

func (f *FtpReader223) AddRegions() {
	panic("implement me")
}

type DatePeriod223 struct {
	From time.Time
	To   time.Time
}

//AddTimeNow Установливаем дату на текущий момент
func (f *FtpReader223) AddTimeNow() {
	f.Data.To = time.Now()
	y, m, d := f.Data.To.Date()
	f.Data.From = time.Date(y, m, d, 0, 0, 0, 0, f.Data.To.Location())
}

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

func (f *FtpReader223) GetAllFolderRegionsDb() []model.SourceRegions {
	regions := f.Db.ReaderRegionsDb()
	return regions
}

// // Установливаем дату на текущий момент
// func (da *DatePeriod223) AddTimeNow223() {
// 	da.To = time.Now()
// 	y, m, d := da.To.Date()
// 	da.From = time.Date(y, m, d, 0, 0, 0, 0, da.To.Location())
// }

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
		f.Logger.ErrorLog("Не могу подключится к FTP серверу - %v", err)
	}
	f.ftp = c
	return c
}

// Start ...
func (f *FtpReader223) Start(config *config.Config) {
	f.Logger.ConfigureLogger(config)
	f.Db.OpenDatabase()
	f.Logger.InfoLog("Сервис запускается ...", "")
}

//Connect223 конфигурирование ftpClienta
func (f *FtpReader223) Connect223(user string, password string, hostname string) (*goftp.Client, error) {
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

func (f *FtpReader223) Start223(config *config.Config) *FtpReader223 {
	ftp, err := f.Connect223(
		config.FTPServer223.Username,
		config.FTPServer223.Password,
		config.FTPServer223.Url223)
	if err != nil {
		log.Printf("Проблемы с соединением - %v", err)
	}
	f.ftp = ftp
	f.Logger.ConfigureLogger(config)
	f.Db.OpenDatabase()
	f.Logger.InfoLog("Сервис запускается ...", "")
	return f
}

// TaskManager ...
func (f *FtpReader223) TaskManager(typeFile string, config *config.Config) {
	//str := "2020-04-18"
	//from, _ := time.Parse(time.RFC3339, str)
	//to := time.Now()
	//now := time.Now()
	//y, m, d := now.Date()
	//from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	//to := time.Now()

	t1 := time.Now()
	f.config = config
	f.Logger.InfoLog("Запуск загрузки ", typeFile)

	listRegions := f.GetListFolderDb()
	for _, region := range listRegions {
		rootPath := f.config.FTPServer223.RootPath
		pathRegions := fmt.Sprint(rootPath + "/" + region)
		listFolder, err := f.ftp.ReadDir(pathRegions)
		if err != nil {
			log.Printf("Не могу прочитать директорию - %v", err)
		}
		// rootPath := "/fcs_regions"
		var massNotice []string
		for _, value := range listFolder {
			if typeFile == "notifications223" {
				pattern := ".+Notice"
				reg := value.Name()
				matched, err := regexp.MatchString(pattern, reg)
				if err != nil {
					log.Printf("Не могу распознать слова - %v", err)
				}
				if matched {
					massNotice = append(massNotice, value.Name())
				}
			} else if typeFile == "protocols223" {
				pattern := ".+Protocol"
				reg := value.Name()
				matched, err := regexp.MatchString(pattern, reg)
				if err != nil {
					log.Printf("Не могу распознать слова - %v", err)
				}
				if matched {
					massNotice = append(massNotice, value.Name())
				}
			}
		}
	}
	t2 := time.Now()
	t3 := t2.Sub(t1)
	f.Logger.InfoLog("Время работы загрузки ", t3.String())
	f.Logger.InfoLog("Загрузка завершена \n", typeFile)
}

//GetListFolder ...
func (f *FtpReader223) GetListFolder() {
	rootPath := f.config.FTPServer223.RootPath
	listFolder, err := f.ftp.ReadDir(rootPath)
	if err != nil {
		log.Printf("Не удается подключиться к FTP серверу - ошибка %v", err)
	}
	for _, value := range listFolder {
		if value.IsDir() == true {
			if f.Db.CheckRegionsDb(value.Name()) == 0 {
				f.Db.AddRegionsDb(value.Name(), "223 ФЗ")
				fmt.Printf("Добавлен регион %v\n", value.Name())
			}
		}
	}
	fmt.Println("Наличие новых регионов проверено")
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
				matched, err := regexp.MatchString(pattern, value.Name())
				if err != nil {
					log.Fatalln(err)
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
						informations = append(informations, informationFile)
						f.Logger.InfoLog("Сейчас обрабатываем файл - ", fullPath)
						return nil
					}, f.Data.From, f.Data.To); err2 != nil {
						f.Logger.ErrorLog("Информация из Walk - %s", err2)
					}
				}
				matched2, err := regexp.MatchString(pattern2, value.Name())
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
						f.Logger.InfoLog("Сейчас обрабатываем файл - ", fullPath)
						return nil
					}, f.Data.From, f.Data.To); err2 != nil {
						f.Logger.ErrorLog("Информация из Walk - %s", err2)
					}
				}
			}
		}
	}
	return informations
}

// func (f *FtpReader223) GetFileInfo(path string, from time.Time, to time.Time, region string, fileTT string, file string) {
// 	fmt.Println(path)
// 	client := f.ftp
// 	common.Walk(client, path, func(fullPath string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			// no permissions is okay, keep walking
// 			if err.(goftp.Error).Code() == 550 {
// 				return nil
// 			}
// 			return err
// 		}

// 		var hash string
// 		res, hash := f.Db.CheckerExistFileDBNotHash(info)
// 		if res == 0 {
// 			id := f.Db.LastID()
// 			var file []byte
// 			hash, file = f.CheckDownloder(id, client, fullPath)
// 			pattern := ".+Notice"
// 			reg := info.Name()
// 			matched, err := regexp.MatchString(pattern, reg)
// 			if err != nil {
// 				log.Printf("Не могу распознать слова - %v", err)
// 			}
// 			if matched {
// 				f.amq.PublishSend(f.config, info, "Notifications223", file, id, region, fullPath, fileTT)
// 			} else {
// 				f.amq.PublishSend(f.config, info, "Protocols223", file, id, region, fullPath, fileTT)
// 			}
// 		}
// 		f.Db.CreateInfoFile(info, region, hash, fullPath, fileTT, file)

// 		return nil
// 	}, from, to, region)
// }

// func (f *FtpReader223) CheckDownloder(id int, client *goftp.Client, fullPath string) (string, []byte) {
// 	// if id != 0 {
// 	pathLocal := common.CreateFolder(f.config, id)

// 	nameFile := common.GenerateID(id)
// 	buf := new(bytes.Buffer)
// 	file, _ := os.Create(f.config.Directory.MainFolder + "/" + pathLocal + nameFile)
// 	defer file.Close()
// 	infoBuf := io.TeeReader(buf, file)
// 	err := client.Retrieve(fullPath, buf)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	var hasher = sha256.New()
// 	_, err = io.Copy(hasher, infoBuf)

// 	if err != nil {
// 		log.Println(err)
// 	}
// 	hash := hex.EncodeToString(hasher.Sum(nil))

// 	fileRead, err2 := ioutil.ReadFile(f.config.Directory.MainFolder + "/" + pathLocal + nameFile)
// 	if err2 != nil {
// 		log.Printf("Не могу прочитать файл \n", err2)
// 	}
// 	return hash, fileRead
// 	// }
// return "", nil
// }

func (f *FtpReader223) FirstChecherRegions() {
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
