package main

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"
)

// === Test helpers ===

func assertError(t testing.TB, got error, wantErr bool) {
	t.Helper()
	if wantErr && got == nil {
		t.Fatal("expected error, got nil")
	}
	if !wantErr && got != nil {
		t.Fatalf("unexpected error: %v", got)
	}
}

func assertString(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertBool(t testing.TB, name string, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

func assertEqual(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// === loadConfig tests ===
// NOTE: loadConfig calls flag.Parse() which reads global os.Args.
// We work around that by setting os.Args and resetting flag.CommandLine.
// In the future we might refactor loadConfig to accept a []string directly.

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		wantDelete      bool
		wantSqueeze     bool
		wantTarget      string
		wantTranslation string
		wantErr         bool
	}{
		{
			name:            "two args: target and translation",
			args:            []string{"cmd", "lo", "bo"},
			wantTarget:      "lo",
			wantTranslation: "bo",
		},
		{
			name:       "delete flag with one arg",
			args:       []string{"cmd", "-d", "abc"},
			wantDelete: true,
			wantTarget: "abc",
		},
		{
			name:        "squeeze flag with one arg",
			args:        []string{"cmd", "-s", "abc"},
			wantSqueeze: true,
			wantTarget:  "abc",
		},
		{
			name:            "delete flag with two args",
			args:            []string{"cmd", "-d", "abc", "xyz"},
			wantDelete:      true,
			wantTarget:      "abc",
			wantTranslation: "xyz",
		},
		{
			name:    "no args",
			args:    []string{"cmd"},
			wantErr: true,
		},
		{
			name:    "one arg without delete flag",
			args:    []string{"cmd", "abc"},
			wantErr: true,
		},
		{
			name:    "three args",
			args:    []string{"cmd", "a", "b", "c"},
			wantErr: true,
		},
		{
			name:        "delete and squeeze together",
			args:        []string{"cmd", "-d", "-s", "abc"},
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args = tc.args

			// Reset global flag state to avoid "flag redefined" panics
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			cfg, err := loadConfig()
			assertError(t, err, tc.wantErr)
			if tc.wantErr {
				return
			}

			assertString(t, cfg.target, tc.wantTarget)
			assertString(t, cfg.translation, tc.wantTranslation)
			assertBool(t, "deleteFlag", cfg.deleteFlag, tc.wantDelete)
			assertBool(t, "squeezeFlag", cfg.squeezeFlag, tc.wantSqueeze)
		})
	}
}

// === buildTranslator tests ===
// These will fail initially because buildTranslator is currently a stub.
// They serve as the TDD guide for the refactor.

func TestBuildTranslator(t *testing.T) {
	t.Run("regular to regular", func(t *testing.T) {
		cfg := config{
			target:      "lo",
			translation: "bo",
		}
		tr, err := buildTranslator(cfg)
		assertError(t, err, false)
		got := strings.Map(tr, "hello")
		assertString(t, got, "hebbo")
	})

	t.Run("regular to regular with delete", func(t *testing.T) {
		cfg := config{
			target:      "lo",
			translation: "bo",
			deleteFlag:  true,
		}
		tr, err := buildTranslator(cfg)
		assertError(t, err, false)
		// Delete means target chars are removed. Using strings.Map,
		// returning -1 drops the character.
		got := strings.Map(tr, "hello")
		assertString(t, got, "he")
	})

	t.Run("range expansion", func(t *testing.T) {
		cfg := config{
			target:      "a-d",
			translation: "e-h",
		}
		tr, err := buildTranslator(cfg)
		assertError(t, err, false)
		got := strings.Map(tr, "abcd")
		assertString(t, got, "efgh")
	})

	t.Run("class specifier lower to upper", func(t *testing.T) {
		cfg := config{
			target:      "[:lower:]",
			translation: "[:upper:]",
		}
		tr, err := buildTranslator(cfg)
		assertError(t, err, false)
		got := strings.Map(tr, "hello")
		assertString(t, got, "HELLO")
	})

	t.Run("class specifier alpha to digit", func(t *testing.T) {
		cfg := config{
			target:      "[:alpha:]",
			translation: "[:digit:]",
		}
		tr, err := buildTranslator(cfg)
		assertError(t, err, false)
		got := strings.Map(tr, "abc")
		assertString(t, got, "999")
	})

	t.Run("invalid class specifier", func(t *testing.T) {
		cfg := config{
			target:      "[:foo:]",
			translation: "[:upper:]",
		}
		_, err := buildTranslator(cfg)
		assertError(t, err, true)
	})

	t.Run("regular target and class specifier translation", func(t *testing.T) {
		cfg := config{
			target:      "od",
			translation: "[:upper:]",
		}
		tr, err := buildTranslator(cfg)
		assertError(t, err, false)
		got := strings.Map(tr, "coding")
		assertString(t, got, "CODing")
	})

	t.Run("class specifier target and regular translation", func(t *testing.T) {
		cfg := config{
			target:      "[:lower:]",
			translation: "xyz",
		}
		tr, err := buildTranslator(cfg)
		assertError(t, err, false)
		got := strings.Map(tr, "abc")
		assertString(t, got, "xyz")
	})

	t.Run("class specifier target and regular translation with extra chars", func(t *testing.T) {
		cfg := config{
			target:      "[:lower:]",
			translation: "xyz",
		}
		tr, err := buildTranslator(cfg)
		assertError(t, err, false)
		// "abc" -> "xyz", then "d" -> "z" (last char repeats)
		got := strings.Map(tr, "abcd")
		assertString(t, got, "xyzz")
	})

	t.Run("emoji rune subst", func(t *testing.T) {
		cfg := config{
			target:      "😊",
			translation: "👀",
		}
		tr, err := buildTranslator(cfg)
		assertError(t, err, false)
		got := strings.Map(tr, "hello😊")
		assertString(t, got, "hello👀")
	})
}

// === run tests ===
// These test the I/O loop. Will fail initially because run is a stub.

func TestRun(t *testing.T) {
	t.Run("single line", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := &config{
			input:  strings.NewReader("hello"),
			output: &buf,
		}
		tr := func(r rune) rune { return r }
		err := cfg.run(tr)
		assertError(t, err, false)
		assertString(t, buf.String(), "hello\n")
	})

	t.Run("multiple lines", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := &config{
			input:  strings.NewReader("hello\nworld"),
			output: &buf,
		}
		tr := func(r rune) rune { return r }
		err := cfg.run(tr)
		assertError(t, err, false)
		assertString(t, buf.String(), "hello\nworld\n")
	})

	t.Run("empty input", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := &config{
			input:  strings.NewReader(""),
			output: &buf,
		}
		tr := func(r rune) rune { return r }
		err := cfg.run(tr)
		assertError(t, err, false)
		assertString(t, buf.String(), "")
	})

	t.Run("translate during run", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := &config{
			input:  strings.NewReader("hello"),
			output: &buf,
		}
		tr := func(r rune) rune {
			if r == 'l' {
				return 'b'
			}
			return r
		}
		err := cfg.run(tr)
		assertError(t, err, false)
		assertString(t, buf.String(), "hebbo\n")
	})
}
