package oigiki

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// processes a tagged string, passing it through fmt.Sprintf before coloring
func ProcessTags(s string, a ...any) string {
	return processString(fmt.Sprintf(s, a...))
}

// add a format tag to a string
func TagString(s, color string) string {
	if s == "" {
		return s
	}
	return "{" + color + "}" + s
}

// strip format tags (NOT ansi) from line
// func StripLine(line string) string {
// 	ret := ""
// 	for _, format := range selectFormats(line) {
// 		lastTag := format.tags[len(format.tags)-1]
// 		ret += line[lastTag.End:format.targetEnd]
// 	}
// 	return ret
// }

func hex2RGB(hex string) (int, int, int) {
	hex = hex[1:]
	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	b, _ := strconv.ParseUint(hex[4:6], 16, 8)
	return int(r), int(g), int(b)
}

func processTag(s, col string) string {
	if strings.HasPrefix(col, "#") {
		/// hex code
		return color.RGB(hex2RGB(col)).Sprint(s)
	}
	if strings.HasPrefix(col, "bg#") {
		return color.BgRGB(hex2RGB(col[2:])).Sprint(s)
	}
	if fn, ok := colorMap[col]; ok {
		return fn(s)
	}
	return s
}

// pass 2:
// process tags in reverse
func processString(line string) string {
	formats := selectFormats(line)
	/// save the target line and update it before substring-replacing
	for i, tag := range slices.Backward(formats) {
		color := line[tag.Start+1 : tag.End-1] /// {(style)}
		if strings.HasPrefix(color, "/") {
			/// if its in the map its one of the ones ansi already has resets for
			/// we only process colors
			if _, ok := colorMap[color]; !ok {
				/// look ahead (behind?) and find the previous color in the string
				/// to reset to
				/// slice off the leading slash too for exact matching
				color = color[1:]
				for ii := i - 1; ii >= 0; ii-- {
					prevTag := formats[ii]
					prevColor := line[prevTag.Start+1 : prevTag.End-1]
					if color == prevColor {
						/// we found our matching color set!
						/// now move back one more (if possible) and use *that* color
						if ii-1 < 0 {
							/// cant move back
							color = ""
							break
						}
						resetTag := formats[ii-1]
						color = line[resetTag.Start+1 : resetTag.End-1]
						break
					}
				}
			}
		}
		formatted := processTag("", color)
		/// TODO: find a better way to verify we have color?
		if len(formatted) > 0 {
			/// chop off the reset code
			/// because we're splitting by \x1b, the furst element is empty
			formatted = "\x1b" + strings.Split(formatted, "\x1b")[1]
		}
		line = line[:tag.Start] + formatted + line[tag.End:]
	}
	return line
}

// pass 1:
// finds format tags by regex, returning their position in the string
func selectFormats(line string) []formattag {
	tagIndices := tagreg.FindAllStringIndex(line, -1)
	tags := make([]formattag, len(tagIndices))
	for i, tag := range tagIndices {
		tags[i] = formattag{Start: tag[0], End: tag[1]}
	}
	return tags
}
