package utils

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func MustBeNilErr(err error, message string, args ...interface{}) {
	if err != nil {
		logrus.Fatalf(message+": %s", err, args)
	}
}

func MkDirFatal(path string) {
	err := os.MkdirAll(path, 0770)
	MustBeNilErr(err, "failed to create dir %s", path)
}

func JoinFilePath(elements ...string) string {
	return filepath.Join(elements...)
}

func SanitizePath(filePath string) string {
	filePath = strings.ReplaceAll(filePath, "/", "_")
	return strings.Trim(filePath, " ")
}
