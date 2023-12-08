package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"yzc/internal/compiler"
)

var logger = log.Default()

func main() {
	files := collectSourceFiles("cmd", "internal")
	fmt.Printf("%v", files)
	fmt.Println("Hi from main")
	compiler.Run(`
	package main
	import "fmt"
    func main() {
       fmt.Println("Hello from test")
	}`)
}

func collectSourceFiles(sourceRoots ...string) []string {
	var files []string
	for _, sourceRoot := range sourceRoots {
		logger.Printf("Walking: %s", sourceRoot)
		_ = filepath.WalkDir(sourceRoot, func(path string, info fs.DirEntry, err error) error {
			if strings.HasSuffix(path, ".go") {
				files = append(files, path)
			}
			return nil
		})
	}
	return files

}
