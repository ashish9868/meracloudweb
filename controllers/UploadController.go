package controllers

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ashish9868/meracloud/lib"
	"github.com/ashish9868/meracloud/models"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type UploadController struct{}

func (u *UploadController) Upload(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	folder := user.Username
	c.Request.ParseMultipartForm(1024 * 1024 * 1024 * 2) // 1gb
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"file": "file is required",
		})
		return
	}

	filename := file.Filename
	upload := models.Upload{}
	upload.Extension = string(file.Filename)[strings.Index(file.Filename, ".")+1 : len(file.Filename)]
	upload.Filename = filename
	uploadPath, _ := upload.GetUploadPath(folder)
	upload.Path = uploadPath
	upload.Size = int(file.Size)
	upload.User = &user

	error := validation.ValidateStruct(&upload,
		validation.Field(&upload.Size, validation.By(func(value interface{}) error {
			maxFileGB := 10
			maxUploadSize := 1024 * 1024 * 1024 * maxFileGB // maxFileMb mb
			if upload.Size > maxUploadSize {
				return errors.New("Maximum upload size is " + (strconv.Itoa(maxFileGB)) + " GB")
			}
			if upload.Size < 1 {
				return errors.New("Cannot upload empty file")
			}
			return nil
		})),
	)

	if error != nil {
		c.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	// all ok
	uploadError := c.SaveUploadedFile(file, upload.Path)

	if uploadError != nil {
		c.JSON(http.StatusInternalServerError, uploadError)
		return
	}

	result := lib.DbInstance().Save(&upload)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, result.Error)
		return
	}
	c.JSON(200, upload)
}

func (u UploadController) GetFileInfo(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	id := c.Param("id")

	upload := models.Upload{}

	result := lib.DbInstance().First(&upload, "id = ? and user_id = ?", id, user.ID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, "file you are looking for not found.")
		return
	}

	c.JSON(http.StatusOK, models.UploadWithLinks{
		Upload: upload,
		Url:    upload.GetUrl(),
	})
}

func (u UploadController) DownloadFile(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	id := c.Param("id")

	upload := models.Upload{}

	result := lib.DbInstance().First(&upload, "id = ? and user_id = ?", id, user.ID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, "file you are looking for not found.")
		return
	}

	u.serveFile(upload.Path, c)
}

func (u *UploadController) serveFile(filePath string, c *gin.Context) {
	_, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	// fileI, _ := file.Stat()
	// mtype, _ := mimetype.DetectFile(filePath)
	// defer file.Close()
	// cacheHour := 12 // 12 hours
	// ts := time.Now().Local().Add(time.Hour*time.Duration(cacheHour)).UTC().Format("Mon 02 Jan 2006 15:04:05") + " GMT"
	// extraHeaders := map[string]string{
	// 	"Content-type":  mtype.String(),
	// 	"Pragma":        "cache",
	// 	"Expires":       ts,
	// 	"Cache-Control": "max-age=" + strconv.Itoa(3600*cacheHour),
	// }
	// c.DataFromReader(http.StatusOK, fileI.Size(), mtype.String(), file, extraHeaders)
	c.File(filePath)
	return
}
