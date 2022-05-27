package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func SignData(key, data string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	//sign := h.Sum(nil)
	//h := sha256.Sum256([]byte(data + key))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifyHash(hash1, hash2 string) bool {

	//TODO обработать ошибки
	sig1, _ := hex.DecodeString(hash1)
	sig2, _ := hex.DecodeString(hash2)
	return hmac.Equal(sig1, sig2)
}
