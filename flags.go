package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func parseFlags(argv []string) error {
	var fp = "";

	var command string;
	var arg string;
	var secArg string;
	var tags []string;
	var verbose = false;

	var parsingError error = nil;

	for i,v := range argv {
		switch v {

		case "-l", "--list":
			command = "list";
		case "-g", "--get":
			command = "get";
			if len(argv) <= i+1 {
				parsingError = fmt.Errorf("invalid operation: `--get` needs one argument.");
				break;
			}
			arg = argv[i+1];	
		case "-o", "--open":
			command = "open";
			if len(argv) <= i+1 {
				parsingError = fmt.Errorf("invalid operation: `--open` needs one argument.");
				break;
			}
			arg = argv[i+1];	
		case "-a", "--add":
			command = "add";
			if len(argv) <= i+1 {
				parsingError = fmt.Errorf("invalid operation: `--add` needs one argument, if option --file is not provided.");
				break;
			}

			if argv[i+1] != "--file" {
				arg = argv[i+1];
			}
		case "--help":
			command = "help";	
		case "--edit":
			command = "edit";
		case "-n", "--name":
			if len(argv) <= i+1 {
				parsingError = fmt.Errorf("invalid operation: `--name` needs one argument.");
				break;
			}
			secArg = argv[i+1];	
		case "-t", "--tags":
			if len(argv) <= i+1 {
				parsingError = fmt.Errorf("invalid operation: `--tags` needs one argument.");
				break;
			}
			tags = strings.Split(argv[i+1], ",");

		case "-v", "--verbose":
			verbose = true;

		case "-f", "--file":
			if len(argv) <= i+1 {
				parsingError = fmt.Errorf("invalid operation: `--file` needs one argument.");
				break;
			}
			fp = argv[i+1];
		}
	}
	if parsingError != nil {
		return parsingError;	
	}

	var err error;

	switch command {
	case "help":
		err = externalCommand("man", "mbm", "1");
		if err != nil {
			return err;
		}
	case "edit":
		dir, err := getConfigDir();
		if err != nil {
			return err;
		}

		editor := os.Getenv("EDITOR");
		err = externalCommand(editor, filepath.Join(dir, "config"));
		if err != nil {
			return err;
		}

	case "list":
		err = listFlag(tags, verbose, fp);
	case "get":
		err = getFlag(arg, verbose, fp);
	case "open":
		err = openFlag(arg, fp);
	case "add":
		err = addFlag(arg, secArg, tags, fp);
	}

	if err != nil {
		return err;
	}
	return nil;
}

/* 
	supported flags: --list, --file, --tags, --verbose
*/
func listFlag(tags []string, verbose bool, fp string) error {
	if fp == "" {
		fp = "default"
	}

	bks, err := parseConfig(fp);	
	if err != nil {
		return fmt.Errorf("failed to parse the file: %v", err);
	}

	var result []Bookmark;

	for _,bk := range bks {
		tagBk := make(map[string]bool);
		for _,t := range bk.Tags {
			tagBk[t] = true;
		}

		match := true;
		for _,t := range tags {
			if !tagBk[t] {
				match = false;
			}
		}

		if match {
			result = append(result, bk);
		}
	}

	if len(tags) == 0 {
		result = bks;
	}

	for _,bk := range result {
		if verbose {
			fmt.Printf("%v %v %v\n", bk.Name, bk.Url, strings.Join(bk.Tags, ","));
		} else {
			fmt.Println(bk.Name);
		}
	}

	return nil;
}

/* 
	supported flags: --add, --file, --tags, --name
*/
func addFlag(url string, name string, tags []string, fp string) error {
	if fp != "" {
		bks, err := parseConfig(fp);
		if err != nil {
			return err;
		}
		saveBookmark(bks...);
		return nil;
	}

	if name == "" {
		name = url;
	}

	bk := Bookmark{name, url, tags};
	saveBookmark(bk);

	return nil;
}

/* 
	supported flags: --open, --file
*/
func openFlag(name string, fp string) error {
	bk, err := openGetFlag(name, fp);
	if err != nil {
		return err;
	}

	if bk.Url == "" {
		return nil;
	}

	cmd := exec.Command("xdg-open", bk.Url);
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run external command: %v", err);
	}

	return nil;
}

/* 
	supported flags: --get, --verbose, --file
*/
func getFlag(name string, verbose bool, fp string) error {
	bk, err := openGetFlag(name, fp);
	if err != nil {
		return err;
	}

	if verbose {
		fmt.Printf("%v %v %v\n", bk.Name, bk.Url, strings.Join(bk.Tags, ","));
	} else {
		fmt.Println(bk.Name);
	}

	return nil;
}

func openGetFlag(name string, fp string) (Bookmark, error) {
	var bks []Bookmark;
	var err error;

	if fp != "" {
		bks, err = parseConfig(fp);
	} else {
		bks, err = parseConfig("default");
	}

	if err != nil {
		return Bookmark{}, fmt.Errorf("failed to parse the config: %v", err);
	}
	
	for _,bk := range bks {
		if bk.Name == name {
			return bk, nil;
		}
	}
	return Bookmark{}, nil;
}

func externalCommand(cmds ...string) error {
	cmd := exec.Command(cmds[0], cmds[1:]...);
	cmd.Stdin = os.Stdin;
	cmd.Stdout = os.Stdout;
	cmd.Stderr = os.Stderr;
	err := cmd.Run();
	if err != nil {
		return err;
	}
	return nil;
}
