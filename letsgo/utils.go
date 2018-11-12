package letsgo

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"net"
	"os"
)

func GetHardwareNo() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	var hardwareAddr string
	for _, inter := range interfaces {
		if inter.HardwareAddr != nil {
			hardwareAddr = inter.HardwareAddr.String()
			break
		}
	}
	if hardwareAddr != "" {
		hostname, _ := os.Hostname()
		hn := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString([]byte(hardwareAddr + "313a61643a" + hostname))))
		if len(hn) > 67 {
			hn = hn[5:8] + hn[62:66] + hn[len(hn)-3:len(hn)]
		} else {
			hn = hn[5:12] + hn[len(hn)-3:len(hn)]
		}
		return hn
	}
	return ""
}

func EncryptAES(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func DecryptAES(src []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockmode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(src))
	blockmode.CryptBlocks(origData, src)
	origData = unpadding(origData)
	return origData, nil
}

func padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
