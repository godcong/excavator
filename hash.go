package excavator

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
)

var tmpDir = "tmp"

func CheckExist(s string) bool {
	_, e := os.Open(GetPath(s))
	return e == nil || os.IsExist(e)
}

func GetPath(s string) string {
	return filepath.Join(tmpDir, s)
}

func MD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
