package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"mbm/internal/config"
	"mbm/internal/types"
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
func ListFlag(tags []string, verbose bool, query string, fp string) error {
    if fp == "" {
        fp = "default"
    }

    bks, err := config.ParseFile(fp);    
    if err != nil {
        return fmt.Errorf("failed to parse the file: %v", err);
    }

    var result []types.Bookmark;

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
func AddFlag(url string, name string, tags []string, fp string) error {
    if name == "" {
        name = url;
    }

	bk := types.Bookmark{Name: name, Url: url, Tags: tags};

	if fp != "" {
		err := config.SaveBookmark(fp, bk);
		if err != nil {
			return err;
		}
		return nil;
	}
	err := config.SaveBookmark("default", bk);
	if err != nil {
		return err;
	}

    return nil;
}

/* 
    supported flags: --open, --file
*/
func OpenFlag(name string, fp string) error {
    bk, err := getBookmarkByName(name, fp);
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
func GetFlag(name string, verbose bool, fp string) error {
    bk, err := getBookmarkByName(name, fp);
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

func getBookmarkByName(name string, fp string) (types.Bookmark, error) {
    var bks []types.Bookmark;
    var err error;

    if fp != "" {
        bks, err = config.ParseFile(fp);
    } else {
        bks, err = config.ParseFile("default");
    }

    if err != nil {
        return types.Bookmark{}, fmt.Errorf("failed to parse the config: %v", err);
    }
    
    for _,bk := range bks {
        if bk.Name == name {
            return bk, nil;
        }
    }
    return types.Bookmark{}, nil;
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

func HelpFlag() error {
	err := externalCommand("man", "mbm", "1");
	if err != nil {
		return err;
	}
	return nil;
}

func EditFlag() error {
	dir, err := config.GetConfigDir();
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

/* 
    supported flags: --import, --file
*/
func ImportFlag(fp string) error {
	if fp == "" {
		return fmt.Errorf("--import flag needs a file");
	}

	bks, err := config.ParseFile(fp);
	if err != nil {
		return err;
	}
	err = config.SaveBookmark("default", bks...);
	if err != nil {
		return err;
	}

	return nil;
}
/* 
    supported flags: --copy, --verbose, --file
*/
func CopyFlag(name string, verbose bool, fp string) error {
	bk, err := getBookmarkByName(name, fp);
	if err != nil {
		return err;
	}

	var cmd *exec.Cmd;
	sessionType := os.Getenv("XDG_SESSION_TYPE");
	switch sessionType {
	case "wayland":
		cmd = exec.Command("wl-copy", bk.Url);
	case "x11":
		cmd = exec.Command("xclip", "-selection", "clipboard", "-i");
		cmd.Stdin = strings.NewReader(bk.Url);
	default:
		return fmt.Errorf("invalid xdg_session_type: `%v`", sessionType);
	}

	if err := cmd.Run(); err != nil {
		return err;
	}

	if verbose {
		fmt.Println(bk.Url);
	}

	return nil;
}

