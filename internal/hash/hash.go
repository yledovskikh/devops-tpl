package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func SignData(key, data string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	sign := h.Sum(nil)
	return hex.EncodeToString(sign)
}

func HashVerify(hash1, hash2 string) bool {
	return hmac.Equal([]byte(hash1), []byte(hash2))
}
