package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type TokKind int;
const (
    TokNot TokKind = iota
    TokAnd
    TokOr
    TokTag
)

type Tok struct {
    Kind TokKind
	Value string
}

func tokenizeQuery(query string) []Tok {

    fields := strings.Fields(query);
    var toks []Tok;

    for _,v := range fields {
        var tok Tok;
        switch v {
        case "and":
            tok.Kind = TokAnd; 
        case "or":
            tok.Kind = TokOr;
        case "not":
            tok.Kind = TokNot;
        default:
            tok.Kind = TokTag;
            tok.Value = v;
        }
        toks = append(toks, tok);
    }

    return toks;
}

func popStack(stack *[]bool) (bool, error) {
	if len(*stack) <= 0 {
		return false, fmt.Errorf("could not pop empty stack");
	}

	last := len(*stack)-1;
	elem := (*stack)[last];
	*stack = append((*stack)[:last], (*stack)[last+1:]...);
	return elem, nil;
}

func applyQueryOnBk(tagBk map[string]bool, toks []Tok) bool {
	var stack []bool;

    for _,tok := range toks {
        switch tok.Kind {
			case TokTag:
				val := tagBk[tok.Value];
				stack = append(stack, val);

            case TokNot:
                a, err := popStack(&stack);
				if err != nil {
					return false;
				}
                stack = append(stack, !a);
            case TokAnd:
                a, err := popStack(&stack);
				if err != nil {
					return false;
				}
                b, err := popStack(&stack);
				if err != nil {
					return false;
				}
                stack = append(stack, a && b);
            case TokOr:
                a, err := popStack(&stack);
				if err != nil {
					return false;
				}
                b, err := popStack(&stack);
				if err != nil {
					return false;
				}
                stack = append(stack, a || b);
        }
    }

    return stack[0];
}

/* 
    supported flags: --list, --file, --tags, --verbose, --query
*/
func listFlag(tags []string, verbose bool, query string, fp string) error {
    if fp == "" {
        fp = "default"
    }

    bks, err := parseFile(fp);    
    if err != nil {
        return fmt.Errorf("failed to parse the file: %v", err);
    }

    var result []Bookmark;

    if query != "" {
        toks := tokenizeQuery(query);
        for _,bk := range bks {
            tagBk := make(map[string]bool);
            for _,t := range bk.Tags {
                tagBk[t] = true;
            }

			match := applyQueryOnBk(tagBk, toks);
			if match {
				result = append(result, bk);
			}
        }
    }

	if len(tags) != 0 {
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
	}

    if len(tags) == 0 && query == "" {
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
    if name == "" {
        name = url;
    }

    bk := Bookmark{name, url, tags};

	if fp != "" {
		err := saveBookmark(fp, bk);
		if err != nil {
			return err;
		}
		return nil;
	}
	err := saveBookmark("default", bk);
	if err != nil {
		return err;
	}

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
        fmt.Println(bk.Url);
    }

    return nil;
}

func openGetFlag(name string, fp string) (Bookmark, error) {
    var bks []Bookmark;
    var err error;

    if fp != "" {
        bks, err = parseFile(fp);
    } else {
        bks, err = parseFile("default");
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

func helpFlag() error {
	err := externalCommand("man", "mbm", "1");
	if err != nil {
		return err;
	}
	return nil;
}

func editFlag() error {
	dir, err := getConfigDir();
	if err != nil {
		return err;
	}

	editor := os.Getenv("EDITOR");
	err = externalCommand(editor, filepath.Join(dir, "config"));
	if err != nil {
		return err;
	}
	return nil;
}

func importFlag(fp string) error {
	if fp == "" {
		return fmt.Errorf("--import flag needs a file");
	}

	bks, err := parseFile(fp);
	if err != nil {
		return err;
	}
	err = saveBookmark("default", bks...);
	if err != nil {
		return err;
	}

	return nil;
}

func copyFlag(arg string, fp string) error {
	_ = arg;
	_ = fp;
	return nil;
}

