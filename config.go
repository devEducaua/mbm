package main

import (
	"os"
	"fmt"
	"strings"
)

const PATH = "./config";

type Group struct {
	Name string
	Bookmarks []Bookmark
}

type Bookmark struct {
	Name string
	Url string
}

func parseConfig() ([]Group, error) {
	dat, err := os.ReadFile(PATH);
	if err != nil {
		return nil, fmt.Errorf("failed to read the file: %v", err);
	}

	content := string(dat);

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

	return groups, nil
}


