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
	// todo: handle nested source dirs, e.g. below `examples` is inside `.` so, ['.' './examples'] would duplicate
	files := collectSourceFiles("examples/simple", "internal")
	logger.Printf("%v\n", files)
	internal.Build(files)

}

func collectSourceFiles(sourceRoots ...string) []internal.SourceFile {
	var files []internal.SourceFile
	for _, sourceRoot := range sourceRoots {
		logger.Printf("Walking: %s", sourceRoot)
		_ = filepath.WalkDir(sourceRoot, func(path string, info fs.DirEntry, err error) error {

			if strings.HasSuffix(path, sourceSuffix) && !info.IsDir() {
				files = append(files, internal.NewSourceFile(sourceRoot, path))
			}
			return nil
		})
	}
	return files
}
