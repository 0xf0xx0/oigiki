// fast ansi color tagging lib, inspired by [chjj/blessed] from node
// and using [github.com/fatih/color]
//
// [chjj/blessed]: https://github.com/chjj/blessed
package oigiki

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

const tAG_DELIM_OPEN = '{'
const tAG_DELIM_CLOSE = '}'
const tAG_CLOSE_MARKER = '/'

const eSCAPE_CODE_RESET = "\x1b[0m"

type decorationFlags struct {
	underline bool
	italic    bool
	bold      bool
}

type TagType uint8

const (
	TagTypeUnknown TagType = iota
	TagTypeReset
	TagTypeColor
	TagTypeBold
	TagTypeItalic
	TagTypeUnderline
)

var (
	/// matches each individual tag
	tagreg  = regexp.MustCompile(`(\{.+?\})`)
	NoColor = func() bool {
		_, x := os.LookupEnv("NO_COLOR")
		return x
	}()
)

func tryPopBack(slice *[]string, value string) bool {
	if len(*slice) < 1 {
		return false
	}

	for i := len(*slice) - 1; i >= 0; i-- {
		if (*slice)[i] == value {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			return true
		}
	}

	return false
}
func findFirstTag(input string) (int, int, string) {
	indexTagStart := strings.IndexByte(input, tAG_DELIM_OPEN)
	if indexTagStart == -1 {
		return -1, -1, ""
	}

	indexTagEnd := strings.IndexByte(input[indexTagStart+1:], tAG_DELIM_CLOSE) + indexTagStart + 1
	/// see TestUnterminatedTag
	if indexTagEnd == 0 {
		indexTagEnd = len(input) - 1
	}
	if indexTagEnd == -1 {
		return -1, -1, ""
	}
	/// see TestUnterminatedTag2
	if indexTagStart+1 >= indexTagEnd {
		return -1, -1, ""
	}

	tagName := input[indexTagStart+1 : indexTagEnd]

	return indexTagStart, indexTagEnd, tagName
}
func isOpeningTag(tagName string) bool {
	return tagName[0] != tAG_CLOSE_MARKER
}
func getRGBEscapeCode(tagName string) (string, bool) {
	code := 38 // fg rgb
	/// matches "bg#hexhex"
	if tagName[0] == 'b' {
		tagName = tagName[2:] // "#hexhex"
		code = 48             // bg rgb code
	}
	if tagName[0] != '#' {
		return "", false
	}

	r, err := strconv.ParseUint(tagName[1:3], 16, 8)
	if err != nil {
		return "", false
	}
	g, err := strconv.ParseUint(tagName[3:5], 16, 8)
	if err != nil {
		return "", false
	}
	b, err := strconv.ParseUint(tagName[5:7], 16, 8)
	if err != nil {
		return "", false
	}

	/// is this actually faster than just a printf? seems negligible
	ret := strings.Builder{}
	ret.Grow(19) /// the full rgb sequence will never be longer than 19 bytes
	ret.WriteString("\x1b[")
	ret.WriteString(strconv.Itoa(code))
	ret.WriteString(";2;") /// its so sadd
	ret.WriteString(strconv.Itoa(int(r)))
	ret.WriteString(";")
	ret.WriteString(strconv.Itoa(int(g)))
	ret.WriteString(";")
	ret.WriteString(strconv.Itoa(int(b)))
	ret.WriteString("m")
	return ret.String(), true
}
func getResetEscapeCode(tagName string) (string, bool) {
	if tagName == "/" || tagName == "reset" {
		return eSCAPE_CODE_RESET, true
	} else {
		return "", false
	}
}
func getColorEscapeCode(tagName string) (string, bool) {
	if tagName[0] == '/' {
		tagName = tagName[1:]
	}
	// Try a known color name
	escapeCode, ok := colorEscapeCodes[tagName]
	if ok {
		return escapeCode, ok
	}

	// Try RGB
	escapeCode, ok = getRGBEscapeCode(tagName)
	return escapeCode, ok
}
func getItalicEscapeCode(tagName string) (string, bool) {
	escapeCode, ok := italicEscapeCodes[tagName]
	return escapeCode, ok
}
func getBoldEscapeCode(tagName string) (string, bool) {
	escapeCode, ok := boldEscapeCodes[tagName]
	return escapeCode, ok
}
func getUnderlineEscapeCode(tagName string) (string, bool) {
	escapeCode, ok := underlineEscapeCodes[tagName]
	return escapeCode, ok
}

// this might be useful, maybe for rgb?
// GetTagEscapeCode returns the ansi escape code for a tag (and its type)
func GetTagEscapeCode(tagName string) (string, TagType) {
	if len(tagName) == 0 {
		return "", TagTypeUnknown
	}
	/// the order is important!
	escapeCode, ok := getResetEscapeCode(tagName)
	if ok {
		return escapeCode, TagTypeReset
	}

	escapeCode, ok = getColorEscapeCode(tagName)
	if ok {
		return escapeCode, TagTypeColor
	}

	escapeCode, ok = getBoldEscapeCode(tagName)
	if ok {
		return escapeCode, TagTypeBold
	}

	escapeCode, ok = getItalicEscapeCode(tagName)
	if ok {
		return escapeCode, TagTypeItalic
	}

	escapeCode, ok = getUnderlineEscapeCode(tagName)
	if ok {
		return escapeCode, TagTypeUnderline
	}

	return "", TagTypeUnknown
}

