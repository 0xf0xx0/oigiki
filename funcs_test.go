package oigiki_test

import (
	"testing"

	"github.com/0xf0xx0/oigiki"
	"github.com/fatih/color"
)

func TestMain(t *testing.T) {
	/// just so its pipable
	color.NoColor = false
}
func TestProcessTags(t *testing.T) {
	str := "{red}red{green}green{blue}blue"
	EXPECTED := "\x1b[31mred\x1b[32mgreen\x1b[34mblue\x1b[0m"
	log(t, str, EXPECTED, false)
}

func TestUnset(t *testing.T) {
	str := "{red}red{green}green{/green}red"
	EXPECTED := "\x1b[31mred\x1b[32mgreen\x1b[31mred\x1b[0m"
	log(t, str, EXPECTED, false)
}

func TestUnsetNesting(t *testing.T) {
	str := "{red}red{green}{blue}{underline}blue{/underline}{/green}{/red}blue"
	EXPECTED := "\x1b[31mred\x1b[32m\x1b[34m\x1b[4mblue\x1b[24mblue\x1b[0m"
	log(t, str, EXPECTED, false)
}

func TestDefaultUnset(t *testing.T) {
	str := "{red}red{blue}blue{/blue}red{/red}default"
	EXPECTED := "\x1b[31mred\x1b[34mblue\x1b[31mred\x1b[39mdefault\x1b[0m"
	log(t, str, EXPECTED, false)
}

func TestOverlap(t *testing.T) {
	str := "{green}green{bold}boldgreen{blue}{/bold}blue{underline}blueunderline"
	EXPECTED := "\x1b[32mgreen\x1b[1mboldgreen\x1b[34m\x1b[22mblue\x1b[4mblueunderline\x1b[0m"
	log(t, str, EXPECTED, false)
}

func TestRGB(t *testing.T) {
	str := "{#b00b69}#b00b69{/#b00b69}"
	EXPECTED := "\x1b[38;2;176;11;105m#b00b69\x1b[39m\x1b[0m"
	log(t, str, EXPECTED, false)
}
func TestBG(t *testing.T) {
	str := "{black}{bgred}oi{bggreen}gi{bgblue}ki{/}{yellow}!"
	EXPECTED := "\x1b[30m\x1b[41moi\x1b[42mgi\x1b[44mki\x1b[0m\x1b[33m!\x1b[0m"
	log(t, str, EXPECTED, false)
}
func TestRGBBG(t *testing.T) {
	str := "{#000000}{bg#ff0000}oi{bg#00ff00}gi{bg#0000ff}ki{/}{#ffff00}!"
	EXPECTED := "\x1b[38;2;0;0;0m\x1b[48;2;255;0;0moi\x1b[48;2;0;255;0mgi\x1b[48;2;0;0;255mki\x1b[0m\x1b[38;2;255;255;0m!\x1b[0m"
	log(t, str, EXPECTED, false)
}

func TestHell(t *testing.T) {
	str := "{red}red{bold}boldred{underline}boldredunderline{/bold}redunderline{/red}fgunderline{/underline}"
	EXPECTED := "\x1b[31mred\x1b[1mboldred\x1b[4mboldredunderline\x1b[22mredunderline\x1b[39mfgunderline\x1b[24m\x1b[0m"
	log(t, str, EXPECTED, false)
}

func TestResetEdgeCase(t *testing.T) {
	str := "{red}red{reset}reset{blue}blue{/blue}fg"
	EXPECTED := "\x1b[31mred\x1b[0mreset\x1b[34mblue\x1b[39mfg\x1b[0m"
	log(t, str, EXPECTED, false)
}

func TestExtraTagTerminators(t *testing.T) {
	str := "{red}}{green}}{blue}}"
	EXPECTED := "\x1b[31m}\x1b[32m}\x1b[34m}\x1b[0m"
	log(t, str, EXPECTED, false)
}
func TestUnterminatedTag(t *testing.T) {
	str := "{red"
	EXPECTED := "\x1b[0m"
	log(t, str, EXPECTED, false)
}
func TestRandomClosingTag(t *testing.T) {
	str := "{red}red{/blue}red"
	EXPECTED := "\x1b[31mredred\x1b[0m"
	log(t, str, EXPECTED, false)
}
func TestInvalidTag(t *testing.T) {
	str := "{blue}blue{snuffleupagus}blue"
	EXPECTED := "\x1b[34mblueblue\x1b[0m"
	log(t, str, EXPECTED, false)
}
func TestRedundantTag(t *testing.T) {
	str := "{green}green{green}green{/green}"
	EXPECTED := "\x1b[32mgreengreen\x1b[0m"
	log(t, str, EXPECTED, false)
}

func TestStripTag(t *testing.T) {
	str := oigiki.StripTags("{red}red{green}green{blue}blue")
	EXPECTED := "redgreenblue"
	log(t, str, EXPECTED, true)
}

func TestTagString(t *testing.T) {
	str := oigiki.TagString("this should be red", "red")
	EXPECTED := "{red}this should be red"
	log(t, str, EXPECTED, true)
}

func BenchmarkProcessing(b *testing.B) {
	str := "{red}red{bold}boldred{underline}bold{green}greenunderline{/bold}{red}redunderline{/red}greenunderline{/underline}" +
		"green{/}reset reset reset{cyan}[{yellow}{#b00b69}{/bold}#b00b69{/#b00b69}yellow]"
	for b.Loop() {
		oigiki.ProcessTags(str)
	}
}

func log(t *testing.T, str string, EXPECTED string, skipprocess bool) {
	t.Logf("input: %q", str)
	processed := str
	if !skipprocess {
		processed = oigiki.ProcessTags(str)
	}
	t.Logf("output: %s", processed)
	if processed != EXPECTED {
		t.Errorf("MISMATCH:\nEXP: %s %q\nGOT: %s %q", EXPECTED, EXPECTED, processed, processed)
	}
}
