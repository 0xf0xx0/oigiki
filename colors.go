package oigiki

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type formattag struct {
	Start, End int
}
type formatmatch struct {
	tags      []formattag
	targetEnd int
}

var (
	/// matches groups of tags, eg `{bold}{green}foobar`
	taggroupreg = regexp.MustCompile(`(\{[#\w\d/]+\})+`)
	/// matches each individual tag
	tagreg = regexp.MustCompile(`(\{[#\w\d/]+\})`)
)
var colorMap = map[string]func(a ...interface{}) string{
	"bold":      color.New(color.Bold).SprintFunc(),
	"underline": color.New(color.Underline).SprintFunc(),
	"italic":    color.New(color.Italic).SprintFunc(),
	"reset":     color.New(color.Reset).SprintFunc(),
	"/":         color.New(color.Reset).SprintFunc(),

	"red":     color.New(color.FgRed).SprintFunc(),
	"green":   color.New(color.FgGreen).SprintFunc(),
	"blue":    color.New(color.FgBlue).SprintFunc(),
	"yellow":  color.New(color.FgYellow).SprintFunc(),
	"cyan":    color.New(color.FgCyan).SprintFunc(),
	"magenta": color.New(color.FgMagenta).SprintFunc(),
	"white":   color.New(color.FgWhite).SprintFunc(),
	"black":   color.New(color.FgBlack).SprintFunc(),

	"redbright":     color.New(color.FgHiRed).SprintFunc(),
	"greenbright":   color.New(color.FgHiGreen).SprintFunc(),
	"bluebright":    color.New(color.FgHiBlue).SprintFunc(),
	"yellowbright":  color.New(color.FgHiYellow).SprintFunc(),
	"cyanbright":    color.New(color.FgHiCyan).SprintFunc(),
	"magentabright": color.New(color.FgHiMagenta).SprintFunc(),
	"whitebright":   color.New(color.FgHiWhite).SprintFunc(),
	"blackbright":   color.New(color.FgHiBlack).SprintFunc(),

	"bgred":     color.New(color.BgRed).SprintFunc(),
	"bggreen":   color.New(color.BgGreen).SprintFunc(),
	"bgblue":    color.New(color.BgBlue).SprintFunc(),
	"bgyellow":  color.New(color.BgYellow).SprintFunc(),
	"bgcyan":    color.New(color.BgCyan).SprintFunc(),
	"bgmagenta": color.New(color.BgMagenta).SprintFunc(),
	"bgwhite":   color.New(color.BgWhite).SprintFunc(),
	"bgblack":   color.New(color.BgBlack).SprintFunc(),

	"bgredbright":     color.New(color.BgHiRed).SprintFunc(),
	"bggreenbright":   color.New(color.BgHiGreen).SprintFunc(),
	"bgbluebright":    color.New(color.BgHiBlue).SprintFunc(),
	"bgyellowbright":  color.New(color.BgHiYellow).SprintFunc(),
	"bgcyanbright":    color.New(color.BgHiCyan).SprintFunc(),
	"bgmagentabright": color.New(color.BgHiMagenta).SprintFunc(),
	"bgwhitebright":   color.New(color.BgHiWhite).SprintFunc(),
	"bgblackbright":   color.New(color.BgHiBlack).SprintFunc(),
}

// processes a tagged string, passing it through fmt.Sprintf before coloring
func ProcessTags(s string, a ...any) string {
	return processString(fmt.Sprintf(s, a...))
}

// add a format tag to a string
func TagString(s, color string) string {
	if s == "" {
		return s
	}
	return "{"+color+"}"+s
}

// strip format tags (NOT ansi) from line
func StripLine(line string) string {
	ret := ""
	for _, format := range selectFormats(line) {
		lastTag := format.tags[len(format.tags)-1]
		ret += line[lastTag.End:format.targetEnd]
	}
	return ret
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
// pass 2
// process tags in reverse
func processString(line string) string {
	formats := selectFormats(line)
	/// TODO: slices.Backward
	slices.Reverse(formats)
	/// find+replace in reverse to avoid indexes jumping around
	for _, fmt := range formats {
		tags := fmt.tags
		/// save tags before flippin the array
		furstTag := tags[0]
		lastTag := tags[len(tags)-1]
		/// TODO: slices.Backward
		slices.Reverse(tags)
		/// save the target line and update it before substring-replacing
		formatted := line[lastTag.End:fmt.targetEnd]
		for _, tag := range tags {
			color := line[tag.Start+1 : tag.End-1] /// {(style)}
			formatted = processTag(formatted, color)
		}
		line = line[:furstTag.Start] + formatted + line[fmt.targetEnd:]
	}
	return line
}

// pass 1
// finds format tags by regex, matching them and their target end idx
func selectFormats(line string) []formatmatch {
	indices := taggroupreg.FindAllStringIndex(line, -1)
	indicesLen := len(indices)
	ret := make([]formatmatch, indicesLen)
	for i, element := range indices {
		endIdx := 0
		if i+1 < indicesLen {
			endIdx = indices[i+1][0]
		} else {
			endIdx = len(line)
		}
		/// this is either '{tag}' or '{tag}{tag}...'
		tagGroup := line[element[0]:element[1]]
		tagIndices := tagreg.FindAllStringIndex(tagGroup, -1)
		tags := make([]formattag, len(tagIndices))
		for ii, tag := range tagIndices {
			/// tag[] is relative to the outer match
			/// make it relative to the whole line
			tags[ii] = formattag{Start: element[0] + tag[0], End: element[0] + tag[1]}
		}
		ret[i] = formatmatch{tags: tags, targetEnd: endIdx}
	}
	return ret
}
