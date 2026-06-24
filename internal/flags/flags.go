
package flags

import (
	"fmt"
	"strconv"
	"strings"
)

type Flag struct {
	Value any
	Usage string
}

type Flags map[string]Flag;

func (flags Flags) Submit(long, short string, value any, usage string) {
	f := &Flag{
		Value: value, 
		Usage: usage,
	};
	flags[long] = *f;
	flags[short] = *f;
}

func (flags Flags) Parse(argv []string) error {
	for i := range argv {
		arg := argv[i];

		var flag Flag;
		ok := false;

		switch {
		case strings.HasPrefix(arg, "--"):
			flag, ok = flags[arg[2:]];
			if !ok {
				continue;
			}
			if err := parseFlag(flag, i, argv); err != nil {
				return err;
			}
		case strings.HasPrefix(arg, "-"):
			withoutPrefix := arg[1:];
			flag, ok = flags[withoutPrefix];
			if !ok {
				continue;
			}
			if err := parseFlag(flag, i, argv); err != nil {
				return err;
			}
		}
	}
	return nil;
}

func parseFlag(flag Flag, i int, argv []string) error {
	arg := argv[i];
	switch v := flag.Value.(type) {
	case *bool:
		*v = true;
	case *string:
		if len(argv) <= i+1 {
			return fmt.Errorf("%v flag need an argument", arg);
		}
		*v = argv[i+1];
	case *int:
		if len(argv) <= i+1 {
			return fmt.Errorf("%v flag need an argument", arg);
		}
		value, err := strconv.Atoi(argv[i+1]);
		if err != nil {
			return fmt.Errorf("%v needs to be an integer", arg);
		}
		*v = value;
	}
	return nil;
}

