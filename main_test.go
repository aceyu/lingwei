package main_test

import (
	"encoding/hex"
	"fmt"
	"lingwei/letsgo"
	"testing"
)

func TestGenKey(t *testing.T) {
	hn := "DNhM2MzDQz"
	aeskey := "HIgtcdRUxqT72582"
	baes, _ := letsgo.EncryptAES([]byte(hn), []byte(aeskey))
	ser := hex.EncodeToString(baes)[7:17]
	fmt.Println(ser)
}
