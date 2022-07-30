package models

import "time"

type Upload struct {
	BaseModel
	Filename      string    `json:"filename" gorm:"type:varchar(512);NOT NULL"`
	AppPath       string    `json:"appPath"`
	Path          string    `json:"path" gorm:"type:varchar(512)"`
	Mime          string    `json:"mime" gorm:"type:varchar(255)"`
	Size          int       `json:"size" gorm:"type:int(11)"`
	Extension     string    `json:"extension" gorm:"type:char(6)"`
	Hash          string    `json:"-"`
	UserID        uint      `json:"-"`
	OriginalAtime time.Time `json:"originalAtime"`
	OriginalMtime time.Time `json:"originalMtime"`
	LocalAtime    time.Time `json:"localAtime"`
	LocalMtime    time.Time `json:"localMtime"`
	User          *User     `json:"-"`
}
