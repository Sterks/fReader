package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	model "github.com/Sterks/FReader/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // ....
)

const (
	host     = "localhost"
	port     = 5432
	user     = "user_ro"
	password = "4r2w3e1q"
	dbname   = "freader"
)

// Store ...
type Store struct {
	db *sqlx.DB
}

// New ...
func New(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

// OpenDb ...
func (s *Store) OpenDb() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("Соединиться не удалось - %s", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("Соединиться не удалось понг - %v", err)
	}
	s.db = db
}

// GetFilesFromDownloader ...
// func (s *Store) GetFilesFromDownloader() []model.File {
// 	// var files []model.File
// 	files := make([]model.File, 0)
// 	rows, err := s.db.Query(`
// 	select
// 	"tF_ID",
// 	"tF_Parent",
// 	"tF_Name",
// 	"tF_Source",
// 	"tF_SourcePath",
// 	"tF_Type",
// 	"tF_SHA256",
// 	"tF_Size",
// 	"tF_DateCreate",
// 	"tF_DateLastCheck",
// 	"tF_DateCreateFromSource",
// 	"tF_DateLastModifyFromSource"
// from
// 	public."tFiles"
// where
// 	"tF_Type" = 10022
// 	and "tF_Source" = 10035
// 	and "tF_Name" ilike '%notification%'
// 	and "tF_DateCreate" > '2019-11-19'
// 	`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for rows.Next() {
// 		var file model.File
// 		err := rows.Scan(
// 			&file.TFID,
// 			&file.TFParent,
// 			&file.TFName,
// 			&file.TFSource,
// 			&file.TFSourcePath,
// 			&file.TFType,
// 			&file.TFSHA256,
// 			&file.TFSize,
// 			&file.TFDateCreate,
// 			&file.TFDateLastCheck,
// 			&file.TFDateCreateFromSource,
// 			&file.TFDateLastModifyFromSource,
// 		)
// 		if err != nil {
// 			log.Printf("Ошибка на маппинге %s", err)
// 		}
// 		files = append(files, file)
// 		if err = rows.Err(); err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	return files
// }

// CheckInfoFile ...
func (s *Store) CheckInfoFile(info os.FileInfo, hash string) int {
	var ident int
	row := s.db.QueryRowx("select f_id from \"Files\" where f_hash = $1 and f_size = $2 and f_name = $3", hash, info.Size(), info.Name())
	row.Scan(&ident)
	// if err != nil {
	// 	log.Printf("Ошибка в методе CheckInfoFile - %v", err)
	// }
	// fmt.Println(ident)
	return ident

}

// CreateInfoFile ...
func (s *Store) CreateInfoFile(info os.FileInfo, region string, Hash string, fullpath string) {

	id := s.CheckInfoFile(info, Hash)
	if id > 0 {
		id := s.CheckInfoFile(info, Hash)
		var f model.File
		t := time.Now().Format("2006-01-02T15:04:05")
		query := fmt.Sprintf("UPDATE public.\"Files\" SET f_date_last_check='%v' where f_id = %v", t, id)
		// fmt.Println(id)
		_, err3 := s.db.Exec(query)
		if err3 != nil {
			log.Fatalf("Обновление не прошло - %v", err3)
		}

		err := s.db.Get(&f, "select * from public.\"Files\" where f_id = $1", id)
		if err != nil {
			log.Printf("Ошибка в функции проверки обновления - %s", err)
		}

		fmt.Printf("Дата последней проверки %v, %v - %v  \n", f.TName, f.TDateCreate, f.TDateLastCheck)
	} else {

		ext := filepath.Ext(info.Name())
		typeFile := s.FindExt(ext)

		query := fmt.Sprintf(`INSERT INTO "Files" 
	(
	 f_parent,
	 f_name,
	 f_area,
	 f_type,
	 f_hash,
	 f_size,
	 f_date_create,
	 f_date_create_from_source,
	 f_fullpath,
	 f_file_is_dir,
	 f_date_last_check)
	 VALUES
	 (NULL, '%v', '%v', %v,
	 '%v', %v, '%v', '%v', '%v',  NULL, '%v')
	`, info.Name(),
			region,
			typeFile,
			Hash,
			info.Size(),
			time.Now().Format("2006-01-02T15:04:05"),
			info.ModTime().Format("2006-01-02T15:04:05"),
			fullpath,
			time.Now().Format("2006-01-02T15:04:05"),
		)
		// fmt.Println(query)
		_, err := s.db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(fullpath)
	}
}

// FindExt ...
func (s *Store) FindExt(ext string) int {
	var ident int
	row := s.db.QueryRowx("select ft_id from \"FilesTypes\" where ft_ext = $1", ext)
	err := row.Scan(&ident)
	if err != nil {
		fmt.Println(err)
	}
	return ident
}
