package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var ignoreFiles = map[string]interface{}{
	".git":       nil,
	".idea":      nil,
	".gitignore": nil,
}

func getLevel(path string) int {
	return strings.Count(path, "/")
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	prefix := "├───"
	lastfix := "└───"


	var dirs []os.DirEntry
	var files []os.DirEntry

	entries, _ := os.ReadDir(path)
	for _, ent := range entries {
		if ent.IsDir() {
			dirs = append(dirs, ent)
		}
		files = append(files, ent)
	}

	var allFiles []os.DirEntry
	if printFiles {
		allFiles = entries
	} else {
		allFiles = dirs
	}


	tabsCount := getLevel(path)


	for pos, file := range entries {
		fileInfo, _ := file.Info()

		// Доп функция
		if _, ok := ignoreFiles[file.Name()]; ok {
			continue
		} else if !file.IsDir() && !printFiles {
			continue
		}


		if tabsCount != 0 {
			fmt.Fprintf(out, "│%s", strings.Repeat("\t", tabsCount))
		}


		if !file.IsDir() && printFiles {
			if pos == len(files)-1 {
				fmt.Fprintf(out, lastfix)
			} else {
				fmt.Fprintf(out, prefix)
			}

			fileSize := strconv.FormatInt(fileInfo.Size(), 10) + "b"
			if fileSize == "0b" {
				fileSize = "empty"
			}
			fmt.Fprintf(out, "%s (%s)\n", file.Name(), fileSize)
		}

		if file.IsDir(){

			if file.Name() == allFiles[len(allFiles)-1].Name() || pos == len(files)-1 {
				fmt.Fprintf(out,"%s%s\n", lastfix, file.Name())
			} else {
				fmt.Fprintf(out,"%s%s\n",prefix, file.Name())
			}



		}
		subPath := fmt.Sprintf("%s/%s", path, file.Name())
		_ = dirTree(out, subPath, printFiles)

	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	//path := os.Args[1]
	path := "testdata"
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
