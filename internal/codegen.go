package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func GenerateCode(tempDir string, boc *boc) (string, error) {
	content := Bytes(boc)
	bocGoName := fmt.Sprintf("%s.go", boc.Name)
	fileName := filepath.Join(tempDir, bocGoName)
	if err := os.WriteFile(fileName, content, 0750); err != nil {
		logger.Fatalf("write error: %q", err)
		return "", err
	}
	return fileName, nil
}

func Bytes(boc *boc) []byte {
	goSourceTemplate, err := template.New("main").Parse(
		`package main

func main() {
	print("Hello {{.Name}} code generator")
}`)
	if err != nil {
		return nil
	}
	var sb strings.Builder
	err = goSourceTemplate.Execute(&sb, boc)
	if err != nil {
		return nil
	}
	return []byte(sb.String())
}
