package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func parseConfig(path string) ([]Group, error) {

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

	var groups []Group;
	
	lines := strings.SplitSeq(content, "\n");
	var currentGroup = Group{Name: "default"};
	for l := range lines {

		if strings.TrimSpace(l) == "" {
			continue;
		}

		if strings.HasPrefix(l, "@@ ") {
			groups = append(groups, currentGroup);
			currentGroup = Group{
				Name: l[3:],
			}
			continue;
		}

		parts := strings.SplitN(l, "=", 2);
		name := strings.TrimSpace(parts[0]);
		url := strings.TrimSpace(parts[1]);
		bk := Bookmark{name, url};

		currentGroup.Bookmarks = append(currentGroup.Bookmarks, bk);
	}
	groups = append(groups, currentGroup);

	return groups, nil
}

func saveBookmark(groupName string, bks ...Bookmark) error {
	configDir, err := getConfigDir();
	if err != nil {
		return err;
	}

	path := filepath.Join(configDir, "config");

	var result string;

	content, err := readFile(path);
	if err != nil {
		return err;
	}

	lines := strings.SplitSeq(content, "\n");
	var currentGroup = Group{Name: "default"};

	var finalLines []string;

	added := false;
	groupFounded := false;

	for l := range lines {
		finalLines = append(finalLines, l);

		if strings.HasPrefix(l, "@@ ") {
			currentGroup = Group{
				Name: l[3:],
			}
			continue;
		}

		if added == false {
			if currentGroup.Name == groupName {
				for _,bk := range bks {
					line := fmt.Sprintf("%v = %v", bk.Name, bk.Url);
					finalLines = append(finalLines, line);
				}
				added = true;
				groupFounded = true;
			}
		}
	}	

	if !groupFounded {
		groupLine := fmt.Sprintf("@@ %v", groupName);
		finalLines = append(finalLines, groupLine);
		for _,bk := range bks {
			line := fmt.Sprintf("%v = %v", bk.Name, bk.Url);
			finalLines = append(finalLines, line);
		}
	}

	result = strings.Join(finalLines, "\n");

	err = os.WriteFile(path, []byte(result), 0664);
	if err != nil {
		return fmt.Errorf("failed to write the file: %v", err);
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

