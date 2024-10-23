package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"yzc/internal"
)

var logger = log.Default()

const sourceSuffix = ".yz"

func main() {
	fmt.Println()
	files := collectSourceFiles(".", "examples/simple", "internal")
	logger.Printf("Collecting source files:\n")
	for _, f := range files {
		logger.Printf("%v", f)
	}
	internal.Build(files)

}

func collectSourceFiles(sourceRoots ...string) []internal.SourceFile {
	logger.Printf("Source directories: %v\n", sourceRoots)
	var files []internal.SourceFile
	seen := make(map[string]bool)
	for _, sourceRoot := range sourceRoots {
		logger.Printf("Walking: %s", sourceRoot)
		_ = filepath.WalkDir(sourceRoot, func(path string, info fs.DirEntry, err error) error {

			if strings.HasSuffix(path, sourceSuffix) && !info.IsDir() {
				if seen[path] == false {
					seen[path] = true
					files = append(files, internal.NewSourceFile(sourceRoot, path))
				}
			}
			return nil
		})
	}
	return files
}
