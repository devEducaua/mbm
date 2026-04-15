
# mbm
**m**anage **b**ook**m**arks in a simple way.

## build
just run ```make```.

to install:
```sh
make install
```

## usage

usage information are all in the manpage.
to open, type:
```sh 
mbm --help
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

