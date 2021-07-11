package fastutf8

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

// Used for generating table for shift-based DFA based on algorithm described in
// https://gist.github.com/pervognsen/218ea17743e1442e59bb60d29b1aa725
func TestPrintStates(t *testing.T) {
	t.Skip()
	var a [256]uint64
	for b := 0; b < 256; b++ {
		state := uint64(0)
		for s := 8; s >= 0; s-- {
			typ := int(utf8d[b])
			state <<= 6
			state |= uint64(utf8d[256+s*16+typ] * 6)
		}
		a[b] = state
	}
	fmt.Printf("%#v\n", a)
}

// Benchmarks and tests are taken from stdlib.
func BenchmarkValidTenASCIIChars(b *testing.B) {
	s := []byte("0123456789")
	b.Run("stdlib", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utf8.Valid(s)
		}
	})
	b.Run("shift", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Valid(s)
		}
	})
}

func BenchmarkValidTenJapaneseChars(b *testing.B) {
	s := []byte("日本語日本語日本語日")
	b.Run("stdlib", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utf8.Valid(s)
		}
	})
	b.Run("shift", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Valid(s)
		}
	})
}

func BenchmarkRuneCountTenASCIIChars(b *testing.B) {
	s := []byte("0123456789")
	b.Run("stdlib", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utf8.RuneCount(s)
		}
	})
	b.Run("shift", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RuneCount(s)
		}
	})
}

func BenchmarkRuneCountTenJapaneseChars(b *testing.B) {
	s := []byte("日本語日本語日本語日")
	b.Run("stdlib", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utf8.RuneCount(s)
		}
	})
	b.Run("shift", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RuneCount(s)
		}
	})
}

// All tests are taken from https://github.com/golang/go/blob/master/src/unicode/utf8/utf8_test.go.
type RuneCountTest struct {
	in  string
	out int
}

var runecounttests = []RuneCountTest{
	{"abcd", 4},
	{"☺☻☹", 3},
	{"1,2,3,4", 7},
	// Skip tests with invalid runes.
	//{"\xe2\x00", 2},
	//{"\xe2\x80", 2},
	//{"a\xe2\x80", 3},
}

func TestRuneCount(t *testing.T) {
	for _, tt := range runecounttests {
		if out := RuneCount([]byte(tt.in)); out != tt.out {
			t.Errorf("RuneCount(%q) = %d, want %d", tt.in, out, tt.out)
		}
	}
}

type ValidTest struct {
	in  string
	out bool
}

var validTests = []ValidTest{
	{"", true},
	{"a", true},
	{"abc", true},
	{"Ж", true},
	{"ЖЖ", true},
	{"брэд-ЛГТМ", true},
	{"☺☻☹", true},
	{"aa\xe2", false},
	{string([]byte{66, 250}), false},
	{string([]byte{66, 250, 67}), false},
	{"a\uFFFDb", true},
	{string("\xF4\x8F\xBF\xBF"), true},      // U+10FFFF
	{string("\xF4\x90\x80\x80"), false},     // U+10FFFF+1; out of range
	{string("\xF7\xBF\xBF\xBF"), false},     // 0x1FFFFF; out of range
	{string("\xFB\xBF\xBF\xBF\xBF"), false}, // 0x3FFFFFF; out of range
	{string("\xc0\x80"), false},             // U+0000 encoded in two bytes: incorrect
	{string("\xed\xa0\x80"), false},         // U+D800 high surrogate (sic)
	{string("\xed\xbf\xbf"), false},         // U+DFFF low surrogate (sic)
}

func TestValid(t *testing.T) {
	for _, tt := range validTests {
		if Valid([]byte(tt.in)) != tt.out {
			t.Errorf("Valid(%q) = %v; want %v", tt.in, !tt.out, tt.out)
		}
	}
}
