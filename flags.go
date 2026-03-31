package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// TODO: --query not implemented
func parseFlags(argv []string) {
	var filepath = "";

	var command string;
	var arg string;
	var secArg string;
	var tags []string;
	var verbose = false;

	for i,arg := range argv {
		switch arg {

		case "-l", "--list":
			command = "list";
		case "-g", "--get":
			command = "get";
			arg = argv[i+1];	
		case "-o", "--open":
			command = "open";
			arg = argv[i+1];	
		case "-a", "--add":
			command = "add";
			arg = argv[i+1];
		case "-n", "--name":
			secArg = argv[i+1];	
		case "-t", "--tags":
			tags = strings.Split(argv[i+1], ",");
		case "-q", "--query":
			secArg = argv[i+1];
			
		case "-f", "--file":
			filepath = argv[i+1];
		}
	}

	var err error;

	switch command {
	case "list":
		err = listFlag(secArg, verbose, filepath);
	case "get":
		err = getFlag(arg, verbose, filepath);
	case "open":
		err = openFlag(arg, filepath);
	case "add":
		err = addFlag(arg, secArg, tags, filepath);
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err);
		os.Exit(1);
	}
}

/* 
	supported flags: --list, --file, --query, --verbose
*/
func listFlag(query string, verbose bool, filepath string) error {
	var bks []Bookmark;
	var err error;

	if filepath != "" {
		bks, err = parseConfig(filepath);	
	} else {
		bks, err = parseConfig("default");	
	}

	if query != "" {
		return fmt.Errorf("query is not implemented yet");
	}

	if err != nil {
		return fmt.Errorf("failed to parse the file: %v", err);
	}

	for _,bk := range bks {
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
func addFlag(url string, name string, tags []string, filepath string) error {
	if filepath != "" {
		bks, err := parseConfig(filepath);
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
func openFlag(name string, filepath string) error {
	bk, err := openGetFlag(name, filepath);
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
func getFlag(name string, verbose bool, filepath string) error {
	bk, err := openGetFlag(name, filepath);
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

func openGetFlag(name string, filepath string) (Bookmark, error) {
	var bks []Bookmark;
	var err error;

	if filepath != "" {
		bks, err = parseConfig(filepath);
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
