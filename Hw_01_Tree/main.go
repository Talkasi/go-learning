package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
)

type node struct {
	name     string
	isDir    bool
	size     int64
	children []*node
}

func main() {
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	includeFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(os.Stdout, path, includeFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, useFiles bool) error {
	var tree node
	if err := readTree(path, &tree, useFiles); err != nil {
		return err
	}
	tree = *tree.children[0]

	sortTree(&tree)
	printTree(out, &tree, "", false, 0)

	return nil
}

func readTree(path string, curNode *node, includeFiles bool) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.New("Unable to read file: " + err.Error())
	}
	defer file.Close()

	fStat, err := file.Stat()
	if err != nil {
		return errors.New("Unable to get file stat: " + err.Error())
	}

	switch fMode := fStat.Mode(); {
	case fMode.IsDir():
		err = os.Chdir(file.Name())
		if err != nil {
			return errors.New("Unable to chdir: " + err.Error())
		}

		newNode := node{
			name:  file.Name(),
			isDir: true,
		}
		curNode.children = append(curNode.children, &newNode)

		entries, err := file.Readdir(-1)
		if err != nil {
			return errors.New("Unable to get dir contents: " + err.Error())
		}

		for i := 0; i < len(entries); i++ {
			readTree(entries[i].Name(), &newNode, includeFiles)
		}
		err = os.Chdir("..")
	default:
		if includeFiles {
			newNode := node{
				name: fStat.Name(),
				size: fStat.Size(),
			}
			curNode.children = append(curNode.children, &newNode)
		}
	}

	return nil
}

func compareStrings(a, b string) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func sortTree(curNode *node) {
	if !curNode.isDir {
		return
	}

	slices.SortFunc(curNode.children, func(a, b *node) int {
		return compareStrings(a.name, b.name)
	})

	for i := 0; i < len(curNode.children); i++ {
		sortTree(curNode.children[i])
	}
}

func getPrefixPart(isLast bool) string {
	if isLast {
		return "└───"
	}

	return "├───"
}

func formatSize(size int64) string {
	if size == 0 {
		return "empty"
	}
	return strconv.FormatInt(size, 10) + "b"
}

func printNode(out io.Writer, curNode *node, prefix string) {
	if curNode.isDir {
		fmt.Fprintf(out, "%s%s\n", prefix, curNode.name)
	} else {
		size := formatSize(curNode.size)
		fmt.Fprintf(out, "%s%s (%s)\n", prefix, curNode.name, size)
	}
}

func printTree(out io.Writer, curNode *node, prefix string, isLast bool, depth int) {
	if depth > 0 {
		printNode(out, curNode, prefix+getPrefixPart(isLast))

		if isLast {
			prefix += "\t"
		} else {
			prefix += "│\t"
		}
	}

	if !curNode.isDir {
		return
	}

	for i := 0; i < len(curNode.children); i++ {
		isLast := i == len(curNode.children)-1
		printTree(out, curNode.children[i], prefix, isLast, depth+1)
	}
}
