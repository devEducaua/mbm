package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	argv := os.Args;	

	if len(argv) < 2 {
		fmt.Fprintf(os.Stderr, "mbm: command not passed\n")
		os.Exit(1);
	}

	err := parseCommand(argv);
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err);
		os.Exit(1);
	}
}

func parseCommand(cmd []string) error {
	switch cmd[1] {
	case "get":
		if len(cmd) < 3 {
			return fmt.Errorf("`get` command needs the `name` argument\n")
		}

		name := cmd[2];

		m, err := parseConfigFile();
		if err != nil {
			return fmt.Errorf("failed to parse the config file: %v", err);
		}	

		url := m[name];
		if url == "" {
			return fmt.Errorf("bookmark with name: `%v` not found", name);
		}

		fmt.Println(url);

	case "add":
		var name, url string;

		if len(cmd) < 3 {
			return fmt.Errorf("`add` command needs the `url` argument\n")
		}

		url = cmd[2];

		if len(cmd) < 4 {
			name = url;
		} else {
			name = cmd[3];
		}

		m, err := parseConfigFile();
		if err != nil {
			return fmt.Errorf("failed to parse the config file: %v", err);
		}
		if m[name] != "" {
			return fmt.Errorf("a bookmark with the name: `%v` already exists", err);
		}

		err = saveBookmark(name, url);
		if err != nil {
			return fmt.Errorf("failed to save the bookmark: %v", err);
		}
	case "list":
		m, err := parseConfigFile();
		if err != nil {
			return fmt.Errorf("failed to parse the config file: %v", err);
		}

		for b := range m {
			fmt.Println(b);
		}

	case "open":
		if len(cmd) < 3 {
			return fmt.Errorf("`open` command needs the `name` argument\n")
		}

		name := cmd[2];

		m, err := parseConfigFile();
		if err != nil {
			return fmt.Errorf("failed to parse the config file: %v", err);
		}

		cmd := exec.Command("xdg-open", m[name]);
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run external command: %v", err);
		}
	default:
		return fmt.Errorf("command: `%v` not found", cmd[1]);
	}

	return nil;
}

func saveBookmark(name string, url string) error {
	const path = "./mbmc";

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

func parseConfigFile() (map[string]string, error) {
	const path = "./mbmc";

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
