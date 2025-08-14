# oigiki (WIP)

ansi color tagging lib, extracted from [axefetch](https://github.com/0xf0xx0/axefetch)  
uses https://github.com/fatih/color

## usage
```go
/// coloring will run till the next tag
/// chaining works too uwu
oigiki.ProcessTags("{red}red{green}green{/green}red") /// "\x1b[31mred\x1b[32mgreen\x1b[31mred\x1b[0m"
oigiki.TagString("this should be red", "red") /// "{red}this should be red"
oigiki.StripLine("{red}{bold}whos that {bg#f100f1}{underline}knockin{/underline} at my door?") /// whos that knockin at my door?
```

### supported colors
everything fatih/color supports, with rgb as hex

## known issues
- color untagging resets to the nearest style, no matter what
    see: TestRGB
## todo
- better processing?
- reset specific colors? (might need a pr to fatih/color)
- rgb tagging? `{r,g,b}`?
