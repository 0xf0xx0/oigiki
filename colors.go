package oigiki

import (
	"strconv"
)

var colorEscapeCodes = map[string]string{
	"fg": getAnsiCode(39),
	"bg": getAnsiCode(49),

	"red":     getAnsiCode(31),
	"green":   getAnsiCode(32),
	"blue":    getAnsiCode(34),
	"yellow":  getAnsiCode(33),
	"cyan":    getAnsiCode(36),
	"magenta": getAnsiCode(35),
	"white":   getAnsiCode(37),
	"black":   getAnsiCode(30),

	"redbright":     getAnsiCode(91),
	"greenbright":   getAnsiCode(92),
	"bluebright":    getAnsiCode(94),
	"yellowbright":  getAnsiCode(93),
	"cyanbright":    getAnsiCode(96),
	"magentabright": getAnsiCode(95),
	"whitebright":   getAnsiCode(97),
	"blackbright":   getAnsiCode(90),

	"bgred":     getAnsiCode(41),
	"bggreen":   getAnsiCode(42),
	"bgblue":    getAnsiCode(44),
	"bgyellow":  getAnsiCode(43),
	"bgcyan":    getAnsiCode(46),
	"bgmagenta": getAnsiCode(45),
	"bgwhite":   getAnsiCode(47),
	"bgblack":   getAnsiCode(40),

	"bgredbright":     getAnsiCode(101),
	"bggreenbright":   getAnsiCode(102),
	"bgbluebright":    getAnsiCode(104),
	"bgyellowbright":  getAnsiCode(103),
	"bgcyanbright":    getAnsiCode(106),
	"bgmagentabright": getAnsiCode(105),
	"bgwhitebright":   getAnsiCode(107),
	"bgblackbright":   getAnsiCode(100),
}

var italicEscapeCodes = map[string]string{
	"italic":  getAnsiCode(3),
	"/italic": getAnsiCode(23),
}
var boldEscapeCodes = map[string]string{
	"bold":  getAnsiCode(1),
	"/bold": getAnsiCode(22),
}

var underlineEscapeCodes = map[string]string{
	"underline":  getAnsiCode(4),
	"/underline": getAnsiCode(24),
}

func getAnsiCode(c int) string {
	return "\x1b[" + strconv.Itoa(c) + "m"
}
