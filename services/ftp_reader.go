package services

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/Sterks/FReader/common"
	"github.com/Sterks/FReader/config"
	"github.com/Sterks/FReader/db"
	"github.com/Sterks/FReader/logger"
	"github.com/Sterks/FReader/router"

	"github.com/secsy/goftp"
)

//FtpReader ...
type FtpReader struct {
	config *config.Config
	ftp    *goftp.Client
	Db     *db.Database
	router *router.WebServer
	logger *logger.Logger
}

// New инициализация сервера
func New(conf *config.Config) *FtpReader {
	return &FtpReader{
		config: &config.Config{},
		Db:     &db.Database{},
		ftp:    &goftp.Client{},
		router: &router.WebServer{},
		logger: &logger.Logger{},
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
	f.Db.OpenDatabase(config, f.logger)
	f.logger.InfoLog("Сервис запускается ...", "")
	return f
}

// GetFileInfo ...
func (f *FtpReader) GetFileInfo(path string, rev bool, down bool, addDb bool, from time.Time, to time.Time, region string, hashReader bool) {
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

		var Hash string
		if rev == true {
			if addDb == true {
				if down == true {
					if hashReader == true {
						id := f.Db.CreateInfoFile(info, region, Hash, fullPath)
						if id != 0 {
							pathLocal := common.CreateFolder(f.config, id)

							// stringID := strconv.Itoa(id)
							nameFile := common.GenerateID(id)
							buf := new(bytes.Buffer)
							file, _ := os.Create(f.config.Directory.MainFolder + "/" + pathLocal + nameFile)
							// file, _ := os.Create(f.config.FileDir + "/" + m + "/" + stringID)
							defer file.Close()
							infoBuf := io.TeeReader(buf, file)
							err = client.Retrieve(fullPath, buf)
							if err != nil {
								log.Println(err)
							}
							var hasher = sha256.New()
							_, err = io.Copy(hasher, infoBuf)
							if err != nil {
								log.Println(err)
							}
							Hash = hex.EncodeToString(hasher.Sum(nil))
						}
						// fmt.Println(fullPath, Hash)
					} else {
						buf := new(bytes.Buffer)
						err = client.Retrieve(fullPath, buf)
						if err != nil {
							log.Println(err)
						}
						var hasher = sha256.New()
						_, err = io.Copy(hasher, buf)
						if err != nil {
							log.Println(err)
						}
						Hash = hex.EncodeToString(hasher.Sum(nil))
						// fmt.Println(fullPath, Hash)
					}
				}
			} else {
				buf := new(bytes.Buffer)
				err = client.Retrieve(fullPath, buf)
				if err != nil {
					log.Println(err)
				}
				var hasher = sha256.New()
				_, err = io.Copy(hasher, buf)
				if err != nil {
					log.Println(err)
				}
				Hash = hex.EncodeToString(hasher.Sum(nil))
				fmt.Println(fullPath, Hash)
				return nil
			}
		} else {
			fmt.Println(fullPath)
		}

		return nil
	}, rev, down, from, to, region, hashReader)
}

// Walk Гуляем по диреториям
func Walk(client *goftp.Client, root string, walkFn filepath.WalkFunc, rev bool, down bool, from time.Time, to time.Time, region string, hashReader bool) (ret error) {
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
func (f *FtpReader) TaskManager(from time.Time, to time.Time, typeFile string, config *config.Config) {
	f.config = config
	f.logger.InfoLog("Запуск загрузки ", typeFile)
	t1 := time.Now()

	listRegions := f.GetListFolderDb()
	for _, region := range listRegions {
		// rootPath := "/fcs_regions"
		rootPath := f.config.Directory.RootPath
		pathServer := fmt.Sprintf("%s/%s/%s", rootPath, region, typeFile)
		f.GetFileInfo(pathServer, true, true, true, from, to, region, true)
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
