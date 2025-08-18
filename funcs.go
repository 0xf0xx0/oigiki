package oigiki

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	/// matches each individual tag
	tagreg = regexp.MustCompile(`(\{[#\w\d/]+\})`)
)

// processes a tagged string into ansi
func ProcessTags(s string) string {
	if color.NoColor {
		return StripLine(s)
	}
	/// we can safely assume we have color
	out := strings.Builder{} /// TODO: can we avoid this?
	out.Grow(len(s))
	tag := strings.Builder{}
	collectingTag := false
	seenTags := make([]string, 1, 8) /// 8 should be plenty for most cases
	seenTags[0] = "fg"               /// see TestDefaultUnset
	for _, char := range s {
		if char == '{' {
			if collectingTag {
				panic("tag inside a tag")
			}
			collectingTag = true
			continue
		}
		if char == '}' && collectingTag {
			tagStr := tag.String()
			color := processTag("", tagStr)

			/// processTag checks the color map for the tag
			/// if its in the map its one of the ones ansi already has resets for,
			/// and color will be populated
			/// we only process closing color tags
			if color == "" && tagStr[0] == '/' {
				/// find the corresponding opening tag
				idx := slices.Index(seenTags, tagStr[1:])
				if idx > -1 {
					/// remove the opening tag from our seen list so it doesnt get picked again
					seenTags = append(seenTags[:idx], seenTags[idx+1:]...)
					/// if its not the last seen tag, select the last seen tag to reset to
					/// see TestUnset and TestDefaultUnset
					l := len(seenTags) - 1
					if idx != l {
						color = processTag("", seenTags[l])
					}
					/// no-op, if idx == l its already the active tag and doesnt need re-applying
					/// see TestUnset2
				}
				/// no-op
			}
			/// also no-op, for unknown tags

			if color != "" {
				out.WriteString("\x1b")
				out.WriteString(strings.Split(color, "\x1b")[1])
			}
			/// i dont like thissssssssssssssssssssssss
			/// TODO: make a way to just ignore non-color tags
			if tagStr[0] != '/' && tagStr != "bold" && tagStr != "underline" && tagStr != "italic" {
				seenTags = append(seenTags, tagStr)
			}
			tag.Reset()
			collectingTag = false
			continue
		}
		if collectingTag {
			tag.WriteRune(char)
		} else {
			out.WriteRune(char)
		}
	}
	/// reset
	out.WriteString("\x1b[0m")
	return out.String()
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
	for _, tag := range slices.Backward(tagreg.FindAllStringIndex(line, -1)) {
		line = line[:tag[0]] + line[tag[1]:]
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
	if col[0] == '#' {
		/// hex code
		return color.RGB(hex2RGB(col)).Sprint(s)
	}
	if col[0] == 'b' && col[1] == 'g' && col[2] == '#' {
		return color.BgRGB(hex2RGB(col[2:])).Sprint(s)
	}
	if fn, ok := colorMap[col]; ok {
		return fn(s)
	}
	return s
}
