package main

import (
	"os"
	"fmt"
	"strings"
	"mbm/internal/commands"
	"mbm/internal/flags"
)

func main() {
	argv := os.Args;	

	err := parseFlags(argv);

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err);
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

	f.Submit("query", "l", &query, "make a query");
	f.Submit("file", "f", &file, "specify file");
	f.Submit("tags", "t", &tags, "list tags");
	f.Submit("name", "n", &name, "specify a bookmakrs name");
	f.Submit("verbose", "v", &verbose, "enable verbose output");

	var command = argv[0];
	switch command {
	case "list":
		tagsArr := strings.Split(tags, ",");
		err = commands.ListFlag(tagsArr, verbose, query, file);
	case "add":
		if len(argv) < 2 {
			return fmt.Errorf("invalid operation: `--add` needs one argument.");
		}
		tagsArr := strings.Split(tags, ",");
		err = commands.AddFlag(argv[1], name, tagsArr, file);
	case "get":
		if len(argv) < 2 {
			return fmt.Errorf("invalid operation: `--get` needs one argument.");
		}
		err = commands.GetFlag(argv[1], verbose, file);
	case "open":
		if len(argv) < 2 {
			return fmt.Errorf("invalid operation: `--open` needs one argument.");
		}
		err = commands.OpenFlag(argv[1], file);
	case "edit":
		err = commands.EditFlag();
	case "copy":
		if len(argv) < 2 {
			return fmt.Errorf("invalid operation: `--copy` needs one argument.");
		}
		err = commands.CopyFlag(argv[1], verbose, file);
	case "import":
		err = commands.ImportFlag(file);
	case "help", "--help":
		err = commands.HelpFlag();
	}

	if err != nil {
		return err;
	}

	if err := f.Parse(argv); err != nil {
		return err;
	}
	return nil;
}

