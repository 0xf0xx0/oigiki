package oigiki_test

import (
	"testing"

	"git.0xf0xx0.eth.limo/0xf0xx0/oigiki"
)

func TestFastMain(t *testing.T) {
	/// just so its pipable
	oigiki.NoColor = false
}
func TestFastFG(t *testing.T) {
	str := "{red}red{green}green{blue}blue"
	EXPECTED := "\x1b[31mred\x1b[32mgreen\x1b[34mblue\x1b[0m"
	logFast(t, str, EXPECTED, false)
}
func TestFastBG(t *testing.T) {
	str := "{white}{bgred}oi{black}{bggreen}gi{white}{bgblue}ki{/}{yellow}!"
	EXPECTED := "\x1b[37m\x1b[41moi\x1b[30m\x1b[42mgi\x1b[37m\x1b[44mki\x1b[0m\x1b[33m!\x1b[0m"
	logFast(t, str, EXPECTED, false)
}

func TestFastMixedFGBG(t *testing.T) {
	str := "{bgblack}{white}mixed {bgblue}bg{/bgblue} and {bgyellow}{black}fg"
	EXPECTED := "\x1b[40m\x1b[37mmixed \x1b[44mbg and \x1b[43m\x1b[30mfg\x1b[0m"
	logFast(t, str, EXPECTED, false)
}

func TestFast256color(t *testing.T) {
	str := "{1}van {2}two {3}three, ah ah ah"
	EXPECTED := "\x1b[38;5;1mvan \x1b[38;5;2mtwo \x1b[38;5;3mthree, ah ah ah\x1b[0m"
	logFast(t, str, EXPECTED, false)
}
func TestFast256colorBG(t *testing.T) {
	str := "{bg0}{28}she bg on my {bg4}fg{/bg4} till i {bg3}wg"
	EXPECTED := "\x1b[48;5;0m\x1b[38;5;28mshe bg on my \x1b[48;5;4mfg till i \x1b[48;5;3mwg\x1b[0m"
	logFast(t, str, EXPECTED, false)
}
func TestFast256colorIndexOOB(t *testing.T) {
	str := "{132453}index out of {-123}bounds"
	EXPECTED := "{132453}index out of {-123}bounds\x1b[0m"
	logFast(t, str, EXPECTED, false)
}

func TestFastTruecolor(t *testing.T) {
	str := "{#b00b69}#b00b69{/#b00b69}"
	EXPECTED := "\x1b[38;2;176;11;105m#b00b69\x1b[0m"
	logFast(t, str, EXPECTED, false)
}
func TestFastTruecolorBG(t *testing.T) {
	str := "{#ffffff}{bg#ff0000}oi{#000000}{bg#00ff00}gi{#ffffff}{bg#0000ff}ki{/}{#ffff00}!"
	EXPECTED := "\x1b[38;2;255;255;255m\x1b[48;2;255;0;0moi\x1b[38;2;0;0;0m\x1b[48;2;0;255;0mgi\x1b[38;2;255;255;255m\x1b[48;2;0;0;255mki\x1b[0m\x1b[38;2;255;255;0m!\x1b[0m"
	logFast(t, str, EXPECTED, false)
}

/// edge cases

func TestFastNoTags(t *testing.T) {
	str := "notags"
	EXPECTED := "notags\x1b[0m"
	logFast(t, str, EXPECTED, false)
}

func TestFastEmpty(t *testing.T) {
	str := ""
	EXPECTED := "\x1b[0m"
	logFast(t, str, EXPECTED, false)
}

func TestFastEmptyTag(t *testing.T) {
	str := "{blue}blue{}blue"
	EXPECTED := "\x1b[34mblue{}blue\x1b[0m"
	logFast(t, str, EXPECTED, false)
}
func TestFastInvalidTag(t *testing.T) {
	str := "{blue}blue{snuffleupagus}blue"
	EXPECTED := "\x1b[34mblue{snuffleupagus}blue\x1b[0m"
	logFast(t, str, EXPECTED, false)
}

func TestFastExtraTagOpeners(t *testing.T) {
	str := "{{red}{{green}{{blue}{"
	EXPECTED := "{{red}{{green}{{blue}{\x1b[0m"
	logFast(t, str, EXPECTED, false)
}

func TestFastExtraTagTerminators(t *testing.T) {
	str := "{red}}{green}}{blue}}"
	EXPECTED := "\x1b[31m}\x1b[32m}\x1b[34m}\x1b[0m"
	logFast(t, str, EXPECTED, false)
}
func TestFastUnterminatedTag(t *testing.T) {
	str := "{red"
	EXPECTED := "{red\x1b[0m"
	logFast(t, str, EXPECTED, false)
}
func TestFastUnterminatedTag2(t *testing.T) {
	str := "trigger{red"
	EXPECTED := "trigger{red\x1b[0m"
	logFast(t, str, EXPECTED, false)
}

func TestFastResetEdge(t *testing.T) {
	str := "{bggreen}green{blue}blue{yellow}yellow{/}reset{bgblue}blue"
	EXPECTED := "\x1b[42mgreen\x1b[34mblue\x1b[33myellow\x1b[0mreset\x1b[44mblue\x1b[0m"
	logFast(t, str, EXPECTED, false)
}

func TestFastRandomClosingTag(t *testing.T) {
	str := "{red}red{/blue}red"
	EXPECTED := "\x1b[31mredred\x1b[0m"
	logFast(t, str, EXPECTED, false)
}

func BenchmarkProcessTagsFast(b *testing.B) {
	for b.Loop() {
		oigiki.ProcessTagsFast(benchString)
	}
}
