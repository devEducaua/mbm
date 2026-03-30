package main

import (
	"os"
	"fmt"
	"strings"
)

func parseConfig(path string) ([]Group, error) {
	defaultPath := "./config"
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

func saveBookmark(bk Bookmark, groupName string) error {
	path := "./config";

	var result string;

	content, err := readFile(path);
	if err != nil {
		return err;
	}

	var groups []Group;
	
	lines := strings.SplitSeq(content, "\n");
	var currentGroup = Group{Name: "default"};

	var finalLines []string;

	added := false;
	for l := range lines {
		finalLines = append(finalLines, fmt.Sprintf("%v\n", l));

		if strings.HasPrefix(l, "@@ ") {
			groups = append(groups, currentGroup);
			currentGroup = Group{
				Name: l[3:],
			}
			continue;
		}

		if added == false {
			if currentGroup.Name == groupName {
				line := fmt.Sprintf("%v = %v\n", bk.Name, bk.Url);
				finalLines = append(finalLines, line);
				added = true;
			}
		}
	}	

	result = strings.Join(finalLines, "");

	err = os.WriteFile(path, []byte(result), 0664);
	if err != nil {
		return fmt.Errorf("failed to write the file: %v", err);
	}

	return nil;
}

