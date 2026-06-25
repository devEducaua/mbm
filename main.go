package main

import (
	"os"
	"fmt"
	"mbm/internal/flags"
	"mbm/internal/commands"
)

func main() {
	argv := os.Args[1:];
	err := parseFlags(argv);

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err);
		os.Exit(1);
	}
}

func parseFlags(argv []string) error {
	if len(argv) == 0 {
		return fmt.Errorf("no command specified");
	}

	f := make(flags.Flags);

	var (
		query string
		file string
		tags string
		name string
		verbose bool

		err error
	)

	f.Var("query", "q", &query);
	f.Var("file", "f", &file);
	f.Var("tags", "t", &tags);
	f.Var("name", "n", &name);
	f.Var("verbose", "v", &verbose);

	if err := f.Parse(argv[1:]); err != nil {
		return err;
	}

	var command = argv[0];
	switch command {
	case "list":
		err = commands.ListFlag(tags, verbose, query, file);
	case "add":
		if len(argv) < 2 {
			return fmt.Errorf("invalid operation: `add` needs one argument.");
		}
		err = commands.AddFlag(argv[1], name, tags, file);
	case "get":
		if len(argv) < 2 {
			return fmt.Errorf("invalid operation: `get` needs one argument.");
		}
		err = commands.GetFlag(argv[1], verbose, file);
	case "open":
		if len(argv) < 2 {
			return fmt.Errorf("invalid operation: `open` needs one argument.");
		}
		err = commands.OpenFlag(argv[1], file);
	case "edit":
		err = commands.EditFlag();
	case "copy":
		if len(argv) < 2 {
			return fmt.Errorf("invalid operation: `copy` needs one argument.");
		}
		err = commands.CopyFlag(argv[1], verbose, file);
	case "import":
		err = commands.ImportFlag(file);
	case "help", "--help":
		err = commands.HelpFlag();
	default:
		return fmt.Errorf("unknown command: `%v`", command);
	}

	if err != nil {
		return err;
	}

	return nil;
}

