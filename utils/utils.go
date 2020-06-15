package utils

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	PRO_NAME = "tomm"
)

func GetProDirAbs() string {
	sbuilder := strings.Builder{}
	sbuilder.WriteString(os.Getenv("GOPATH"))
	sbuilder.WriteString(string(filepath.Separator))
	sbuilder.WriteString("src")
	sbuilder.WriteString(string(filepath.Separator))
	sbuilder.WriteString(PRO_NAME)
	sbuilder.WriteString(string(filepath.Separator))

	path := sbuilder.String()
	return path
}
