package domain

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// fbc1fd86ab53d52c3ffeb6529aea9676e14bc52b792414c32f5612b4eb2c9745:1567618546
// JSv5M7FZ5iPHnHiLXR1QbhnMcdoY/wvEae4a76KrGBxeHruFb1S90d4GkwsoQQU4R1zqEdSa0KMGflriF2dHj5XWm4Zp6OBxLp6BJFUhqpQxEBEr5yl4sxEHadgssvVWfWtDKe0bENU=
// JSv5M7FZ5iPHnHiLXR1QbhnMcDAU/wvEae4a76KrGBxeHruFb1S90d4GkwsoQQU4R1zqEdSa0KMGflriF2dHj5XWm4Zp6OBxLp6BJFUhqpQxEBEr5yl4sxEHadgssvVWfWtDKe0bENU=
// JSv5M7FZ5iPHnHiLXR1QbhnMcNEF/wvEae4a76KrGBxeHruFb1S90d4GkwsoQQU4R1zqEdSa0KMGflriF2dHj5XWm4Zp6OBxLp6BJFUhqpQxEBEr5yl4sxEHadgssvVWfWtDKe0bENU=

type HashToken struct {
	Secret []byte
}

func NewHMACHashToken(secret string) (*HashToken, error) {
	return &HashToken{Secret: []byte(secret)}, nil
}

func (tk *HashToken) Create(uString string, tokenExpTime int64) (string, error) {
	h := hmac.New(sha256.New, []byte(tk.Secret))
	data := fmt.Sprintf("%s:%d", uString, tokenExpTime)
	h.Write([]byte(data))
	token := hex.EncodeToString(h.Sum(nil)) + ":" + strconv.FormatInt(tokenExpTime, 10)
	return token, nil
}

func (tk *HashToken) Check(headerToken string, cookieToken string) (bool, error) {
	headerData, tokenExp, err := splitToken(headerToken)
	if err != nil {
		return false, err
	}

	headerMAC, err := checkTokenData(headerData, tokenExp)
	if err != nil {
		return false, err
	}

	cookieData, tokenExp, err := splitToken(headerToken)
	if err != nil {
		return false, err
	}

	cookieMAC, err := checkTokenData(cookieData, tokenExp)
	if err != nil {
		return false, err
	}

	return hmac.Equal(headerMAC, cookieMAC), nil
}

func splitToken(token string) (string, int64, error) {
	tokenData := strings.Split(token, ":")
	if len(tokenData) != 2 {
		return "", 0, fmt.Errorf("bad token data")
	}

	tokenExp, err := strconv.ParseInt(tokenData[1], 10, 64)
	if err != nil {
		return "", 0, err
	}

	return tokenData[0], tokenExp, nil
}

func checkTokenData(data string, exp int64) ([]byte, error) {
	if exp < time.Now().Unix() {
		return []byte{}, fmt.Errorf("token expired")
	}

	mac, err := hex.DecodeString(data)
	if err != nil {
		return []byte{}, err
	}

	return mac, nil
}
