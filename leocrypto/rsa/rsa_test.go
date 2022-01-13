package rsa

import (
	"fmt"
	"testing"
)

const (
	publicKeyPath  string = "client_public.pem"
	privateKeyPath string = "client_private.pem"
)

func TestGenRsaKey(t *testing.T) {
	err := GenRsaKey("", "", 1024*4)
	if err != nil {
		fmt.Println("fail")
	} else {
		fmt.Println("success")
	}
}

func TestRsaEncrypt(t *testing.T) {
	rowData := "hello"
	data, err := RsaEncrypt("./public.pem", []byte(rowData))
	println(data)
	if err != nil {
		t.Fatal(err)
	}
	decryData, err := RsaDecrypt("./private.pem", []byte(data))
	println(string(decryData))
	if err != nil {
		t.Fatal(err)
	}
	if rowData != string(decryData) {
		t.Fail()
	}
}
