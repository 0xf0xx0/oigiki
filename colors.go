package oigiki

import (
	"github.com/fatih/color"
)

var colorMap = map[string]func(a ...interface{}) string{
	"bold":      color.New(color.Bold).SprintFunc(),
	"underline": color.New(color.Underline).SprintFunc(),
	"italic":    color.New(color.Italic).SprintFunc(),

	"/bold":      color.New(color.ResetBold).SprintFunc(),
	"/italic":    color.New(color.ResetItalic).SprintFunc(),
	"/underline": color.New(color.ResetUnderline).SprintFunc(),
	"reset":      color.New(color.Reset).SprintFunc(),

	"fg": color.New(39).SprintFunc(),
	"bg": color.New(49).SprintFunc(),

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
