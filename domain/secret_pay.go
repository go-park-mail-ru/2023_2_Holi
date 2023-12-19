package domain

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"sort"
	"strings"
)

const secret = `L9KkXz48oYV2+zpEer9DPy7/`

func CreateParametersString(r *http.Request) string {
	values := r.URL.Query()

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
