// ansi color tagging system similar to blessed from node
package oigiki

import (
	"regexp"
	"slices"
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
type TagType int

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
	tagreg = regexp.MustCompile(`(\{.+?\})`)
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

	tagName := input[indexTagStart+1 : indexTagEnd]

	return indexTagStart, indexTagEnd, tagName
}
func isOpeningTag(tagName string) bool {
	return tagName[0] != tAG_CLOSE_MARKER
}

func validateRgbString(input string) bool {
	if input[0] != '#' {
		return false
	}
	l := len(input)
	if l != 7 {
		return false
	}

	for i := 1; i < l; i++ {
		char := input[i]
		if !(char >= 'a' && char <= 'f') && !(char >= '0' && char <= '9') {
			return false
		}
	}

	return true
}
func getRGBEscapeCode(tagName string) (string, bool) {
	code := 38 // fg rgb
	/// matches "bg#hexhex"
	if tagName[0] == 'b' {
		tagName = tagName[2:]
		code = 48 // bg rgb
	}
	if !validateRgbString(tagName) {
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

	ret := strings.Builder{}
	ret.Grow(19) /// the full rgb sequence will never be longer than 19 bytes
	ret.WriteString("\x1b[")
	ret.WriteString(strconv.Itoa(code))
	ret.WriteString(";2;")
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
// returns the ansi escape code for a tag (and its type)
func GetTagEscapeCode(tagName string) (string, TagType) {
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

// process a tagged string into ansi
func ProcessTags(input string) string {
	s := strings.Builder{}

	// A list of colors that have been pushed via opening tags. Closing tags will pop the most recently-pushed entry of that name from the stack
	//
	var colorEscapeCodeStack = []string{colorEscapeCodes["bg"], colorEscapeCodes["fg"]}
	// The last-written color; used to prevent redundant writes
	var lastColorEscapeCode = colorEscapeCodes["fg"]
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
		switch tagType {
		case TagTypeReset:
			{
				s.WriteString(tagEscapeCode)
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
			/// ignore it
		}

		// Jump forward in the input string to one past the closing delimiter
		inputIndexStart += substrTagIndexEnd + 1
	}

	s.WriteString(eSCAPE_CODE_RESET)
	return s.String()
}

// add a format tag to a string
func TagString(s, tag string) string {
	if s == "" {
		return s
	}
	return "{" + tag + "}" + s
}

// strip format tags (NOT ansi) from line
func StripTags(s string) string {
	for _, tag := range slices.Backward(tagreg.FindAllStringIndex(s, -1)) {
		s = s[:tag[0]] + s[tag[1]:]
	}
	return s
}
