// Package oigiki is a fast ansi color tagging lib, inspired by [blessed] from node.
//
// [blessed]: https://github.com/chjj/blessed
package oigiki

import (
	"os"
	"strconv"
	"strings"
)

// Set to true to force-disable color, or false to force-enable.
// By default, it respects [NO_COLOR], [FORCE_COLOR], and [CLICOLOR].
//
// [NO_COLOR]: https://no-color.org
// [FORCE_COLOR]: https://force-color.org
// [CLICOLOR]: https://bixense.com/clicolors/
var NoColor = func() bool {
	no_color, _ := os.LookupEnv("NO_COLOR")
	force_color, _ := os.LookupEnv("FORCE_COLOR")
	clicolor_force, _ := os.LookupEnv("CLICOLOR_FORCE")
	/// "Command-line software which adds ANSI color to its output by default
	/// should check for a `NO_COLOR` environment variable that, when present
	/// and not an empty string (regardless of its value), prevents the addition of ANSI color."
	//
	/// "Command-line software which outputs colored text
	/// should check for a FORCE_COLOR environment variable. When this variable is present
	/// and not an empty string (regardless of its value), it should force the addition of ANSI color."
	//
	/// so basically NoColor = true if NO_COLOR && !(FORCE_COLOR || CLICOLOR_FORCE)
	return (no_color != "") && !(force_color != "" || clicolor_force != "")
}()

const (
	tAG_DELIM_OPEN   = '{'
	tAG_DELIM_CLOSE  = '}'
	tAG_CLOSE_MARKER = '/'
)

type decorationFlags struct {
	underline bool
	italic    bool
	bold      bool
}

// TagType is an enum of ansi tag types.
type TagType uint8

const (
	TagTypeUnknown TagType = iota
	TagTypeReset
	TagTypeColor
	TagTypeBold
	TagTypeItalic
	TagTypeUnderline
)

// GetTagEscapeCode returns the ansi escape code for a tag (and its type).
// Mainly useful for 256- and Truecolor.
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

