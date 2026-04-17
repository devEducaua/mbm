package main

import (
    "fmt"
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

        case "-q", "--query":
            if len(argv) <= i+1 {
                parsingError = fmt.Errorf("invalid operation: `--query` needs one argument.");
                break;
            }
            arg = argv[i+1];

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
        err = helpFlag();
    case "edit":
        err = editFlag();
    case "list":
        err = listFlag(tags, verbose, arg, fp);
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

