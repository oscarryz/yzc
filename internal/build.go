package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type SourceFile struct {
	Root         string
	Path         string
	AbsolutePath string
}

func NewSourceFile(root, path, absolutePath string) SourceFile {
	return SourceFile{
		root, path, absolutePath,
	}
}

var logger = log.Default()

const (
	target_dir = "target/"
)

func Build(input []SourceFile, keepGeneratedSource bool) {
	// read source file
	// tokenize
	// create ast
	// check / validate?
	// rewrite (desugar)
	// IR
	// generate code
	// compile the code

	tmpDir, cleanup := createTempDir(keepGeneratedSource)
	defer cleanup()

	for _, sourceFile := range input {
		fmt.Println()
		logger.Printf("Processing: %s\n", sourceFile.AbsolutePath)
		content, e := os.ReadFile(sourceFile.AbsolutePath)

		tokens, e := Tokenize(sourceFile.Path, string(content))
		boc, e := Parse(sourceFile.Path, tokens)
		if e != nil {
			logger.Fatal(e)
		}
		// ir
		logger.Printf("IR: %v\n", boc)

		// generate code
		fileName, e := GenerateCode(tmpDir, boc)

		if e != nil {
			log.Fatalf("%q", e)
			return
		}
		logger.Printf("go build %s\n", fileName)
		// compile the code
		gobuild(boc.Name, fileName)

	}
}

func gobuild(name, fileName string) {
	_ = os.MkdirAll(target_dir, 0750)
	outputFile := fmt.Sprintf("%s%s", target_dir, name)
	//logger.Printf("Generated %s", outputFile)

	cmd := exec.Command("go", "build", "-o", outputFile, fileName)
	output, _ := cmd.CombinedOutput()
	if len(output) > 0 {
		logger.Println(string(output))
	}
}

// createTempDir creates a temporary directory for generated source files.
// It returns the path to the directory and a cleanup function that removes the directory.
// If keepGeneratedSource is true, the source is created under the ./generated and the cleanup function is a no-op.
func createTempDir(keepGeneratedSource bool) (string, func()) {
	generatedDir := ""
	if keepGeneratedSource {
		generatedDir = "generated"
		_ = os.MkdirAll(generatedDir, 0750)
	}
	tmpDir, e := os.MkdirTemp(generatedDir, "yzc_generated_go")

	if e != nil {
		logger.Fatalf("%q", e)
	}

	cleanup := func() {
		if keepGeneratedSource {
			logger.Printf("Generated source files are in %s\n", tmpDir)
		} else {
			_ = os.RemoveAll(tmpDir)
		}
	}

	return tmpDir, cleanup
}
