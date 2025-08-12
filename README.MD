# oigiki (WIP)

ansi color tagging lib, extracted from [axefetch](https://github.com/0xf0xx0/axefetch)  
uses https://github.com/fatih/color

## usage
```go
/// coloring will run till the next tag
/// chaining works too uwu
oigiki.ProcessTags("{bgwhite}{black}hello{reset} {#b00b69}world!") /// "\x1b[31mhello \x1b[0m\x1b[38;2;176;11;105mworld!\x1b[0;22;0;0;0m"
oigiki.TagString("this should be red", "red") /// "{red}this should be red"
oigiki.StripLine("{red}{bold}whos that {bg#f100f1}{underline}knockin{/underline} at my door?") /// whos that knockin at my door?
```

### supported colors
everything fatih/color supports, with rgb as hex

## todo
- better processing?
- reset specific colors? (might need a pr to fatih/color)
- rgb tagging? `{r,g,b}`?
