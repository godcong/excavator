package excavator

import (
	"crypto/md5"
	"fmt"
	"os"
)

func CheckExist(s string) bool {
	_, e := os.Open(s)
	return os.IsExist(e)
}

func MD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
