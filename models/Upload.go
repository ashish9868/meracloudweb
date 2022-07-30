package models

import (
	"os"
	"strconv"
	"time"

	"github.com/ashish9868/meracloud/utils"
)

type Upload struct {
	BaseModel
	Filename      string    `json:"filename" gorm:"type:varchar(512);NOT NULL"`
	Path          string    `json:"path" gorm:"type:varchar(512)"`
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

type UploadWithLinks struct {
	Upload
	Url string `json:"url"`
}

func (u *Upload) GetUploadPath(baseFolder string) (string, bool) {

	path := u.Path
	if utils.IsEmpty(u.Path) {
		path = utils.UPLOAD_DIR + "/" + baseFolder + "/" + u.Filename
		if !utils.IsEmpty(baseFolder) {
			u.CreateFolder(baseFolder)
			path = utils.UPLOAD_DIR + "/" + baseFolder + "/" + u.Filename
		}
	}
	exists := false
	_, err := os.Stat(path)
	if err == nil {
		exists = true
	}
	return path, exists
}

func (u *Upload) CreateFolder(name string) bool {
	error := os.MkdirAll(u.GetPathToUploadDir(name), 0755)
	return error == nil
}

func (u *Upload) GetPathToUploadDir(foldername string) string {
	return utils.UPLOAD_DIR + "/" + foldername
}

func (u *Upload) GetUrl() string {
	return "/download/" + strconv.FormatUint(uint64(u.ID), 10)
}