func noColor() {
	NoColor = true
}
func YesColor() {
	NoColor = false
}

// process a tagged string into ansi
func ProcessTags(input string) string {
	s := strings.Builder{}
	s.Grow(len(input))

	// A list of colors that have been pushed via opening tags.
	// Closing tags will pop the most recently-pushed entry of that name from the stack
	colorEscapeCodeStack := make([]string, 2, 8)
	colorEscapeCodeStack[0] = colorEscapeCodes["bg"]
	colorEscapeCodeStack[1] = colorEscapeCodes["fg"]
	// The last-written color; used to prevent redundant writes
	lastColorEscapeCode := colorEscapeCodeStack[1]
	// A list of "decoration flags"; used to track redundant calls to decoration flag tags
	var decorationFlags decorationFlags

	inputIndexStart := 0
	for {
		currentSubstr := input[inputIndexStart:]

		// Find the indices of the start and end tag delimiters within the current substr
		substrTagIndexStart, substrTagIndexEnd, tagName := findFirstTag(currentSubstr)
		if substrTagIndexStart < 0 {
			s.WriteString(currentSubstr)
			break
		} else if substrTagIndexEnd < 0 {
			return ""
		}

		// Write all contents in the substr prior to tag open to the output string
		s.WriteString(currentSubstr[:substrTagIndexStart])

		tagEscapeCode, tagType := GetTagEscapeCode(tagName)

		if NoColor {
			tagEscapeCode = ""
		}
		switch tagType {
		case TagTypeReset:
			{
				/// ansi reset clears everything
				colorEscapeCodeStack = colorEscapeCodeStack[:2]
				if lastColorEscapeCode != tagEscapeCode {
					s.WriteString(tagEscapeCode)
					lastColorEscapeCode = tagEscapeCode
				}
			}
		case TagTypeColor:
			if isOpeningTag(tagName) {
				// Push the escape code to the stack and write it to the output if needed
				colorEscapeCodeStack = append(colorEscapeCodeStack, tagEscapeCode)
				if lastColorEscapeCode != tagEscapeCode {
					s.WriteString(tagEscapeCode)
					lastColorEscapeCode = tagEscapeCode
				}
			} else {
				// Pop the escape code from the stack
				ok := tryPopBack(&colorEscapeCodeStack, tagEscapeCode)

				if ok {
					// write most recent color
					topColorEscapeCode := colorEscapeCodeStack[len(colorEscapeCodeStack)-1]
					if lastColorEscapeCode != topColorEscapeCode {
						s.WriteString(topColorEscapeCode)
						lastColorEscapeCode = topColorEscapeCode
					}
				}
				/// otherwise ignore random closing tags
				/// MAYBE: also add to output string?
			}
		case TagTypeBold:
			// Only write escape code if we would otherwise change the state of the flag
			if isOpeningTag(tagName) && !decorationFlags.bold {
				decorationFlags.bold = true
				s.WriteString(tagEscapeCode)
			} else if !isOpeningTag(tagName) && decorationFlags.bold {
				decorationFlags.bold = false
				s.WriteString(tagEscapeCode)
			}
		case TagTypeItalic:
			// Only write escape code if we would otherwise change the state of the flag
			if isOpeningTag(tagName) && !decorationFlags.italic {
				decorationFlags.italic = true
				s.WriteString(tagEscapeCode)
			} else if !isOpeningTag(tagName) && decorationFlags.italic {
				decorationFlags.italic = false
				s.WriteString(tagEscapeCode)
			}
		case TagTypeUnderline:
			// Only write escape code if we would otherwise change the state of the flag
			if isOpeningTag(tagName) && !decorationFlags.underline {
				decorationFlags.underline = true
				s.WriteString(tagEscapeCode)
			} else if !isOpeningTag(tagName) && decorationFlags.underline {
				decorationFlags.underline = false
				s.WriteString(tagEscapeCode)
			}
		case TagTypeUnknown:
			/// just print it
			s.WriteString(currentSubstr[substrTagIndexStart : substrTagIndexEnd+1])
		}

		// Jump forward in the input string to one past the closing delimiter
		inputIndexStart += substrTagIndexEnd + 1
	}

	if !NoColor {
		s.WriteString(eSCAPE_CODE_RESET)
	}
	return s.String()
}

// prefix a string with a format tag
func TagString(s, tag string) string {
	if s == "" {
		return s
	}
	return "{" + tag + "}" + s
}

// strip format tags (NOT ansi) from line
func StripTags(input string) string {
	/// this is basically just processTags without the processing
	s := strings.Builder{}
	s.Grow(len(input))
	inputIndexStart := 0
	for {
		currentSubstr := input[inputIndexStart:]

		// Find the indices of the start and end tag delimiters within the current substr
		substrTagIndexStart, substrTagIndexEnd, _ := findFirstTag(currentSubstr)

		if substrTagIndexStart < 0 {
			s.WriteString(currentSubstr)
			break
		} else if substrTagIndexEnd < 0 {
			return ""
		}

		/// write the preceding chars...
		s.WriteString(currentSubstr[:substrTagIndexStart])
		/// then jump past
		inputIndexStart += substrTagIndexEnd + 1
	}
	return s.String()
}
