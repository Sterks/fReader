package services

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	router2 "github.com/Sterks/fReader/router"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/Sterks/Pp.Common.Db/db"

	"github.com/Sterks/fReader/amqp"
	"github.com/Sterks/fReader/common"
	"github.com/Sterks/fReader/config"
	"github.com/Sterks/fReader/logger"
	"github.com/secsy/goftp"
)

//FtpReader ...
type FtpReader struct {
	config *config.Config
	ftp    *goftp.Client
	Db     *db.Database
	router *router2.WebServer
	logger *logger.Logger
	amq    *amqp.ProducerMQ
}

// New инициализация сервера
func New(conf *config.Config) *FtpReader {
	return &FtpReader{
		config: &config.Config{},
		Db:     &db.Database{},
		ftp:    &goftp.Client{},
		router: &router2.WebServer{},
		logger: &logger.Logger{},
		amq:    &amqp.ProducerMQ{},
	}
}

//Connect конфигурирование ftpClienta
func (f *FtpReader) Connect() (*goftp.Client, error) {
	ftpServ := goftp.Config{
		User:     "free",
		Password: "free",
		// Logger:   os.Stderr,
		Timeout: 3 * time.Minute,
	}
	c, err := goftp.DialConfig(ftpServ, "ftp.zakupki.gov.ru:21")
	if err != nil {
		return nil, err
	}
	f.ftp = c
	return c, nil
}

// Start ...
func (f *FtpReader) Start(config *config.Config) *FtpReader {
	ftp, err := f.Connect()
	if err != nil {
		log.Printf("Проблемы с соединением - %v", err)
	}
	f.ftp = ftp
	f.logger.ConfigureLogger(config)
	f.Db.OpenDatabase()
	f.logger.InfoLog("Сервис запускается ...", "")
	return f
}

// GetFileInfo ...
func (f *FtpReader) GetFileInfo(path string, from time.Time, to time.Time, region string, fileTT string) {
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
			f.amq.PublishSend(f.config, info, "Files", file, id, region, fullPath, fileTT)
		}
		f.Db.CreateInfoFile(info, region, hash, fullPath, fileTT)

		return nil
	}, from, to, region)
}

// Walk Гуляем по диреториям
func Walk(client *goftp.Client, root string, walkFn filepath.WalkFunc, from time.Time, to time.Time, region string) (ret error) {
	dirsToCheck := make(chan string, 100)

	var workCount int32 = 1
	dirsToCheck <- root

	for dir := range dirsToCheck {
		go func(dir string) {
			files, err := client.ReadDir(dir)

			if err != nil {
				if err = walkFn(dir, nil, err); err != nil && err != filepath.SkipDir {
					ret = err
					close(dirsToCheck)
					return
				}
			}

			for _, file := range files {
				if file.ModTime().After(from) && file.ModTime().Before(to) && file.IsDir() == false {
					if err = walkFn(path.Join(dir, file.Name()), file, nil); err != nil {
						if file.IsDir() && err == filepath.SkipDir {
							continue
						}
						ret = err
						close(dirsToCheck)
						return
					}
				}

				if file.IsDir() {
					atomic.AddInt32(&workCount, 1)
					dirsToCheck <- path.Join(dir, file.Name())
				}
			}

			atomic.AddInt32(&workCount, -1)
			if workCount == 0 {
				close(dirsToCheck)
			}
		}(dir)
	}

	return ret
}

// GetListFolderFtp ...
func (f *FtpReader) GetListFolderFtp() []string {
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
func (f *FtpReader) GetListFolderDb() []string {
	var listFolder []string
	listRegDb := f.Db.ReaderRegionsDb()
	for _, value := range listRegDb {
		if value.RID != 0 {
			listFolder = append(listFolder, value.RName)
		}
	}
	return listFolder
}

//GetListFolder ...
func (f *FtpReader) GetListFolder() {
	rootPath := "/fcs_regions"
	listFolder, err := f.ftp.ReadDir(rootPath)
	if err != nil {
		log.Printf("Не удается подключиться к FTP серверу - ошибка %v", err)
	}
	for _, value := range listFolder {
		if value.IsDir() == true {
			if f.Db.CheckRegionsDb(value.Name()) == 0 {
				f.Db.AddRegionsDb(value.Name())
				fmt.Printf("Добавлен регион %v\n", value.Name())
			}
		}
	}
}

// TaskManager ...
func (f *FtpReader) TaskManager(typeFile string, config *config.Config) {
	//str := "2020-03-06"
	//from, _ := time.Parse(time.RFC3339, str)
	//to := time.Now()
	now := time.Now()
	y, m, d := now.Date()
	from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	to := time.Now()
	f.FirstChecherRegions()

	f.config = config
	f.logger.InfoLog("Запуск загрузки ", typeFile)
	t1 := time.Now()

	listRegions := f.GetListFolderDb()
	for _, region := range listRegions {
		// rootPath := "/fcs_regions"
		rootPath := f.config.Directory.RootPath
		pathServer := fmt.Sprintf("%s/%s/%s", rootPath, region, typeFile)
		f.GetFileInfo(pathServer, from, to, region, typeFile)
	}
	t2 := time.Now()
	t3 := t2.Sub(t1)
	f.logger.InfoLog("Время работы загрузки ", t3.String())
	f.logger.InfoLog("Загрузка завершена \n", typeFile)
}

//FirstChecherRegions ...
func (f *FtpReader) FirstChecherRegions() {
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
func (f *FtpReader) CheckDownloder(id int, client *goftp.Client, fullPath string) (string, []byte) {
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
