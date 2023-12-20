package internal

import (
	"fmt"
	"os"
)

func generateCode(a *program) error {
	d := "generated/"
	err := os.MkdirAll(d, 0750)
	if err != nil {
		logger.Fatalf("%q", err)
	}
	content := a.Bytes()
	if err := os.WriteFile("generated/main.go", content, 0750); err != nil {
		fmt.Printf("write error: %q", err)
		return err
	}
	return nil
}
