package common

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Sterks/fReader/config"
	"github.com/sirupsen/logrus"
)

// GenerateID - Герерация строки длинной 12 символов
func GenerateID(ident int) string {
	ident = ident + 1
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
func PathLocalFile(config config.Config, id int) string {
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
