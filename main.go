package main

import (
	"fmt"
	"os"
)

type Bookmark struct {
	Name string
	Url string
	Tags []string
}

func main() {
	argv := os.Args;	

	err := parseFlags(argv);

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err);
		os.Exit(1);
	}
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
