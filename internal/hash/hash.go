package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func SignData(key, data string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
