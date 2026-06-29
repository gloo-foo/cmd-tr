package alias_test

import (
	"slices"
	"testing"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/testable"

	tr "github.com/gloo-foo/cmd-tr/alias"
)

// The alias package re-exports the constructor and flag constants under
// unprefixed names. A mis-wired re-export (say, Delete bound to the disabled
// constant, or Tr bound to the wrong function) compiles cleanly, so only
// behavior can prove the wiring. Each test exercises one re-export and asserts
// the GNU tr output it must produce.

func assertLines(t *testing.T, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func runLines(t *testing.T, cmd gloo.Command[[]byte, []byte], input string) []string {
	t.Helper()
	lines, err := testable.TestLines(cmd, input)
	if err != nil {
		t.Fatal(err)
	}
	return lines
}

func TestAlias_TrTranslates(t *testing.T) {
	// Tr must be the translating constructor: a-z -> A-Z uppercases.
	got := runLines(t, tr.Tr("a-z", "A-Z"), "hello\n")
	assertLines(t, got, []string{"HELLO"})
}

func TestAlias_DeleteRemovesSet1(t *testing.T) {
	// Delete is -d: the vowels in set1 are dropped.
	got := runLines(t, tr.Tr("aeiou", "", tr.Delete), "hello world\n")
	assertLines(t, got, []string{"hll wrld"})
}

func TestAlias_SqueezeCollapsesRuns(t *testing.T) {
	// Squeeze is -s: runs of spaces collapse to one.
	got := runLines(t, tr.Tr(" ", "", tr.Squeeze), "hello   world\n")
	assertLines(t, got, []string{"hello world"})
}

func TestAlias_ComplementSelectsOutsideSet1(t *testing.T) {
	// Complement is -c paired with -d: keep only the lowercase letters.
	got := runLines(t, tr.Tr("a-z", "", tr.Delete, tr.Complement), "Hello World 123\n")
	assertLines(t, got, []string{"elloorld"})
}

func TestAlias_DisabledFlagsMatchDefault(t *testing.T) {
	// The No* constants are the disabled forms: they must behave exactly like
	// passing no flag at all (plain translation).
	got := runLines(t, tr.Tr("a-z", "A-Z", tr.NoDelete, tr.NoSqueeze, tr.NoComplement), "hello\n")
	assertLines(t, got, []string{"HELLO"})
}
