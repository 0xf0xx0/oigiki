package oigiki_test

import (
	"testing"

	"github.com/0xf0xx0/oigiki"
	"github.com/fatih/color"
)

func TestProcessTags(t *testing.T) {
	color.NoColor = false

	str := oigiki.ProcessTags("{red}red{green}green{blue}blue")
	EXPECTED := "\x1b[31mred\x1b[32mgreen\x1b[34mblue\x1b[0m"
	t.Log(str)
	if str != EXPECTED {
		t.Errorf("MISMATCH: %q %q", EXPECTED, str)
	}
}

func TestUnset(t *testing.T) {
	str := oigiki.ProcessTags("{red}red{green}green{/green}red")
	EXPECTED := "\x1b[31mred\x1b[32mgreen\x1b[31mred\x1b[0m"
	t.Log(str)
	if str != EXPECTED {
		t.Errorf("MISMATCH: %q %q", EXPECTED, str)
	}
}

func TestMixed(t *testing.T) {
	str := oigiki.ProcessTags("{green}green{bold}boldgreen{blue}{/bold}blue{underline}blueunderline")
	EXPECTED := "\x1b[32mgreen\x1b[1mboldgreen\x1b[34m\x1b[22mblue\x1b[4mblueunderline\x1b[0m"
	t.Log(str)
	if str != EXPECTED {
		t.Errorf("MISMATCH: %q %q", EXPECTED, str)
	}
}

func TestRGB(t *testing.T) {
	str := oigiki.ProcessTags("{green}green{bold}boldgreen{#b00b69}{/bold}#b00b69{/#b00b69}greenunderline")
	EXPECTED := ""
	t.Log(str)
	if str != EXPECTED {
		t.Errorf("MISMATCH: %q %q", EXPECTED, str)
	}
}

func TestHell(t *testing.T) {
	str := oigiki.ProcessTags("{red}red{bold}boldred{underline}boldredunderline{/bold}redunderline{/red}underline{/underline}")
	EXPECTED := "\x1b[31mred\x1b[1mboldred\x1b[4mboldredunderline\x1b[22mredunderline\x1b[39munderline\x1b[24m\x1b[0m"
	t.Log(str)
	if str != EXPECTED {
		t.Errorf("MISMATCH: %q %q", EXPECTED, str)
	}
}
func TestFunnyQuirk(t *testing.T) {
	str := oigiki.ProcessTags("{green}1 text {blue}2 test {red}red text{/blue}bloo... wait no green?")
	EXPECTED := "\x1b[32m1 text \x1b[34m2 test \x1b[31mred text\x1b[32mbloo... wait no green?\x1b[0m"
	t.Log(str)
	if str != EXPECTED {
		t.Errorf("MISMATCH: %q %q", EXPECTED, str)
	}
}
