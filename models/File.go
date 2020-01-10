package model

import (
	"database/sql"
	"time"
)

//File - структура пользователя
type File struct {
	TID                   int            `json:"TID" db:"f_id"`
	TParent               sql.NullInt32  `json:"TParent" db:"f_parent"`
	TName                 string         `json:"TName" db:"f_name"`
	TArea                 sql.NullString `json:"TArea" db:"f_area"`
	TType                 int            `json:"TType" db:"f_type"`
	THash                 string         `json:"THash" db:"f_hash"`
	TSize                 int            `json:"TSize" db:"f_size"`
	TDateCreate           time.Time      `json:"TDateCreate" db:"f_date_create"`
	TDateCreateFromSource time.Time      `json:"TDateCreateFromSource" db:"f_date_create_from_source"`
	TFullpath             string         `json:"TFullpath" db:"f_fullpath"`
	TDateLastCheck        time.Time      `json:"TDateLastCheck" db:"f_date_last_check"`
	TFileIsDir            sql.NullString `json:"TFileIsDir" db:"f_file_is_dir"`
}
