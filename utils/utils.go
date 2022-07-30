package utils

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const UPLOAD_DIR = "uploads"

func RandomBase64Token() string {
	const (
		iterations = 3
	)
	tokens := [iterations]string{}
	for i := 0; i < iterations; i++ {
		token := uuid.NewString() + uuid.NewString() + strconv.FormatInt(time.Now().UnixNano(), 10)
		hash := sha512.Sum512([]byte(token))
		tokens[i] = base64.StdEncoding.EncodeToString(hash[:])
	}
	token := strings.Join(tokens[:], "")
	return base64.StdEncoding.EncodeToString([]byte(token))
}

func CreateSha512Hash(data string) string {
	hash := sha512.Sum512([]byte(data))
	return hex.EncodeToString(hash[:])
}

func RandomUUID() string {
	return uuid.NewString()
}

func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func HasKey(dict map[string]interface{}, key string) bool {
	if _, ok := dict[key]; ok {
		return true
	}
	return false
}

func JsonEncode(obj interface{}) string {
	data, err := json.Marshal(obj)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func ParseInt(s string) int {
	value, error := strconv.Atoi(s)
	if error != nil {
		return 0
	}
	return value
}

func FloatToString(i float64) string {
	return strconv.FormatFloat(i, 'f', -1, 64)
}

func GetBaseUrl(c *gin.Context) string {
	scheme := "http://"
	if !IsEmpty(c.Request.URL.Scheme) {
		scheme = c.Request.URL.Scheme + "://"
	}
	return scheme + c.Request.Host
}
