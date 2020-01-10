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
	"sync"
	"sync/atomic"
	"time"

	"github.com/Sterks/FReader/config"
	"github.com/Sterks/FReader/services/db"
	"github.com/patrickmn/go-cache"
	"github.com/secsy/goftp"
)

//FtpReader ...
type FtpReader struct {
	config *config.Config
	ftp    *goftp.Client
	Db     *db.Store
}

// New инициализация сервера
func New(conf *config.Config) *FtpReader {
	return &FtpReader{
		config: conf,
		Db:     &db.Store{},
		ftp:    &goftp.Client{},
	}
}

//Connect конфигурирование ftpClienta
func (f *FtpReader) Connect() (*goftp.Client, error) {
	ftpServ := goftp.Config{
		User:     "free",
		Password: "free",
		// Logger:   os.Stderr,
		Timeout: 100 * time.Second,
	}
	c, err := goftp.DialConfig(ftpServ, "ftp.zakupki.gov.ru:21")
	if err != nil {
		return nil, err
	}
	f.ftp = c
	return c, nil
}

// Start ...
func (f *FtpReader) Start() *FtpReader {
	ftp, err := f.Connect()
	if err != nil {
		log.Printf("Проблемы с соединением - %v", err)
	}
	f.ftp = ftp

	f.Db.OpenDb()

	log.Println("Сервис запускается ...")

	return f
}

// GetFileInfo ...
func (f *FtpReader) GetFileInfo(path string, rev bool, down bool, addDb bool, from time.Time, to time.Time, region string, hashReader bool) {
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
						buf := new(bytes.Buffer)
						file, _ := os.Create(info.Name())
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
				f.Db.CreateInfoFile(info, region, Hash, fullPath)
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

// GetListFolder ...
func (f *FtpReader) GetListFolder(c *cache.Cache) []string {
	rootPath := "/fcs_regions"
	var listFolder []os.FileInfo
	massFiles, found := c.Get("massFiles")
	if found {
		listFolder = massFiles.([]os.FileInfo)
	} else {
		listFolder, _ = f.ftp.ReadDir(rootPath)
		// if err2 != nil {
		// 	log.Printf("Соединение - %v", err2)
		// }
	}

	var listPath []string
	for _, value := range listFolder {
		if value.IsDir() == true {
			listPath = append(listPath, value.Name())
			c.Set("massFiles", value.Name(), cache.DefaultExpiration)
		}
	}

	log.Printf("Получен список папок в /fcs_regions")

	return listPath
}

// TaskManager ...
func (f *FtpReader) TaskManager(from time.Time, to time.Time, typeFile string, wg *sync.WaitGroup, c *cache.Cache) {
	defer wg.Done()
	listRegions := f.GetListFolder(c)
	for _, region := range listRegions {
		rootPath := "/fcs_regions"
		pathServer := fmt.Sprintf("%s/%s/%s", rootPath, region, typeFile)
		f.GetFileInfo(pathServer, true, true, true, from, to, region, false)
	}
}
