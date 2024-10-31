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
	//files := collectSourceFiles("phantom")
	//files := collectSourceFiles("README.md")
	//files := collectSourceFiles( ".", "examples/simple")
	//files := collectSourceFiles(".")
	files := collectSourceFiles("examples/simple")
	logger.Printf("Collecting source files:\n")
	for _, f := range files {
		logger.Printf("%v", f)
	}
	internal.Build(files, false)

}

// collectSourceFiles walks through the provided source directories and collects all source files
// with the specified suffix. It returns a slice of SourceFile structs representing the collected files.
// If any errors occur during the directory walk, the function logs the error and terminates the program.
//
// The following validations are performed:
// - The sourceRoots are valid directories.
// - The sourceRoots don't contain duplicate source files (e.g. a source directory is not a subdirectory of another source directory).
//
// Parameters:
// - sourceRoots: variadic string arguments representing the root directories to search for source files.
//
// Returns:
// - []internal.SourceFile: a slice of SourceFile structs representing the collected source files.
func collectSourceFiles(sourceRoots ...string) []internal.SourceFile {
	var files []internal.SourceFile
	seen := make(map[string]internal.SourceFile)
	for _, currentRoot := range sourceRoots {
		//logger.Printf("Walking source directory: %s", currentRoot)
		walkError := filepath.WalkDir(currentRoot, func(path string, info fs.DirEntry, err error) error {

			if err != nil {
				return fmt.Errorf("Reading source directory, %v\n"+
					"Hint: Check all the directories in the source path exists. Source path: %s", err, sourceRoots)
			}
			if info == nil || path == currentRoot && !info.IsDir() {
				return fmt.Errorf("Not a directory: %s\n"+
					"Hint: Check all the directories in the source path exists. Source path: %s", path, sourceRoots)

			}

			if strings.HasSuffix(path, sourceSuffix) && !info.IsDir() {

				if strings.HasPrefix(currentRoot, "./") {
					currentRoot = currentRoot[2:]
				}
				if strings.HasSuffix(currentRoot, "/") {
					currentRoot = currentRoot[:len(currentRoot)-1]
				}
				afp, _ := filepath.Abs(path)

				path, _ = strings.CutPrefix(path, currentRoot)
				if strings.HasPrefix(path, "/") {
					path = path[1:]
				}

				file := internal.NewSourceFile(currentRoot, path, afp)
				if seen[afp] == (internal.SourceFile{}) {
					seen[afp] = internal.NewSourceFile(currentRoot, path, afp)
				} else {
					first := seen[afp]
					return fmt.Errorf("Duplicate source files\n%s (source directory:\"%s\") and %s (souce directory:\"%s\") are the same file: %s\n"+
						"Hint: Check a source directory is not a subdirectory of another source directory. Source directories: %s ", first.Path, first.Root, file.Path, file.Root, afp, sourceRoots)
				}
				files = append(files, file)
			}
			return nil
		})
		if walkError != nil {
			logger.Fatalf("%v", walkError)
		}
	}
	return files
}
