package test

import (
	"fmt"
	"rpc-go-netty/utils/aes"
	"testing"
)

func TestAesEncoder(t *testing.T) {

	cal := "o53cnwtKwHd8XT8hjohYCsWJl6dhySqA38mIXxx/4EXrQV6j6ZqfdLxR2cZ1W5Waqpxn2dj2ynm2/7lGfJSzMA=="
	// 揭秘
	decrypt, _ := aes.Decrypt(cal)
	fmt.Println(decrypt)

}
