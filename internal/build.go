package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type SourceFile struct {
	Root string
	Path string
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
	// create ast
	// check / validate?
	// rewrite (desugar)
	// IR
	// generate code
	// compile the code

	for _, sourceFile := range input {
		fmt.Println()
		logger.Printf("Processing: %s\n", sourceFile.Path)
		content, e := os.ReadFile(sourceFile.Path)

		fileName, _ := strings.CutPrefix(sourceFile.Path, sourceFile.Root)
		tokens, e := Tokenize(fileName, string(content))
		a, e := Parse(fileName, tokens)
		if e != nil {
			logger.Fatal(e)
		}
		// ir
		e = GenerateCode(a)
		if e != nil {
			log.Fatalf("%q", e)
			return
		}

	}
	logger.Printf("go build\n")
	gobuild()
}

func gobuild() {
	cmd := exec.Command("go", "build", "-C", "generated/", "-o", "i-was-generated", "main.go")
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))
}
