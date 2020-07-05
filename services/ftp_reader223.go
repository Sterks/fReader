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
	"github.com/Sterks/fReader/amqp"
	"github.com/Sterks/fReader/common"
	"github.com/Sterks/fReader/config"
	"github.com/Sterks/fReader/logger"
	"github.com/secsy/goftp"
)

type FtpReader223 struct {
	config *config.Config
	ftp    *goftp.Client
	Db     *db.Database
	logger *logger.Logger
	amq    *amqp.ProducerMQ
}

func NewFtpReader223(conf *config.Config) *FtpReader223 {
	return &FtpReader223{
		config: conf,
		Db:     &db.Database{},
		ftp:    &goftp.Client{},
		logger: &logger.Logger{},
		amq:    &amqp.ProducerMQ{},
	}
}

//Connect44 конфигурирование ftpClienta
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
	f.logger.ConfigureLogger(config)
	f.Db.OpenDatabase()
	f.logger.InfoLog("Сервис запускается ...", "")
	return f
}

// TaskManager ...
func (f *FtpReader223) TaskManager(typeFile string, config *config.Config) {
	//str := "2020-04-18"
	//from, _ := time.Parse(time.RFC3339, str)
	//to := time.Now()
	now := time.Now()
	y, m, d := now.Date()
	from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	to := time.Now()

	f.config = config
	f.logger.InfoLog("Запуск загрузки ", typeFile)
	t1 := time.Now()

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
		for _, value := range massNotice {
			pathServer := fmt.Sprintf("%s/%s/%s", rootPath, region, value)
			f.GetFileInfo(pathServer, from, to, region, value, typeFile)
		}
	}
	t2 := time.Now()
	t3 := t2.Sub(t1)
	f.logger.InfoLog("Время работы загрузки ", t3.String())
	f.logger.InfoLog("Загрузка завершена \n", typeFile)
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

// GetFileInfo ...
func (f *FtpReader223) GetFileInfo(path string, from time.Time, to time.Time, region string, fileTT string, file string) {
	fmt.Println(path)
	client := f.ftp
	Walk(client, path, func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			// no permissions is okay, keep walking
			if err.(goftp.Error).Code() == 550 {
				return nil
			}
			return err
		}

		var hash string
		res, hash := f.Db.CheckerExistFileDBNotHash(info)
		if res == 0 {
			id := f.Db.LastID()
			var file []byte
			hash, file = f.CheckDownloder(id, client, fullPath)
			pattern := ".+Notice"
			reg := info.Name()
			matched, err := regexp.MatchString(pattern, reg)
			if err != nil {
				log.Printf("Не могу распознать слова - %v", err)
			}
			if matched {
				f.amq.PublishSend(f.config, info, "Notifications223", file, id, region, fullPath, fileTT)
			} else {
				f.amq.PublishSend(f.config, info, "Protocols223", file, id, region, fullPath, fileTT)
			}
		}
		f.Db.CreateInfoFile(info, region, hash, fullPath, fileTT, file)

		return nil
	}, from, to, region)
}

func (f *FtpReader223) CheckDownloder(id int, client *goftp.Client, fullPath string) (string, []byte) {
	// if id != 0 {
	pathLocal := common.CreateFolder(f.config, id)

	nameFile := common.GenerateID(id)
	buf := new(bytes.Buffer)
	file, _ := os.Create(f.config.Directory.MainFolder + "/" + pathLocal + nameFile)
	defer file.Close()
	infoBuf := io.TeeReader(buf, file)
	err := client.Retrieve(fullPath, buf)
	if err != nil {
		log.Println(err)
	}
	var hasher = sha256.New()
	_, err = io.Copy(hasher, infoBuf)

	if err != nil {
		log.Println(err)
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	fileRead, err2 := ioutil.ReadFile(f.config.Directory.MainFolder + "/" + pathLocal + nameFile)
	if err2 != nil {
		log.Printf("Не могу прочитать файл \n", err2)
	}
	return hash, fileRead
	// }
	// return "", nil
}
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
