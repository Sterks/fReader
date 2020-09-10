package models

import (
	"github.com/Sterks/Pp.Common.Db/models"
	"os"
)

// InfornationFile Инфромация о файле
type InformationFile struct {
	Inform os.FileInfo
	Fullpath string
	Hash string
	Raw []byte
	Region models.SourceRegions
	TypeFile string
}
