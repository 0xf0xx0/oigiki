package oigiki_test

import (
	"testing"

	"github.com/0xf0xx0/oigiki"
	"github.com/fatih/color"
)

func TestProcessTags(t *testing.T) {
	color.NoColor = false

	str := oigiki.ProcessTags("{red}red{green}green{blue}blue{/}")
	EXPECTED := "\x1b[31mred\x1b[32mgreen\x1b[34mblue\x1b[0m"
	t.Log(str)
	if str != EXPECTED {
		t.Errorf("MISMATCH: %q %q", EXPECTED, str)
	}
}

func TestUnset(t *testing.T) {
	str := oigiki.ProcessTags("{red}red{green}green{/green}red{/}")
	EXPECTED := "\x1b[31mred\x1b[32mgreen\x1b[31mred\x1b[0m"
	t.Log(str)
	if str != EXPECTED {
		t.Errorf("MISMATCH: %q %q", EXPECTED, str)
	}
}

func TestFailedUnset(t *testing.T) {
	str := oigiki.ProcessTags("unset{green}green{/green}unset{/}")
	EXPECTED := "unset\x1b[32mgreenunset\x1b[0m"
	t.Log(str)
	if str != EXPECTED {
		t.Errorf("MISMATCH: %q %q", EXPECTED, str)
	}
}
func TestMixed(t *testing.T) {
	str := oigiki.ProcessTags("{green}green{bold}boldgreen{blue}{/bold}blue{underline}blueunderline{/}")
	EXPECTED := "\x1b[32mgreen\x1b[1mboldgreen\x1b[34m\x1b[22mblue\x1b[4mblueunderline\x1b[0m"
	t.Log(str)
	if str != EXPECTED {
		t.Errorf("MISMATCH: %q %q", EXPECTED, str)
	}
}
