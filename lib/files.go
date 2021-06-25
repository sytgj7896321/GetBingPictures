package mylib

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
