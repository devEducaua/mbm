package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func parseConfig(path string) ([]Bookmark, error) {

	configDir, err := getConfigDir();
	if err != nil {
		return nil, err;
	}

	defaultPath := filepath.Join(configDir, "config");

	if path == "default" {
		path = defaultPath;
	}
	
	content, err := readFile(path);
	if err != nil {
		return nil, err;
	}

	var bks []Bookmark;

	lines := strings.Split(content, "\n");
	for i,l := range lines {
		if strings.TrimSpace(l) == "" {
			continue;
		}

		parts := strings.SplitN(l, " ", 3);
		var tags = []string{};

		if len(parts) < 2 {
			fmt.Printf("failed to parse the bookmark line: %v\n", i);
			continue;
		}

		name := strings.TrimSpace(parts[0]);
		url := strings.TrimSpace(parts[1]);

		if len(parts) == 3 {
			tagsStr := strings.TrimSpace(parts[2]);
			tags = strings.Split(tagsStr, ",");
		}

		bk := Bookmark{name, url, tags};
		bks = append(bks, bk);
	}

	return bks, nil
}

func saveBookmark(bks ...Bookmark) error {
	configDir, err := getConfigDir();
	if err != nil {
		return err;
	}

	path := filepath.Join(configDir, "config");

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644);
	if err != nil {
		return err;
	}
	defer f.Close();

	for _,bk := range bks {
		line := fmt.Sprintf("%v %v %v\n", bk.Name, bk.Url, strings.Join(bk.Tags, ","));

		_, err = f.WriteString(line);
		if err != nil {
			return err;
		}
	}

	return nil;
}

func getConfigDir() (string, error) {
	home, err := os.UserHomeDir();
	if err != nil {
		return "", err;
	}

	path := filepath.Join(home, ".config", "mbm");

	err = os.MkdirAll(path, 0755);
	if err != nil {
		return "", err;
	}

	return path, nil;	
}

