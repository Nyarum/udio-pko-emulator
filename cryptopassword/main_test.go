package cryptopassword

import (
	"crypto/md5"
	"encoding/hex"
	"testing"
)

func TestEncryptPassword(t *testing.T) {
	md5Init := md5.Sum([]byte("testtest"))
	passMd5 := hex.EncodeToString(md5Init[:])
	_, err := EncryptPassword(string([]byte(passMd5)[:24]), "[07-03 00:36:03:588]")

	if err != nil {
		t.Error("Error encrypt password, err - ", err)
	}
}
