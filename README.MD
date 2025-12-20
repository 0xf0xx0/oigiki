# oigiki
[![Go Reference](https://pkg.go.dev/badge/github.com/0xf0xx0/oigiki.svg)](https://pkg.go.dev/github.com/0xf0xx0/oigiki)

fast ansi color tagging lib, inspired by [blessed](https://github.com/chjj/blessed) from node   
extracted from [axefetch](https://github.com/0xf0xx0/axefetch) and improved

## features
- [`NO_COLOR`](https://no-color.org/) support
- all 16 colors (normal + bright)
- truecolor (as rrrggbb hex)

## usage
```go
/// coloring will run until overridden
/// chaining works too uwu
oigiki.ProcessTags("{red}red{green}green{/green}red") /// "\x1b[31mred\x1b[32mgreen\x1b[31mred\x1b[0m"
/// unknown tags are ignored
oigiki.ProcessTags("{red}red {totallyrealcolor}still red") /// "\x1b[31mred {totallyrealcolor}still red\x1b[0m"

oigiki.TagString("this should be red", "red") /// "{red}this should be red"
oigiki.StripTags("{red}{bold}whos that {bg#f100f1}{underline}knockin{/underline} at my door?") /// whos that knockin at my door?
```

### truecolor
```go
/// hex must be the full rrggbb
"{#ff00ff}ff00ff"
/// bg is easy too
"{bg#ff00ff}{black}ff00ff"
```
