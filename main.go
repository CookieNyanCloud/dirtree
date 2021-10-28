package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Node struct {
	path     string
	pref     string
	lastPref string
}

type Stack []Node

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(str Node) {
	*s = append(*s, str)
}

func (s *Stack) Pop() (Node, bool) {
	if s.IsEmpty() {
		return Node{}, false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
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

func dirTree(out *os.File, path string, printFiles bool) error {

	stack := Stack{Node{
		path:     path,
		pref:     "",
		lastPref: "├───",
	}}
	j := 0

	for len(stack) > 0 {
		cur, ok := stack.Pop()
		if !ok {
			fmt.Println(ok)
			continue
		}

		openNode, err := os.Open(cur.path)
		defer openNode.Close()

		if err != nil {
			fmt.Println(err)
			return err
		}
		stat, err := openNode.Stat()
		if err != nil {
			return err
		}
		j++
		if stat.IsDir() {

			names, err := openNode.Readdirnames(0)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if j != 1 {
				_, _ = out.Write([]byte(cur.pref + cur.lastPref + stat.Name() + "\n"))
				j++
			}
			sort.Strings(names)
			sort.Sort(sort.Reverse(sort.StringSlice(names)))
			for i, filename := range names {
				//fmt.Println(names,filename)
				if strings.Contains(filename, ".git") || strings.Contains(filename, ".idea") {
					continue
				}
				var lastPref, pref string
				//fmt.Println(names,filename)
				//if i == len(names)-1 {
				if i == 0 {
					lastPref = "└───"
				} else {
					lastPref = "├───"
				}
				if cur.lastPref == "└───" {
					pref = cur.pref + "\t"
				} else {
					pref = cur.pref + "│\t"
				}
				if j == 1 {
					pref = ""
				}
				stack.Push(Node{
					path:     openNode.Name() + string(filepath.Separator) + filename,
					pref:     pref,
					lastPref: lastPref,
				})
			}

		} else if printFiles {
			sizeInt := stat.Size()
			var size string
			if sizeInt > 0 {
				size = " (" + strconv.FormatInt(sizeInt, 10) + "b)"
			} else {
				size = " (empty)"
			}
			_, _ = out.Write([]byte(cur.pref + cur.lastPref + stat.Name() + size + "\n"))
		}

	}
	return nil
}
