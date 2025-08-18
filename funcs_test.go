package oigiki_test

import (
	"testing"

	"github.com/0xf0xx0/oigiki"
	"github.com/fatih/color"
)

func TestInit(t *testing.T) {
	/// just so its pipable
	color.NoColor = false
}
func TestProcessTags(t *testing.T) {
	str := "{red}red{green}green{blue}blue"
	EXPECTED := "\x1b[31mred\x1b[32mgreen\x1b[34mblue\x1b[0m"
	log(t, str, EXPECTED)
}

func TestUnset(t *testing.T) {
	str := "{red}red{green}green{/green}red"
	EXPECTED := "\x1b[31mred\x1b[32mgreen\x1b[31mred\x1b[0m"
	log(t, str, EXPECTED)
}

func TestUnset2(t *testing.T) {
	str := "{red}red{green}{blue}{underline}blue{/underline}{/green}{/red}blue"
	EXPECTED := "\x1b[31mred\x1b[32m\x1b[34m\x1b[4mblue\x1b[24mblue\x1b[0m"
	log(t, str, EXPECTED)
}

func TestDefaultUnset(t *testing.T) {
	str := "{red}red{blue}blue{/blue}red{/red}default"
	EXPECTED := "\x1b[31mred\x1b[34mblue\x1b[31mred\x1b[39mdefault\x1b[0m"
	log(t, str, EXPECTED)
}

func TestMixed(t *testing.T) {
	str := "{green}green{bold}boldgreen{blue}{/bold}blue{underline}blueunderline"
	EXPECTED := "\x1b[32mgreen\x1b[1mboldgreen\x1b[34m\x1b[22mblue\x1b[4mblueunderline\x1b[0m"
	log(t, str, EXPECTED)
}

func TestRGB(t *testing.T) {
	str := "{green}green{bold}boldgreen{#b00b69}{/bold}#b00b69{underline}{/#b00b69}greenunderline"
	EXPECTED := "\x1b[32mgreen\x1b[1mboldgreen\x1b[38;2;176;11;105m\x1b[22m#b00b69\x1b[4m\x1b[32mgreenunderline\x1b[0m"
	log(t, str, EXPECTED)
}

func TestHell(t *testing.T) {
	str := "{red}red{bold}boldred{underline}boldredunderline{/bold}redunderline{/red}fgunderline{/underline}"
	EXPECTED := "\x1b[31mred\x1b[1mboldred\x1b[4mboldredunderline\x1b[22mredunderline\x1b[39mfgunderline\x1b[24m\x1b[0m"
	log(t, str, EXPECTED)
}

func TestRobustness(t *testing.T) {
	str := "{red}}{green}}{blue}}"
	EXPECTED := "\x1b[31m}\x1b[32m}\x1b[34m}\x1b[0m"
	log(t, str, EXPECTED)
}

func BenchmarkProcessing(b *testing.B) {
	str := "{red}red{bold}boldred{underline}bold{green}greenunderline{/bold}{red}redunderline{/red}greenunderline{/underline}"+
	"green{/}reset reset reset{cyan}[{yellow}{#b00b69}{/bold}#b00b69{/#b00b69}yellow]"
	for b.Loop() {
		oigiki.ProcessTags(str)
	}
}

func log(t *testing.T, str string, EXPECTED string) {
	t.Logf("input: %q", str)
	processed := oigiki.ProcessTags(str)
	t.Log(processed)
	if processed != EXPECTED {
		t.Errorf("MISMATCH:\nEXPECTED: %s %q\nGOT: %s %q", EXPECTED, EXPECTED, processed, processed)
	}
}
