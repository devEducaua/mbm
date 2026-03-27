package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	argv := os.Args;	

	if len(argv) < 2 {
		fmt.Fprintf(os.Stderr, "bk: command not passed\n")
		os.Exit(1);
	}

	parseCommand(argv);
}

func parseCommand(cmd []string) {
	switch cmd[1] {
	case "get":
		if len(cmd) < 3 {
			fmt.Fprintf(os.Stderr, "bk: `get` command needs the `name` argument\n")
			os.Exit(1);
		}

		url, err := getUrlByName(cmd[2]);
		if err != nil {
			panic(err);
		}

		fmt.Println(url);

	case "add":
		var name, url string;

		if len(cmd) < 3 {
			fmt.Fprintf(os.Stderr, "bk: `add` command needs the url argument\n");
			os.Exit(1);
		}

		url = cmd[2];

		if len(cmd) < 4 {
			name = url;
		} else {
			name = cmd[3];
		}

		m, err := parseBkFile();
		if err != nil {
			panic(err);
		}
		if m[name] != "" {
			fmt.Println("bookmark with this name already exists");
			os.Exit(1);
		}

		err = saveBookmark(name, url);
		if err != nil {
			panic(err);
		}
	case "list":
		m, err := parseBkFile();
		if err != nil {
			panic(err);
		}

		for b := range m {
			fmt.Println(b);
		}

	case "open":
	default:
		os.Exit(1);
	}
}

func saveBookmark(name string, url string) error {
	const path = "./bks";

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644);
	if err != nil {
		return err;
	}

	defer f.Close();

	line := fmt.Sprintf("%v %v\n", name, url);
	_, err = f.WriteString(line);	
	if err != nil {
		return err;
	}

	return nil;
}

func getUrlByName(name string) (string, error) {
	m, err := parseBkFile();
	if err != nil {
		return "", err;
	}

	return m[name], nil;
}

func parseBkFile() (map[string]string, error) {
	const path = "./bks";

	dat, err := os.ReadFile(path);
	if err != nil {
		return nil, fmt.Errorf("failed to read the file: %v", err);
	}

	lines := strings.Split(string(dat), "\n");

	m := make(map[string]string);

	for _,l := range lines {
		if strings.TrimSpace(l) == "" {
			continue;
		}

		parts := strings.SplitN(l, " ", 2);
		name := parts[0];
		url := parts[1];
		m[name] = url;
	}

	return m, nil;
}

