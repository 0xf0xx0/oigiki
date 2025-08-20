package oigiki

import (
	"strconv"

	"github.com/fatih/color"
)

var colorEscapeCodes = map[string]string{
	"fg":    getAnsiCode(39),
	"bg":    getAnsiCode(49),

	"red":     getAnsiCode(color.FgRed),
	"green":   getAnsiCode(color.FgGreen),
	"blue":    getAnsiCode(color.FgBlue),
	"yellow":  getAnsiCode(color.FgYellow),
	"cyan":    getAnsiCode(color.FgCyan),
	"magenta": getAnsiCode(color.FgMagenta),
	"white":   getAnsiCode(color.FgWhite),
	"black":   getAnsiCode(color.FgBlack),

	"redbright":     getAnsiCode(color.FgHiRed),
	"greenbright":   getAnsiCode(color.FgHiGreen),
	"bluebright":    getAnsiCode(color.FgHiBlue),
	"yellowbright":  getAnsiCode(color.FgHiYellow),
	"cyanbright":    getAnsiCode(color.FgHiCyan),
	"magentabright": getAnsiCode(color.FgHiMagenta),
	"whitebright":   getAnsiCode(color.FgHiWhite),
	"blackbright":   getAnsiCode(color.FgHiBlack),

	"bgred":     getAnsiCode(color.BgRed),
	"bggreen":   getAnsiCode(color.BgGreen),
	"bgblue":    getAnsiCode(color.BgBlue),
	"bgyellow":  getAnsiCode(color.BgYellow),
	"bgcyan":    getAnsiCode(color.BgCyan),
	"bgmagenta": getAnsiCode(color.BgMagenta),
	"bgwhite":   getAnsiCode(color.BgWhite),
	"bgblack":   getAnsiCode(color.BgBlack),

	"bgredbright":     getAnsiCode(color.BgHiRed),
	"bggreenbright":   getAnsiCode(color.BgHiGreen),
	"bgbluebright":    getAnsiCode(color.BgHiBlue),
	"bgyellowbright":  getAnsiCode(color.BgHiYellow),
	"bgcyanbright":    getAnsiCode(color.BgHiCyan),
	"bgmagentabright": getAnsiCode(color.BgHiMagenta),
	"bgwhitebright":   getAnsiCode(color.BgHiWhite),
	"bgblackbright":   getAnsiCode(color.BgHiBlack),
}

var italicEscapeCodes = map[string]string{
	"italic":  getAnsiCode(color.Italic),
	"/italic": getAnsiCode(color.ResetItalic),
}
var boldEscapeCodes = map[string]string{
	"bold":  getAnsiCode(color.Bold),
	"/bold": getAnsiCode(color.ResetBold),
}

var underlineEscapeCodes = map[string]string{
	"underline":  getAnsiCode(color.Underline),
	"/underline": getAnsiCode(color.ResetUnderline),
}

func getAnsiCode(c color.Attribute) string {
	return "\x1b[" + strconv.Itoa(int(c)) + "m"
}
