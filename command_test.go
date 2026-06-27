package command_test

import (
	"fmt"
	"testing"

	command "github.com/gloo-foo/cmd-tr"
	"github.com/gloo-foo/testable"
)

func TestTr_BasicTranslation(t *testing.T) {
	lines, err := testable.TestLines(
		command.Tr("abc", "xyz"),
		"abc\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != "xyz" {
		t.Fatalf("got %q, want [xyz]", lines)
	}
}

func TestTr_RangeExpansion(t *testing.T) {
	lines, err := testable.TestLines(
		command.Tr("a-z", "A-Z"),
		"hello\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != "HELLO" {
		t.Fatalf("got %q, want [HELLO]", lines)
	}
}

func TestTr_Delete(t *testing.T) {
	lines, err := testable.TestLines(
		command.Tr("aeiou", "", command.TrDelete),
		"hello world\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != "hll wrld" {
		t.Fatalf("got %q, want [hll wrld]", lines)
	}
}

func TestTr_Squeeze(t *testing.T) {
	lines, err := testable.TestLines(
		command.Tr(" ", "", command.TrSqueeze),
		"hello   world\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != "hello world" {
		t.Fatalf("got %q, want [hello world]", lines)
	}
}

func TestTr_TranslateAndSqueeze(t *testing.T) {
	lines, err := testable.TestLines(
		command.Tr("a", "x", command.TrSqueeze),
		"aaa bbb aaa\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != "x bbb x" {
		t.Fatalf("got %q, want [x bbb x]", lines)
	}
}

func TestTr_Set1LongerThanSet2(t *testing.T) {
	lines, err := testable.TestLines(
		command.Tr("abc", "x"),
		"abcdef\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != "xxxdef" {
		t.Fatalf("got %q, want [xxxdef]", lines)
	}
}

func TestTr_MultipleLines(t *testing.T) {
	lines, err := testable.TestLines(
		command.Tr("aeiou", "AEIOU"),
		"hello\nworld\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2", len(lines))
	}
	if lines[0] != "hEllO" {
		t.Errorf("line 0: got %q, want %q", lines[0], "hEllO")
	}
	if lines[1] != "wOrld" {
		t.Errorf("line 1: got %q, want %q", lines[1], "wOrld")
	}
}

func TestTr_NoMatchingChars(t *testing.T) {
	lines, err := testable.TestLines(
		command.Tr("xyz", "ABC"),
		"hello\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != "hello" {
		t.Fatalf("got %q, want [hello]", lines)
	}
}

func TestTr_EmptyInput(t *testing.T) {
	lines, err := testable.TestLines(
		command.Tr("a", "b"),
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 0 {
		t.Fatalf("got %q, want empty", lines)
	}
}

func TestTr_DeleteRange(t *testing.T) {
	lines, err := testable.TestLines(
		command.Tr("0-9", "", command.TrDelete),
		"abc123def456\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != "abcdef" {
		t.Fatalf("got %q, want [abcdef]", lines)
	}
}

func TestTr_SqueezeAfterTranslate(t *testing.T) {
	lines, err := testable.TestLines(
		command.Tr("ae", "xx", command.TrSqueeze),
		"aabbee\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	// a->x, e->x, then squeeze x: "xxbbxx" -> "xbbx"
	if len(lines) != 1 || lines[0] != "xbbx" {
		t.Fatalf("got %q, want [xbbx]", lines)
	}
}

// ==============================================================================
// Complement (-c)
// ==============================================================================

func TestTr_Complement_Delete(t *testing.T) {
	// Delete chars NOT in set1 (keep only lowercase letters)
	lines, err := testable.TestLines(
		command.Tr("a-z", "", command.TrDelete, command.TrComplement),
		"Hello World 123\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != "elloorld" {
		t.Fatalf("got %q, want [elloorld]", lines)
	}
}

func TestTr_Complement_Translate(t *testing.T) {
	// Translate chars NOT in set1 to last char of set2
	lines, err := testable.TestLines(
		command.Tr("a-z", "_", command.TrComplement),
		"Hello World 123\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != "_ello__orld____" {
		t.Fatalf("got %q, want [_ello__orld____]", lines)
	}
}

func TestTr_Complement_Delete_Digits(t *testing.T) {
	// Keep only digits
	lines, err := testable.TestLines(
		command.Tr("0-9", "", command.TrDelete, command.TrComplement),
		"abc123def456\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != "123456" {
		t.Fatalf("got %q, want [123456]", lines)
	}
}

// ==============================================================================
// Examples
// ==============================================================================

func ExampleTr() {
	lines, _ := testable.TestLines(
		command.Tr("a-z", "A-Z"),
		"hello world\n",
	)
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// HELLO WORLD
}

func ExampleTr_delete() {
	lines, _ := testable.TestLines(
		command.Tr("aeiou", "", command.TrDelete),
		"hello world\n",
	)
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// hll wrld
}

func ExampleTr_squeeze() {
	lines, _ := testable.TestLines(
		command.Tr(" ", "", command.TrSqueeze),
		"hello   world\n",
	)
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// hello world
}
