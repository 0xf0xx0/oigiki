# oigiki

fast ansi color tagging lib, similar to [blessed](https://github.com/chjj/blessed) from node  
uses https://github.com/fatih/color  
extracted from [axefetch](https://github.com/0xf0xx0/axefetch)

## usage
```go
/// coloring will run until overridden
/// chaining works too uwu
oigiki.ProcessTags("{red}red{green}green{/green}red") /// "\x1b[31mred\x1b[32mgreen\x1b[31mred\x1b[0m"
/// unknown tags are ignored
oigiki.ProcessTags("{red}red {totallyrealcolor}still red") /// "\x1b[31mred still red\x1b[0m"

oigiki.TagString("this should be red", "red") /// "{red}this should be red"
oigiki.StripLine("{red}{bold}whos that {bg#f100f1}{underline}knockin{/underline} at my door?") /// whos that knockin at my door?
```

## supported colors
everything fatih/color supports, with rgb as hex
```go
/// hex must be the full rrggb
"{#ff00ff}ff00ff"
/// bg is easy too
"{bg#ff00ff}{black}ff00ff"
```

## todo
- rgb tagging? `{r,g,b}`?
