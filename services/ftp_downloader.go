package services

import (
	"github.com/Sterks/Pp.Common.Db/models"
	"github.com/Sterks/fReader/config"
	model "github.com/Sterks/fReader/pkg/models"
	"github.com/secsy/goftp"
)

// DownloaderFtp ...
type DownloaderFtp interface {
	AddTimeNow(config *config.Config)
	Connect(config *config.Config) *goftp.Client
	Start(config *config.Config)
	GetListFolder()
	GetAllFolderRegionsDb() []models.SourceRegions
	GetFileInfo([]models.SourceRegions, string) []model.InformationFile
	CheckDownloder([]model.InformationFile) []model.InformationFile
}

// StartService Interface
func StartService(down DownloaderFtp, config *config.Config, typeFile string) {
	down.AddTimeNow(config)
	down.Connect(config)
	down.Start(config)
	down.GetListFolder()
	listRegions := down.GetAllFolderRegionsDb()
	listFilesNotification := down.GetFileInfo(listRegions, typeFile)
	_ = down.CheckDownloder(listFilesNotification)
}

// StartService223 Interface
func StartService223(down DownloaderFtp, config *config.Config, typeFile string) {
	down.AddTimeNow(config)
	down.Connect(config)
	down.Start(config)
	down.GetListFolder()
	listRegions := down.GetAllFolderRegionsDb()
	listFilesNotification := down.GetFileInfo(listRegions, typeFile)
	_ = down.CheckDownloder(listFilesNotification)
}
