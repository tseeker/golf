package golf

import (
	"errors"
	"strings"
	"testing"
)

func ensurePanic(t *testing.T, errorMessage string, f func()) {
	t.Run(errorMessage, func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil || r.(error).Error() != errorMessage {
				t.Errorf("GOT: %v; WANT: %v", r, errorMessage)
			}
		}()
		f()
	})
}

func TestParseEmpty(t *testing.T) {
	resetParser()
	if got, want := parse(""), error(nil); got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestParseUnknownShortOption(t *testing.T) {
	resetParser()
	got, want := parse("-a"), errors.New("unknown option: 'a'")
	if got == nil || got.Error() != want.Error() {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestParseUnknownLongOption(t *testing.T) {
	resetParser()
	got, want := parse("--version"), errors.New("unknown option: \"version\"")
	if got == nil || got.Error() != want.Error() {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestParseComplex(t *testing.T) {
	resetParser()

	a := Int("l", "limit", 0, "limit results")
	b := Bool("v", "verbose", false, "print verbose info")
	c := String("s", "servers", "", "ask servers")

	if got, want := parse("-l 4 -v -s host1,host2 some other arguments"), error(nil); got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := *a, 4; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := *b, true; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := *c, "host1,host2"; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	// FIXME difficult to test when Args() invokes os.Args...

	// fmt.Fprintf(os.Stderr, "%#v\n", Args())

	// if got, want := strings.Join(Args(), " "), "some other arguments"; got != want {
	// 	t.Errorf("GOT: %v; WANT: %v", got, want)
	// }
}

// TODO: test for `--` to stop parsing

func TestParseStopsAfterDoubleHyphen(t *testing.T) {
	resetParser()

	a := Int("l", "limit", 0, "limit results")
	b := Bool("v", "verbose", false, "print verbose info")

	if got, want := parse("-l4 -- --verbose some other arguments"), error(nil); got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := *a, 4; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := *b, false; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestParseConfused(t *testing.T) {
	t.Skip("only works with long command line")

	resetParser()

	a := Int("l", "limit", 0, "limit results")
	b := Bool("v", "verbose", false, "print verbose info")

	// ??? It would be nice to test for exact string match to automate ensuring
	// proper argument is returned in the error message.
	if got, want := parse("-vl"), "cannot parse argument"; !strings.HasPrefix(got.Error(), want) {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := *a, 0; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := *b, true; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestParseHyphenAfterShort(t *testing.T) {
	t.Skip("only works with long command line")

	resetParser()

	a := Int("l", "limit", 0, "limit results")
	b := Bool("v", "verbose", false, "print verbose info")

	// ??? It would be nice to test for exact string match to automate ensuring
	// proper argument is returned in the error message.
	if got, want := parse("-v-l"), "cannot parse argument"; !strings.HasPrefix(got.Error(), want) {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := *a, 0; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := *b, true; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestPanicsWhenAttemptToRedefineFlag(t *testing.T) {
	ensurePanic(t, "cannot add option that duplicates short flag: 'f'", func() {
		_ = Uint("f", "flubber", 0, "some example flag")
		_ = Uint("f", "blubber", 0, "some example flag")
	})

	ensurePanic(t, "cannot add option that duplicates long flag: \"flubber\"", func() {
		_ = Uint("f", "flubber", 0, "some example flag")
		_ = Uint("b", "flubber", 0, "some example flag")
	})
}
