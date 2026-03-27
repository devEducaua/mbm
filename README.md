
# mbm
**m**anage **b**ook**m**arks in a simple way.

# commands

- list -> list all bookmarks name on stdout
- open "name" -> open bookmark with xdg-open
- get "name" -> returns an url based on the name
- add "url" [name] -> save a new bookmark

# usage in scripts

example usage of mbm with dmenu and fzf

for dmenu:
```bash
mbm open $(mbm list | dmenu -i)
```

for fzf:
```bash
mbm open $(mbm list | fzf)
```

