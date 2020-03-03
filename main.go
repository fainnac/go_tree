package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func tabs(deep int, tabsType []bool) string {
	var res string
	for i := 0; i < deep; i++ {
		if tabsType[i] {
			res += "    "
		} else {
			res += "│" + "   "
		}
	}

	return res
}

func Size(size int64) string {
	if size != 0 {
		return fmt.Sprintf(" (%db)", size)
	} else {
		return " (empty)"
	}
}

func myReadDirWithFiles(path string, out *os.File, deep int, tabsType *[]bool) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	if deep >= len(*tabsType) {
		*tabsType = append(*tabsType, false)
	}

	for i, file := range files {
		if file.IsDir() {
			if (i + 1) == len(files) {
				(*tabsType)[deep] = true
				out.WriteString(tabs(deep, *tabsType) + "└── " + file.Name() + "/\n")
				myReadDirWithFiles(path+"/"+file.Name(), out, deep+1, tabsType)
			} else {
				(*tabsType)[deep] = false
				out.WriteString(tabs(deep, *tabsType) + "├── " + file.Name() + "/\n")
				myReadDirWithFiles(path+"/"+file.Name(), out, deep+1, tabsType)
			}
		} else {
			if (i + 1) == len(files) {
				(*tabsType)[deep] = true
				out.WriteString(tabs(deep, *tabsType) + "└── " + file.Name() + Size(file.Size()) + "\n")
			} else {
				(*tabsType)[deep] = false
				out.WriteString(tabs(deep, *tabsType) + "├── " + file.Name() + Size(file.Size()) + "\n")
			}
		}
	}
}

func myReadDir(path string, out *os.File, deep int, tabsType *[]bool) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	dirsIndex := make([]os.FileInfo, 0)
	for _, v := range files {
		if v.IsDir() {
			dirsIndex = append(dirsIndex, v)
		}
	}

	if deep >= len(*tabsType) {
		*tabsType = append(*tabsType, false)
	}

	for i, file := range dirsIndex {
		if (i + 1) == len(dirsIndex) {
			(*tabsType)[deep] = true
			out.WriteString(tabs(deep, *tabsType) + "└── " + file.Name() + "/\n")
			myReadDir(path+"/"+file.Name(), out, deep+1, tabsType)
		} else {
			(*tabsType)[deep] = false
			out.WriteString(tabs(deep, *tabsType) + "├── " + file.Name() + "/\n")
			myReadDir(path+"/"+file.Name(), out, deep+1, tabsType)
		}
	}
}

func dirTree(out *os.File, path string, printFiles bool) error {
	start := strings.Replace(path+"/\n", "//", "/", 1)
	out.WriteString(start)

	var tabsType = make([]bool, 0)
	if printFiles {
		myReadDirWithFiles(path, out, 0, &tabsType)
	} else {
		myReadDir(path, out, 0, &tabsType)
	}
	return nil
}

func main() {
	var path string

	out := os.Stdout
	defer out.Close()

	printFiles := flag.Bool("f", false, "Print files")
	flag.Parse()
	args := flag.Args()

	if len(args) > 1 {
		out.WriteString("Too many path's...\n")
		return
	} else if len(args) == 1 {
		path = args[0]
	} else {
		path = "."
	}

	err := dirTree(out, path, *printFiles)
	if err != nil {
		panic(err.Error())
	}
}
