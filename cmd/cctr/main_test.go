package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestProcessLines(t *testing.T) {
	t.Run("Sub small and capital c-letters", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			input:       strings.NewReader("Coding Challenges"),
			target:      []rune("C"),
			translation: []rune("c"),
			output:      &buf,
		}

		cfg.translateCmd()

		got := buf.String()
		want := "coding challenges"

		assertEqual(t, got, want)
	})

	t.Run("Sub various letters and numbers, multiline", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			input:       strings.NewReader("Coding Challenges\nhello123"),
			target:      []rune("lo12"),
			translation: []rune("bo34"),
			output:      &buf,
		}

		cfg.translateCmd()

		got := buf.String()
		want := "Coding Chabbenges\nhebbo343"

		assertEqual(t, got, want)
	})

	t.Run("emoji rune subst", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			input:       strings.NewReader("hello😊"),
			target:      []rune("😊"),
			translation: []rune("👀"),
			output:      &buf,
		}

		cfg.translateCmd()

		got := buf.String()
		want := "hello👀"

		assertEqual(t, got, want)
	})

	t.Run("range expression", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			input:       strings.NewReader("Coding Challenges\nhello123æøå"),
			target:      []rune("a-d"),
			translation: []rune("e-h"),
			output:      &buf,
		}

		cfg.translateCmd()

		got := buf.String()
		want := "Cohing Chellenges\nhello123æøå"

		assertEqual(t, got, want)
	})

	t.Run("range expression, mixed", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			input:       strings.NewReader("abcdefghijklmnop"),
			target:      []rune("abc-f"),
			translation: []rune("ghi-l"),
			output:      &buf,
		}

		cfg.translateCmd()

		got := buf.String()
		want := "ghijklghijklmnop"

		assertEqual(t, got, want)
	})

	t.Run("range expression ignored", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			input:       strings.NewReader("Coding Challenges\nhello123æøå"),
			target:      []rune("d-a"),
			translation: []rune("h-e"),
			output:      &buf,
		}

		cfg.translateCmd()

		got := buf.String()
		want := "Cohing Chellenges\nhello123æøå"

		assertEqual(t, got, want)
	})

	t.Run("special chars", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			input:       strings.NewReader("Coding =%098*\nhello123æøå"),
			target:      []rune("æ%"),
			translation: []rune("å=."),
			output:      &buf,
		}

		cfg.translateCmd()

		got := buf.String()
		want := "Coding ==098*\nhello123åøå"

		assertEqual(t, got, want)
	})

	t.Run("class specifier lower to upper", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			input:       strings.NewReader("Coding Challenge"),
			target:      []rune("[:lower:]"),
			translation: []rune("[:upper:]"),
			output:      &buf,
		}

		cfg.translateCmd()

		got := buf.String()
		want := "CODING CHALLENGE"

		assertEqual(t, got, want)
	})

	t.Run("class specifier alpha to digit", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			input:       strings.NewReader("Coding Challenge123%.?"),
			target:      []rune("[:alpha:]"),
			translation: []rune("[:digit:]"),
			output:      &buf,
		}

		// tr output: 299999 2999999991239999_9999.99.?
		// don't understand logic behind that, implemented my own mapping

		cfg.translateCmd()

		got := buf.String()
		want := "293896 270994964123%.?"

		assertEqual(t, got, want)
	})

	t.Run("regular target and class specifier translation", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			input:       strings.NewReader("Coding HELLO Goodbye 123"),
			target:      []rune("od"),
			translation: []rune("[:upper:]"),
			output:      &buf,
		}

		cfg.translateCmd()

		got := buf.String()
		want := "CODing HELLO GOODbye 123"
		assertEqual(t, got, want)
	})

	t.Run("class specifier target and regular translation, multiline", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			input:       strings.NewReader("abcd abc\ndef def\nabc abc\ndef"),
			target:      []rune("[:lower:]"),
			translation: []rune("xyz"),
			output:      &buf,
		}

		cfg.translateCmd()

		got := buf.String()
		want := "xyzz xyz\nzzz zzz\nxyz xyz\nzzz"

		assertEqual(t, got, want)
	})

	t.Run("class specifier target and regular translation, normal letters", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			input:       strings.NewReader("coding HELLO abc Good 123"),
			target:      []rune("[:lower:]"),
			translation: []rune("xyz"),
			output:      &buf,
		}

		cfg.translateCmd()

		got := buf.String()
		want := "xyzzzz HELLO zzx Gyyz 123"

		assertEqual(t, got, want)
	})

	t.Run("delete regular", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			input:       strings.NewReader("hello..."),
			target:      []rune("e."),
			translation: []rune(""),
			deleteFlag:  true,
			output:      &buf,
		}

		cfg.translateCmd()

		got := buf.String()
		want := "hllo"

		assertEqual(t, got, want)
	})

	t.Run("delete class spec", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			input:       strings.NewReader("hello..."),
			target:      []rune("[:punct:]"),
			translation: []rune(""),
			deleteFlag:  true,
			output:      &buf,
		}

		cfg.translateCmd()

		got := buf.String()
		want := "hello"

		assertEqual(t, got, want)
	})
	/*
		t.Run("squeeze", func(t *testing.T) {
			var buf bytes.Buffer
			cfg := config{
				input:       strings.NewReader("hello...oo"),
				target:      []rune("l."),
				squeezeFlag: true,
				output:      &buf,
			}

			cfg.translateCmd()

			got := buf.String()
			want := "helo.oo"

			assertEqual(t, got, want)
		})
	*/
	/*
		t.Run("class specifier target and regular translation", func(t *testing.T) {
			var buf bytes.Buffer
			cfg := config{
				input:       strings.NewReader("coding HELLO abc Goodbye 123"),
				target:      "[:lower:]",
				translation: "💚🥰🐹😊👀",
				output:      &buf,
			}

			cfg.translateCmd()

			got := buf.String()
			want := "💚🥰🐹😊👀👀 HELLO 👀👀💚 G🥰🥰🐹👀👀👀 123"

			assertEqual(t, got, want)
		})
	*/
}

// coding👀HELLO👀abc👀Goodbye👀123🥰
//
// 🐹👀ding 👀👀👀👀👀 abc 👀oodbye🥰 123
// coding 👀👀👀👀👀 abc 👀oodbye 123

func assertEqual(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
