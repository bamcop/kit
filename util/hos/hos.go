package hos

import (
	"os"
)

func MustPwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return pwd
}

func MustUserHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return dir
}

func Exist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
