package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type SourceFile struct {
	root string
	path string
}

func NewSourceFile(root string, path string) SourceFile {
	return SourceFile{
		root, path,
	}
}

var logger = log.Default()

func Build(input []SourceFile) {
	// read source file
	// tokenize
	// create ast -> validate?
	// second pass
	// generate code
	// compile the code

	for _, sourceFile := range input {
		content, e := os.ReadFile(sourceFile.path)
		tokens := tokenize(e, content)
		a, e := parse(tokens)
		// ir
		e = generateCode(a)
		if e != nil {
			log.Fatalf("%q", e)
			return
		}

	}
	gobuild()
}

func gobuild() {
	cmd := exec.Command("go", "build", "-C", "generated/", "-o", "i-was-generated", "main.go")
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))
}
