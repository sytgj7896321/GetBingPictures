package myfile

import (
	"fmt"
	"os"
)

func CreatePath(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.Mkdir(path, 0666)
		if err != nil {
			return err
		}
		fmt.Println("Path '" + path + "' created")
		return nil
	}
	return err
}

func OpenFile(name string) (*os.File, error) {
	_, err := os.Stat(name)
	if err == nil {
		fp, err := os.Open(name)
		return fp, err
	}
	if os.IsNotExist(err) {
		fp, err := os.Create(name)
		return fp, err
	}
	return nil, err
}
