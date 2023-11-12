package test

import (
	"fmt"
	"rpc-go-netty/utils/aes"
	"testing"
)

func TestAesEncoder(t *testing.T) {

	cal := "VYGThW0MXPf4v88IKP/o4g=="
	// 揭秘
	decrypt, _ := aes.Decrypt(cal)
	fmt.Println(decrypt)

}
