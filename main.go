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

var lvl = -1

func dirTree(out io.Writer, path string, printFiles bool) error {

	lvl++
	files, _ := os.ReadDir(path)

	for pos, file := range files {
		fileInfo, _ := file.Info()

		// Доп функция
		if _, ok := ignoreFiles[file.Name()]; ok {
			continue
		} else if !file.IsDir() && !printFiles {
			continue
		}

		if lvl != 0 {
			fmt.Fprintf(out, "%s", strings.Repeat("│\t", lvl))
		}

		if !file.IsDir() && printFiles {
			if pos == len(files)-1 {
				fmt.Fprintf(out, "└───")
			} else {
				fmt.Fprintf(out, "├───")
			}

			fileSize := strconv.FormatInt(fileInfo.Size(), 10) + "b"
			if fileSize == "0b" {
				fileSize = "empty"
			}
			fmt.Fprintf(out, "%s (%s)\n", file.Name(), fileSize)
		}

		if file.IsDir() {
			if file.Name() == files[len(files)-2].Name() || pos == len(files)-1 {
				fmt.Fprintf(out, "└───")
			} else {
				fmt.Fprintf(out, "├───")
			}

			fmt.Fprintf(out, "%s\n", file.Name())

			subPath := fmt.Sprintf("%s/%s", path, file.Name())
			_ = dirTree(out, subPath, printFiles)
		}

	}
	lvl--
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
