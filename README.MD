# oigiki

[![Go Reference](https://pkg.go.dev/badge/git.0xf0xx0.eth.limo/0xf0xx0/oigiki.svg)](https://pkg.go.dev/git.0xf0xx0.eth.limo/0xf0xx0/oigiki)

fast ansi color tagging lib, inspired by [blessed](https://github.com/chjj/blessed) from node  
extracted from [axefetch](https://github.com/0xf0xx0/axefetch) and improved  
minimum go version: 1.10

## features

- [`NO_COLOR`](https://no-color.org/) support
- [`FORCE_COLOR`](https://force-color.org/) support (takes precedence over `NO_COLOR`)
- [`CLICOLOR_FORCE`](https://bixense.com/clicolors/) support (takes precedence over `NO_COLOR`)
- all 16 colors (normal + bright)
- 256color (as index, 0-255)
- truecolor (as rrrggbb hex)

## import

```sh
go get git.0xf0xx0.eth.limo/0xf0xx0/oigiki
```

## usage

```go
/// coloring will run until overridden
/// chaining works too uwu
oigiki.ProcessTags("{red}red{green}green{/green}red") /// -> "\x1b[31mred\x1b[32mgreen\x1b[31mred\x1b[0m"
/// unknown tags are ignored
oigiki.ProcessTags("{red}red {totallyrealcolor}still red") /// -> "\x1b[31mred {totallyrealcolor}still red\x1b[0m"

oigiki.TagString("this should be red", "red") /// -> "{red}this should be red"
oigiki.StripTags("{red}{bold}whos that {bg#f100f1}{underline}knockin{/underline} at my door?") /// -> whos that knockin at my door?
```

### 256color

```go
"{bg0}{7}yin{bg7}{0}yang" /// -> "\x1b[48;5;0m\x1b[38;5;7myin\x1b[48;5;7m\x1b[38;5;0myang\x1b[0m"
```

### truecolor

```go
/// hex must be the full rrggbb
"{bg#d60270}{#0038a8}po{/}{#9b4f96}go{bg#0038a8}{#d60270}lo{/}!" /// -> "\x1b[48;2;214;2;112m\x1b[38;2;0;56;168mpo\x1b[0m\x1b[38;2;155;79;150mgo\x1b[48;2;0;56;168m\x1b[38;2;214;2;112mlo\x1b[0m!\x1b[0m"
```
