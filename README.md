# A fork of Better tree (arguably)

Modified version of [bt](https://github.com/LeperGnome/bt):

- Solarized theme (terminal colors need to be set to Solarized!)
- Make preview optional
- Remove file info header
- Remove help

## Installation

```bash
go install git.sr.ht/~ekoeppen/bt/cmd/bt@latest
```

Or from source

```bash
make install
```

## Usage

```bash
bt [flags] [directory]

Flags:
  -i    In-place render (without alternate screen)
  -pad uint
        Edge padding for top and bottom (default 5)
  -p    Enable preview
```

Key bindings:

| key           | desc                                                   |
|---------------|--------------------------------------------------------|
| j / arr down  | Select next child                                      |
| k / arr up    | Select previous child                                  |
| h / arr left  | Move up a dir                                          |
| l / arr right | Enter selected directory                               |
| d             | Move selected child (then 'p' to paste)                |
| y             | Copy selected child (then 'p' to paste)                |
| D             | Delete selected child                                  |
| if / id       | Create file (if) / directory (id) in current directory |
| r             | Rename selected child                                  |
| e             | Edit selected file in $EDITOR                          |
| gg            | Go to top most child in current directory              |
| G             | Go to last child in current directory                  |
| enter         | Collapse / expand selected directory                   |
| esc           | Clear error message / stop current operation           |
| ?             | Toggle help                                            |
| q / ctrl+c    | Exit                                                   |

## Motivation

