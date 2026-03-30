package main

import (
	"fmt"
	"os"
)

type Group struct {
	Name string
	Bookmarks []Bookmark
}

type Bookmark struct {
	Name string
	Url string
}

func main() {
	argv := os.Args;	

	parseFlags(argv);
}

func readFile(path string) (string, error) {
	dat, err := os.ReadFile(path);
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil;
		}
		return "", fmt.Errorf("failed to read the file: %v", err);
	}

	return string(dat), nil;
}
