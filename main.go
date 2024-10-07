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

func getLevelsBeforeNoLastFolder(path string) int {
	pathS := strings.Split(path, "/")

	var name string
	var paths []string

	for i := 0; i < len(pathS); i++ {
		name = name + pathS[i] + "/"
		paths = append(paths, name)
	}

	var counter = 1
	for i := len(paths) - 1; i >= 0; i-- {
		if strings.Count(paths[i], "/") == 1 {
			continue
		}

		x := 1
		if len(paths) != 1 {
			x = i - 1
		}

		dirs, _ := combineFilesDirs(paths[x])

		islastDir := dirs[len(dirs)-1].Name() == pathS[i] || (len(dirs) == 1 && dirs[0].Name() == pathS[i])

		if strings.Count(paths[i], "/") > 1 && islastDir {
			counter++
		}

	}

	if counter >= getLevel(path) {
		return 0
	}
	return getLevel(path) - counter
}

func combineFilesDirs(path string) ([]os.DirEntry, []os.DirEntry) {
	var dirs []os.DirEntry
	var files []os.DirEntry

	entries, _ := os.ReadDir(path)

	for _, ent := range entries {
		if ent.IsDir() {
			dirs = append(dirs, ent)
		}
		files = append(files, ent)

	}

	return dirs, files

}

func getLevel(path string) int {
	return strings.Count(path, "/")
}

func isBaseFolderLast(path string, printFiles bool) bool {
	if !strings.Contains(path, "/") {
		return false
	}

	basePath := strings.SplitAfterN(path, "/", 100)[0]
	basePath = strings.TrimSuffix(basePath, "/")

	subBasePath := strings.SplitAfterN(path, "/", 100)[1]
	subBasePath = strings.TrimSuffix(subBasePath, "/")

	dirs, files := combineFilesDirs(basePath)
	if printFiles {
		return subBasePath == files[len(files)-1].Name()
	}
	return subBasePath == dirs[len(dirs)-1].Name()
}

func isPrevFolderLast(path string) bool {
	var dirs []os.DirEntry
	var files []os.DirEntry
	var cuttedPath, upFolder string

	counteSplitters := strings.Count(path, "/")
	if counteSplitters == 0 {
		cuttedPath = path
	} else if counteSplitters == 1 {
		cuttedPath = strings.SplitAfterN(path, "/", counteSplitters)[0]
	} else {

		cuttedPathParts := strings.Split(path, "/")
		upFolder = cuttedPathParts[len(cuttedPathParts)-1]
		for i := 0; i <= counteSplitters-1; i++ {
			cuttedPath += cuttedPathParts[i] + "/"
		}

	}

	dirs, files = combineFilesDirs(cuttedPath)

	if len(dirs) == 1 {
		return true
	} else if upFolder == files[len(files)-1].Name() {
		return true
	} else if len(dirs) > 0 && upFolder == dirs[len(dirs)-1].Name() {
		return true
	} else {
		return false
	}

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

		prevFolder := isPrevFolderLast(path)
		noLastFolderLevel := getLevelsBeforeNoLastFolder(path)

		if tabsCount != 0 && !prevFolder {
			fmt.Fprintf(out, "%s", strings.Repeat("│\t", tabsCount))
		} else if isBaseFolderLast(path, printFiles) {
			fmt.Fprintf(out, "%s", strings.Repeat("\t", tabsCount))
		} else if tabsCount != 0 {
			if noLastFolderLevel != 0 {
				fmt.Fprintf(out, "%s%s", strings.Repeat("│\t│", noLastFolderLevel),
					strings.Repeat("\t", tabsCount-noLastFolderLevel))
			} else {

				fmt.Fprintf(out, "│%s", strings.Repeat("\t", tabsCount))
			}
		}

		if !file.IsDir() && printFiles {
			if pos == len(files)-1 {
				fmt.Fprint(out, lastfix)
			} else {
				fmt.Fprint(out, prefix)
			}

			fileSize := strconv.FormatInt(fileInfo.Size(), 10) + "b"
			if fileSize == "0b" {
				fileSize = "empty"
			}
			fmt.Fprintf(out, "%s (%s)\n", file.Name(), fileSize)
		}

		if file.IsDir() {

			if file.Name() == allFiles[len(allFiles)-1].Name() || pos == len(entries)-1 {
				fmt.Fprintf(out, "%s%s\n", lastfix, file.Name())
			} else {

				fmt.Fprintf(out, "%s%s\n", prefix, file.Name())
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
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
