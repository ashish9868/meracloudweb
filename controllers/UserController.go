package controllers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/ashish9868/meracloud/lib"
	"github.com/ashish9868/meracloud/models"
	"github.com/ashish9868/meracloud/utils"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type RegisterDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserController struct{}

func (u UserController) Register(c *gin.Context) {
	var payload RegisterDto
	c.ShouldBind(&payload)

	err := validation.ValidateStruct(&payload,
		validation.Field(&payload.Username, validation.Required, validation.Length(5, 255)),
		validation.Field(&payload.Password, validation.Required, validation.Length(5, 255)),
		validation.Field(&payload.Username, validation.By(func(value interface{}) error {
			foundUser := models.User{}
			lib.DbInstance().First(&foundUser, "username = ?", payload.Username)

			if foundUser.ID > 0 {
				return errors.New("Username already exists.")
			}
			return nil
		})),
	)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{
		Username: payload.Username,
		Password: hashPassword(payload.Password),
	}
	result := lib.DbInstance().Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"global": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, user)

}

func (u UserController) Login(c *gin.Context) {
	var payload LoginDto
	c.ShouldBind(&payload)

	err := validation.ValidateStruct(&payload,
		validation.Field(&payload.Username, validation.Required, validation.Length(5, 255)),
		validation.Field(&payload.Password, validation.Required, validation.Length(5, 255)),
	)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	result := lib.DbInstance().First(&user, "username = ? and password = ?", payload.Username, hashPassword(payload.Password))

	if !(user.ID > 0) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"global": "Invalid username and/or password."})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"global": "Invalid username and/or password."})
		return
	}

	user.Token = utils.RandomBase64Token()
	user.TokenExpiry = time.Now().Add(time.Hour * time.Duration(24))

	result = lib.DbInstance().Save(user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"global": "Server error."})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u UserController) Logout(c *gin.Context) {
	user, ok := c.Get("user")

	if ok {
		userObj := user.(models.User)
		userObj.Token = ""
		userObj.TokenExpiry = time.Now().Add(time.Hour * time.Duration(-48))
		lib.DbInstance().Save(userObj)
	}
	c.JSON(200, gin.H{
		"success": ok,
	})
}

func hashPassword(password string) string {
	return utils.CreateSha512Hash(os.Getenv("SALT") + password + os.Getenv("SALT"))
}
