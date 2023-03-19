package aes

import (
	"fmt"
	"testing"
)

func TestAse(t *testing.T) {
	str := []byte("12345678")
	pwd, err := EnPwdCode(str)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pwd)
	bytes, err := DePwdCode(pwd)
	if err != nil {
		t.Error(err)
	}
	if string(bytes) != "12345678" {
		t.Error("加密解密错误")
	}
}
