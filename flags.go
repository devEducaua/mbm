package main

import (
	"fmt"
	"os"
	"os/exec"
)

func parseFlags(argv []string) {

	// TODO: considerate turn it on a enum
	// default, group, file
	var mode = "default";
	var option = "";

	var command string;
	var name string;
	var url string;

	for i,arg := range argv {
		switch arg {

		case "-l", "--list":
			command = "list";
		case "-e", "--get":
			command = "get";
			name = argv[i+1];	
		case "-o", "--open":
			command = "open";
			name = argv[i+1];	

		case "-a", "--add":
			command = "add";
			url = argv[i+1];
		case "-n", "--name":
			name = argv[i+1];	

		case "-g", "--group":
			mode = "group";
			option = argv[i+1];
		case "-f", "--file":
			mode = "file";
			option = argv[i+1];
		}
	}

	var err error;

	switch command {
	case "list":
		err = listFlag(mode, option);
	case "get":
		err = getFlag(name, mode, option);
	case "open":
		err = openFlag(name, mode, option);
	case "add":
		if name == "" {
			name = url;
		}
		addFlag(url, name, mode, option);
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err);
		os.Exit(1);
	}
}

func listFlag(mode string, option string) error {
	var groups []Group;
	var err error;

	if mode == "file" {
		groups, err = parseConfig(option);
		if err != nil {
			return fmt.Errorf("failed to parse the file: %v", err);
		}
	}

	if mode == "default" || mode == "group" {
		groups, err = parseConfig("default");
		if err != nil {
			return fmt.Errorf("failed to parse the config: %v", err);
		}
	}

	var groupFounded = false;
	if mode == "group" {
		for _,group := range groups {
			if group.Name == option {
				for _,b := range group.Bookmarks {
					fmt.Println(b.Name);
					groupFounded = true;
				}
			}
		}

		if groupFounded {
			return fmt.Errorf("group not found: %v", option);
		}
		return nil;
	}

	for _,group := range groups {
		for _,b := range group.Bookmarks {
			fmt.Println(b.Name);
		}
	}
	
	return nil;
}

func addFlag(url string, name string, mode string, option string) error {
	bk := Bookmark{name, url};

	if mode == "default" {
		saveBookmark("default", bk);
	} else {
		saveBookmark(option, bk);
	}

	return nil;
}

func openFlag(arg string, mode string, option string) error {
	result, err := openGetFlag(arg, mode, option);
	if err != nil {
		return err;
	}

	if result == "" {
		return nil;
	}

	cmd := exec.Command("xdg-open", result);
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run external command: %v", err);
	}

	return nil;
}

func getFlag(arg string, mode string, option string) error {

	result, err := openGetFlag(arg, mode, option);
	if err != nil {
		return err;
	}

	fmt.Println(result);

	return nil;
}

func openGetFlag(arg string, mode string, option string) (string, error) {
	var groups []Group;
	var err error;
	if mode == "file" {
		groups, err = parseConfig(option);
	} else {
		groups, err = parseConfig("default");
	}

	if err != nil {
		return "", fmt.Errorf("failed to parse the config: %v", err);
	}
	
	var result string;

	for _,group := range groups {
		for _,bk := range group.Bookmarks {
			if bk.Name == arg {
				result = bk.Url;
			}
		}
	}

	return result, nil;
}
