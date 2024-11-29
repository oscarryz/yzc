package internal

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func GenerateCode(tempDir string, boc *Boc, bocGoName string) (string, error) {
	content, e := Bytes(boc, bocGoName)
	if e != nil {
		logger.Fatalf("generate code error: %q", e)
		return "", e

	}
	fileName := filepath.Join(tempDir, bocGoName)
	if err := os.WriteFile(fileName, content, 0750); err != nil {
		logger.Fatalf("write error: %q", err)
		return "", err
	}
	return fileName, nil
}

func Bytes(boc *Boc, name string) ([]byte, error) {
	goSourceTemplate, err := template.New("main").Parse(
		`package main

func main() {
	print("Hello {{.}} code generator")
}`)
	if err != nil {
		return nil, err
	}
	var sb strings.Builder
	err = goSourceTemplate.Execute(&sb, name)
	if err != nil {
		return nil, err
	}
	return []byte(sb.String()), nil
}
