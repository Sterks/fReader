package common

import (
	"fmt"
	"github.com/secsy/goftp"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/Sterks/fReader/config"
	"github.com/sirupsen/logrus"
)

// GenerateID - Герерация строки длинной 12 символов
func GenerateID(ident int) string {
	word := strconv.Itoa(ident)
	ch := len(word)
	nool := 12 - ch
	var ap string
	ap = word
	for i := 0; i < nool; i++ {
		ap = fmt.Sprintf("0%s", ap)
	}
	return ap
}

//PathLocalFile Путь до файла на файловой системе
func PathLocalFile(config *config.Config, id int) string {
	stringID := GenerateID(id)
	lv1 := fmt.Sprint(stringID[0:3])
	lv2 := fmt.Sprint(stringID[3:6])
	lv3 := fmt.Sprint(stringID[6:9])
	strFormat := fmt.Sprint(config.Directory.MainFolder + "/" + lv1 + "/" + lv2 + "/" + lv3 + "/" + stringID)
	return strFormat
}

// CreateFolder ...
func CreateFolder(config *config.Config, ident int) string {
	saveDir := config.Directory.MainFolder
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		logrus.Errorf("Не могу создать директорию - %v\n", err)
	}
	stringID := GenerateID(ident)
	lv1 := fmt.Sprint(stringID[0:3])
	lv2 := fmt.Sprint(stringID[3:6])
	lv3 := fmt.Sprint(stringID[6:9])
	// lv4 := fmt.Sprintln(stringID[9:12])
	if err := os.MkdirAll(saveDir+"/"+lv1, 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(saveDir+"/"+lv1+"/"+lv2, 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(saveDir+"/"+lv1+"/"+lv2+"/"+lv3, 0755); err != nil {
		log.Fatal(err)
	}
	// if err := os.MkdirAll(lv4, 0755); err != nil {
	// 	log.Fatal(err)
	// }
	path := fmt.Sprintf("%s/%s/%s/", lv1, lv2, lv3)
	return path
}

// GetLocalPath ...
func GetLocalPath(config *config.Config, ident int) string {
	rootPath := config.Directory.MainFolder
	word := strconv.Itoa(ident)
	stringID := GenerateID(ident)
	lv1 := fmt.Sprint(stringID[0:3])
	lv2 := fmt.Sprint(stringID[3:6])
	lv3 := fmt.Sprint(stringID[6:9])
	s := fmt.Sprint(rootPath + "/" + lv1 + "/" + lv2 + "/" + lv3 + "/" + word)
	return s
}

// Walk Гуляем по диреториям
func Walk(client *goftp.Client, root string, walkFn filepath.WalkFunc, from time.Time, to time.Time) (ret error) {
	dirsToCheck := make(chan string, 100)

	var workCount int32 = 1
	dirsToCheck <- root

	for dir := range dirsToCheck {
		go func(dir string) {
			files, err := client.ReadDir(dir)
			if err != nil {
				if err = walkFn(dir, nil, err); err != nil && err != filepath.SkipDir {
					Logging("Log.txt", root, dir, filepath.Dir(dir) )
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
						Logging("Log.txt", root, dir, file.Name() )
						ret = err
						close(dirsToCheck)
						return
					}
				}

				isDir := file.IsDir()
				if isDir {
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

func Logging(path string, root string, dir string, filen string) error {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		log.Println("Не могу создать файл для лога", err)
		return err
	}
	st := fmt.Sprintf("Директория где обраатывается - %s, еще папка - %s, файл - %s ", root, dir, filen)
	_, _ = file.Write([]byte(st))
	return nil
}