// ProcessTags processes a tagged input string.
func ProcessTags(input string) string {
	ret := strings.Builder{}
	ret.Grow(len(input))

	// A list of colors that have been pushed via opening tags.
	// Closing tags will pop the most recently-pushed entry of that name from the relevant stack
	fgColorStack := make([]string, 1, 8)
	bgColorStack := make([]string, 1, 8)

	fgColorStack[0] = colorEscapeCodes["fg"]
	bgColorStack[0] = colorEscapeCodes["bg"]

	// The last-written color; used to prevent redundant writes
	lastFgColor := fgColorStack[0]
	lastBgColor := bgColorStack[0]

	// A list of "decoration flags"; used to track redundant calls to decoration flag tags
	var decorationFlags decorationFlags

	inputIndexStart := 0
	for {
		currentSubstr := input[inputIndexStart:]

		// Find the indices of the start and end tag delimiters within the current substr
		substrTagIndexStart, substrTagIndexEnd, tagName := findFirstTag(currentSubstr)
		if substrTagIndexStart < 0 {
			ret.WriteString(currentSubstr)
			break
		}

		// Write all contents in the substr prior to tag open to the output string
		ret.WriteString(currentSubstr[:substrTagIndexStart])

		tagEscapeCode, tagType := GetTagEscapeCode(tagName)

		switch tagType {
		case TagTypeReset:
			{
				/// reset internal state
				fgColorStack = fgColorStack[:1]
				bgColorStack = bgColorStack[:1]

				fgColorStack[0] = colorEscapeCodes["fg"]
				bgColorStack[0] = colorEscapeCodes["bg"]

				lastFgColor = fgColorStack[0]
				lastBgColor = bgColorStack[0]

				ret.WriteString(tagEscapeCode)
			}
		case TagTypeColor:
			if NoColor {
				break
			}

			if isOpeningTag(tagName) {
				stack := &fgColorStack
				lastColor := &lastFgColor
				if strings.HasPrefix(tagName, "bg") {
					stack = &bgColorStack
					lastColor = &lastBgColor
				}
				// Push the escape code to the stack and write it to the output if needed
				*stack = append(*stack, tagEscapeCode)
				if *lastColor != tagEscapeCode {
					ret.WriteString(tagEscapeCode)
					*lastColor = tagEscapeCode
				}
			} else {
				stack := &fgColorStack
				lastColor := &lastFgColor
				if strings.HasPrefix(tagName, "/bg") {
					stack = &bgColorStack
					lastColor = &lastBgColor
				}
				// Pop the escape code from the stack
				ok := tryPopBack(stack, tagEscapeCode)

				stackLen := len(*stack)
				if ok && stackLen > 0 {
					// write most recent color
					topColorEscapeCode := (*stack)[stackLen-1]
					if *lastColor != topColorEscapeCode {
						ret.WriteString(topColorEscapeCode)
						*lastColor = topColorEscapeCode
					}
				}
				/// otherwise ignore random closing tags
				/// MAYBE: also add to output string?
			}
		case TagTypeBold:
			// Only write escape code if we would otherwise change the state of the flag
			if isOpeningTag(tagName) && !decorationFlags.bold {
				decorationFlags.bold = true
				ret.WriteString(tagEscapeCode)
			} else if !isOpeningTag(tagName) && decorationFlags.bold {
				decorationFlags.bold = false
				ret.WriteString(tagEscapeCode)
			}
		case TagTypeItalic:
			// Only write escape code if we would otherwise change the state of the flag
			if isOpeningTag(tagName) && !decorationFlags.italic {
				decorationFlags.italic = true
				ret.WriteString(tagEscapeCode)
			} else if !isOpeningTag(tagName) && decorationFlags.italic {
				decorationFlags.italic = false
				ret.WriteString(tagEscapeCode)
			}
		case TagTypeUnderline:
			// Only write escape code if we would otherwise change the state of the flag
			if isOpeningTag(tagName) && !decorationFlags.underline {
				decorationFlags.underline = true
				ret.WriteString(tagEscapeCode)
			} else if !isOpeningTag(tagName) && decorationFlags.underline {
				decorationFlags.underline = false
				ret.WriteString(tagEscapeCode)
			}
		case TagTypeUnknown:
			/// just print it
			ret.WriteString(currentSubstr[substrTagIndexStart : substrTagIndexEnd+1])
		}

		// Jump forward in the input string to one past the closing delimiter
		inputIndexStart += substrTagIndexEnd + 1
	}

	ret.WriteString("\x1b[0m")

	return ret.String()
}

// TagString prefixes a string with the specified format tag.
func TagString(input, tag string) string {
	if input == "" || tag == "" {
		return input
	}
	return "{" + tag + "}" + input
}

// StripTags strips format tags (NOT ansi) from the input.
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
		}

		/// write the preceding chars...
		s.WriteString(currentSubstr[:substrTagIndexStart])
		/// then jump past
		inputIndexStart += substrTagIndexEnd + 1
	}
	return s.String()
}

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
func getResetEscapeCode(tagName string) (string, bool) {
	if tagName == "/" || tagName == "reset" {
		return "\x1b[0m", true
	} else {
		return "", false
	}
}
func get256colorEscapeCode(tagName string) (string, bool) {
	code := 38 // fg rgb
	/// matches "bg<number>"
	if tagName[0] == 'b' {
		tagName = tagName[2:] // "<number>"
		code = 48             // bg rgb code
	}
	num, err := strconv.Atoi(tagName)
	if err != nil || num < 0 || num > 255 {
		/// i dont like that the error is gobbled but oh well
		return "", false
	}
	ret := strings.Builder{}
	ret.Grow(11)
	ret.WriteString("\x1b[")
	ret.WriteString(strconv.Itoa(code))
	ret.WriteString(";5;")
	ret.WriteString(strconv.Itoa(num))
	ret.WriteString("m")
	return ret.String(), true
}
func getTruecolorEscapeCode(tagName string) (string, bool) {
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
func getColorEscapeCode(tagName string) (string, bool) {
	if tagName[0] == '/' {
		tagName = tagName[1:]
	}
	// Try a known color name
	escapeCode, ok := colorEscapeCodes[tagName]
	if ok {
		return escapeCode, ok
	}

	// Try Truecolor
	escapeCode, ok = getTruecolorEscapeCode(tagName)
	if ok {
		return escapeCode, ok
	}

	/// and try 256color
	escapeCode, ok = get256colorEscapeCode(tagName)
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
