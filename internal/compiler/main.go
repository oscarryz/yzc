package compiler

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var logger = log.Default()

func Build(input []SourceFile) {

	for _, sourceFile := range input {
		content, err := os.ReadFile(sourceFile.path)
		if err != nil {
			logger.Fatalf("%q", err)

		}
		d := "generated/" + filepath.Dir(sourceFile.path)
		logger.Printf("About to create: %s", d)
		err = os.MkdirAll(d, 0750)
		if err != nil {
			logger.Fatalf("%q", err)
		}
		if err := os.WriteFile("generated/"+sourceFile.path, content, 0750); err != nil {
			fmt.Printf("write error: %q", err)
			return
		}

	}
	cmd := exec.Command("go", "build", "generated")
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
	}
	return

}

type SourceFile struct {
	root string
	path string
}

func NewSourceFile(root string, path string) SourceFile {
	return SourceFile{
		root, path,
	}
}
