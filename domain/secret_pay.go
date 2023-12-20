package domain

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

func dbParamsfromEnvUsr() string {
	secret := os.Getenv("SECRET_YOOMONEY")
	return fmt.Sprintf("secret=%s", secret)
}
func CreateParametersString(r *http.Request) string {
	values := r.URL.Query()
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	secret := dbParamsfromEnvUsr()
	var keys []string
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var parameters []string
	for _, key := range keys {
		if key != "sha1_hash" {
			parameters = append(parameters, key+"="+values.Get(key))
		}
	}
	parameters = append(parameters, "notification_secret="+secret)

	return strings.Join(parameters, "&")
}

func CalculateSHA1Hash(parameters string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(parameters))
	if err != nil {
		return "", err
	}
	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)
	return hashString, nil
}
