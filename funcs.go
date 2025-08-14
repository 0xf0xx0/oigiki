package oigiki

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type formattag struct {
	Start, End int
}

var (
	/// matches each individual tag
	tagreg = regexp.MustCompile(`(\{[#\w\d/]+\})`)
)

// processes a tagged string into ansi
func ProcessTags(s string) string {
	return processString(s)
}

// add a format tag to a string
func TagString(s, tag string) string {
	if s == "" {
		return s
	}
	return "{" + tag + "}" + s
}

// strip format tags (NOT ansi) from line
func StripLine(line string) string {
	for _, tag := range slices.Backward(selectFormats(line)) {
		line = line[:tag.Start] + line[tag.End:]
	}
	return line
}

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
	tags := selectFormats(line)
	for i, tag := range slices.Backward(tags) {
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
					prevTag := tags[ii]
					prevColor := line[prevTag.Start+1 : prevTag.End-1]
					if color == prevColor {
						/// we found our matching color set!
						/// now move back one more (if possible) and use *that* color
						/// FIXME: actually find the previous color
						/// i do NOT want to be 3 loops deep
						if ii-1 < 0 {
							/// cant move back, just reset
							color = "fg"
							if strings.HasPrefix(color, "bg") {
								color = "bg"
							}
							break
						}
						resetTag := tags[ii-1]
						// print(color + " -> ")
						color = line[resetTag.Start+1 : resetTag.End-1]
						// print(color + "\n")
						break
					}
				}
			}
		}
		c := processTag("", color)
		/// TODO: find a better way to verify we have color?
		if len(c) > 0 {
			/// chop off the reset code
			/// because we're splitting by \x1b, the furst element is empty
			c = "\x1b" + strings.Split(c, "\x1b")[1]
		}
		line = line[:tag.Start] + c + line[tag.End:]
	}
	return line + "\x1b[0m"
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
