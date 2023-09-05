package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

func ParseYamlFile(filePath string, obj interface{}) error {
	byteData, err := ioutil.ReadFile(GetAbsDir(filePath))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(byteData, obj)
}
func GetAbsDir(relativePath string) string {
	dir := filepath.Dir(os.Args[0])
	return path.Join(dir, relativePath)
}

func GetLogFileName() string {
	timeNow := time.Now()
	t1, t2, t3, t4, t5 := fmt.Sprint(timeNow.Year()), fmt.Sprint(timeNow.Month()), fmt.Sprint(timeNow.Day()), fmt.Sprint(timeNow.Hour()), fmt.Sprint(timeNow.Minute())
	return strings.Join([]string{t1, t2, t3, t4, t5}, "-") + ".log"
}
