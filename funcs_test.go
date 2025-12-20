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

func TestNoTags(t *testing.T) {
	str := "notags"
	EXPECTED := "notags\x1b[0m"
	log(t, str, EXPECTED, false)
}

func TestEmpty(t *testing.T) {
	str := ""
	EXPECTED := "\x1b[0m"
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
	str := "{white}{bgred}oi{black}{bggreen}gi{white}{bgblue}ki{/}{yellow}!"
	EXPECTED := "\x1b[37m\x1b[41moi\x1b[30m\x1b[42mgi\x1b[37m\x1b[44mki\x1b[0m\x1b[33m!\x1b[0m"
	log(t, str, EXPECTED, false)
}
func TestRGBBG(t *testing.T) {
	str := "{#ffffff}{bg#ff0000}oi{#000000}{bg#00ff00}gi{#ffffff}{bg#0000ff}ki{/}{#ffff00}!"
	EXPECTED := "\x1b[38;2;255;255;255m\x1b[48;2;255;0;0moi\x1b[38;2;0;0;0m\x1b[48;2;0;255;0mgi\x1b[38;2;255;255;255m\x1b[48;2;0;0;255mki\x1b[0m\x1b[38;2;255;255;0m!\x1b[0m"
	log(t, str, EXPECTED, false)
}

func TestHell(t *testing.T) {
	str := "{red}red{bold}boldred{underline}boldredunderline{/bold}redunderline{/red}fgunderline{/underline}"
	EXPECTED := "\x1b[31mred\x1b[1mboldred\x1b[4mboldredunderline\x1b[22mredunderline\x1b[39mfgunderline\x1b[24m\x1b[0m"
	log(t, str, EXPECTED, false)
}

func TestEmptyColorStack(t *testing.T) {
	str := "{/bg}{/fg}{/blue}"
	EXPECTED := "\x1b[0m"
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
	EXPECTED := "{red\x1b[0m"
	log(t, str, EXPECTED, false)
}
func TestUnterminatedTag2(t *testing.T) {
	str := "trigger{red"
	EXPECTED := "trigger{red\x1b[0m"
	log(t, str, EXPECTED, false)
}
func TestRandomClosingTag(t *testing.T) {
	str := "{red}red{/blue}red"
	EXPECTED := "\x1b[31mredred\x1b[0m"
	log(t, str, EXPECTED, false)
}
func TestInvalidTag(t *testing.T) {
	str := "{blue}blue{snuffleupagus}blue{}"
	EXPECTED := "\x1b[34mblue{snuffleupagus}blue{}\x1b[0m"
	log(t, str, EXPECTED, false)
}
func TestEmptyTag(t *testing.T) {
	str := "{blue}blue{}blue"
	EXPECTED := "\x1b[34mblue{}blue\x1b[0m"
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
func TestTagEmptyString(t *testing.T) {
	str := oigiki.TagString("", "red")
	EXPECTED := ""
	log(t, str, EXPECTED, true)
}

func TestGetCode(t *testing.T) {
	str, tag := oigiki.GetTagEscapeCode("bold")
	if str != "\x1b[1m" || tag != oigiki.TagTypeBold {
		t.Fatalf("bold isnt bold (str: %q, tag: %d)", str, tag)
	}

	str, tag = oigiki.GetTagEscapeCode("italic")
	if str != "\x1b[3m" || tag != oigiki.TagTypeItalic {
		t.Fatalf("italic isnt italic (str: %q, tag: %d)", str, tag)
	}

	str, tag = oigiki.GetTagEscapeCode("underline")
	if str != "\x1b[4m" || tag != oigiki.TagTypeUnderline {
		t.Fatalf("underline isnt underline (str: %q, tag: %d)", str, tag)
	}

	str, tag = oigiki.GetTagEscapeCode("red")
	if str != "\x1b[31m" || tag != oigiki.TagTypeColor {
		t.Fatalf("red isnt red (str: %q, tag: %d)", str, tag)
	}

	str, tag = oigiki.GetTagEscapeCode("#ff00ff")
	if str != "\x1b[38;2;255;0;255m" || tag != oigiki.TagTypeColor {
		t.Fatalf("rgb isnt rgb (str: %q, tag: %d)", str, tag)
	}

	str, tag = oigiki.GetTagEscapeCode("bg#ff00ff")
	if str != "\x1b[48;2;255;0;255m" || tag != oigiki.TagTypeColor {
		t.Fatalf("bg rgb isnt bg rgb (str: %q, tag: %d)", str, tag)
	}


	str, tag = oigiki.GetTagEscapeCode("#LL00ff")
	if str != "" || tag != oigiki.TagTypeUnknown {
		t.Fatalf("invalid rgb should be invalid (furst byte) (str: %q, tag: %d)", str, tag)
	}

	str, tag = oigiki.GetTagEscapeCode("#ffLLff")
	if str != "" || tag != oigiki.TagTypeUnknown {
		t.Fatalf("invalid rgb should be invalid (second byte) (str: %q, tag: %d)", str, tag)
	}

	str, tag = oigiki.GetTagEscapeCode("#ff00LL")
	if str != "" || tag != oigiki.TagTypeUnknown {
		t.Fatalf("invalid rgb should be invalid (third byte) (str: %q, tag: %d)", str, tag)
	}

	str, tag = oigiki.GetTagEscapeCode("")
	if str != "" || tag != oigiki.TagTypeUnknown {
		t.Fatalf("the empty tag should be unknown (str: %q, tag: %d)", str, tag)
	}
}

func TestNoColor(t *testing.T) {
	oigiki.NoColor = true
	str := oigiki.ProcessTags("{red}this {italic}should{/italic} be plain")
	oigiki.NoColor = false
	EXPECTED := "this \x1b[3mshould\x1b[23m be plain\x1b[0m"
	log(t, str, EXPECTED, true)
}

func FuzzTagging(f *testing.F) {
	testcases := []string{
		"{red}red{green}green{blue}blue",
		"{red}{bold}boldred{/bold}red",
		"{}invalidtag{invalidtag}{/}{ momo ompos pmo}",
		"{rwefwge{rfwgege{WGgwggeg{wgwgg}wg}wgrwgg{wwgwgr",
		"trigger{red",
		"}wgwehhwg{",
	}
	for _, tc := range testcases {
        f.Add(tc)
    }
    f.Fuzz(func(t *testing.T, a string) {

    })
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
