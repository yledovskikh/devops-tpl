package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func SignData(key, data string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data + key))
	sign := h.Sum(nil)
	return hex.EncodeToString(sign)
}

func VerifyHash(hash1, hash2 string) bool {

	//TODO обработать ошибки
	sig1, _ := hex.DecodeString(hash1)
	sig2, _ := hex.DecodeString(hash2)
	return hmac.Equal(sig1, sig2)
}
