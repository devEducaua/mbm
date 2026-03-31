
# mbm
**m**anage **b**ook**m**arks in a simple way.

## build
just run ```make```.

to install:
```sh
make install
```

## usage

```sh
mbm --list # will list all the bookmarks in ~/.config/mbm/config.

mbm --list --file ./bks.txt # will list all the bookmarks in the file ./bks.txt.

mbm --add https://example.com --name example --tags test,example # will save a new bookmark.

mbm open example # will open the bookmark with `xdg-open`.j
```

## usage in scripts

example usage of mbm with dmenu and fzf

for dmenu:
```sh
mbm --open $(mbm --list | dmenu)
```

for fzf:
```sh
mbm --open $(mbm --list | fzf)
```

