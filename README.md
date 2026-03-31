
# mbm
**m**anage **b**ook**m**arks in a simple way.

## build
just run ```make```.

to install:
```sh
make install
```

## usage

mbm has two groups of flags:
- command flags: --add, --list, --get, --open
- mode flags: --file, --group

each command flag can be combined with a mode flag, meaning that each command has 3 variations.

examples using `--list`:

```sh
mbm --list # will list all the bookmarks in ~/.config/mbm/config.

mbm --list --group links # will list all the bookmarks in the `links` group.

mbm --list --file ./bks.txt # will list all the bookmarks in the file ./bks.txt.
```

## usage in scripts

example usage of mbm with dmenu and fzf

for dmenu:
```sh
mbm open $(mbm --list | dmenu -i)
```

for fzf:
```sh
mbm open $(mbm --list | fzf)
```

