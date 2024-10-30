package internal

import (
	"fmt"
	"os"
)

func GenerateCode(a *boc) error {
	logger.Printf("Generating code for:\n %s", a)
	d := "generated/"
	err := os.MkdirAll(d, 0750)
	if err != nil {
		logger.Fatalf("%q", err)
	}
	content := a.Bytes()
	fileName := fmt.Sprintf("%s%s.go", d, a.name)
	if err := os.WriteFile(fileName, content, 0750); err != nil {
		logger.Printf("write error: %q", err)
		return err
	}
	return nil
}
