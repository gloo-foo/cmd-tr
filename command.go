package command

import (
	"strings"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// lineTransform rewrites one input line into its translated form.
type lineTransform func(string) string

// Tr returns a Command that translates, deletes, or squeezes characters.
// from is the source character set (set1), to is the replacement character set
// (set2). Character ranges like "a-z" are expanded automatically.
//
// Flags:
//   - TrDelete (-d):     delete characters in set1
//   - TrSqueeze (-s):    squeeze repeated characters
//   - TrComplement (-c): operate on the complement of set1
func Tr(from, to string, opts ...any) gloo.Command[[]byte, []byte] {
	f := gloo.NewParameters[gloo.File, flags](opts...).Flags
	transform := compose(plan(expandSet(from), expandSet(to), f))
	return patterns.Map(func(line []byte) ([]byte, error) {
		return []byte(transform(string(line))), nil
	})
}

// plan selects the ordered transforms a given flag set applies. The returned
// stages are run left to right by compose.
func plan(set1, set2 string, f flags) []lineTransform {
	stages := []lineTransform{primary(set1, set2, f)}
	if bool(f.squeeze) {
		stages = append(stages, squeezeStage(set1, set2))
	}
	return stages
}

// primary chooses the delete-or-translate stage that runs before squeezing. A
// run with neither -d nor a set2 to translate into leaves the line unchanged.
func primary(set1, set2 string, f flags) lineTransform {
	if bool(f.del) {
		return deleteWith(member(set1, bool(f.complement)))
	}
	if set2 != "" {
		return translateWith(set1, set2, bool(f.complement))
	}
	return identity
}

// squeezeStage builds the -s stage. The squeezed set is set2 when translating
// into one, otherwise set1, matching GNU tr.
func squeezeStage(set1, set2 string) lineTransform {
	squeezeSet := set2
	if squeezeSet == "" {
		squeezeSet = set1
	}
	return squeezeWith(squeezeSet)
}

// compose folds a list of transforms into one applied left to right.
func compose(stages []lineTransform) lineTransform {
	return func(s string) string {
		for _, stage := range stages {
			s = stage(s)
		}
		return s
	}
}

// identity returns its input unchanged.
func identity(s string) string { return s }

// member reports membership in set, inverted when complement is set. The
// predicate names the characters a delete keeps and a complement translate
// leaves untouched.
func member(set string, complement bool) func(rune) bool {
	return func(r rune) bool {
		return strings.ContainsRune(set, r) != complement
	}
}

// deleteWith drops every rune for which drop reports true, keeping the rest.
func deleteWith(drop func(rune) bool) lineTransform {
	return func(s string) string {
		return filterRunes(s, func(r rune) bool { return !drop(r) })
	}
}

// translateWith maps each rune through a translation table. With complement,
// every rune outside set1 becomes the last rune of set2.
func translateWith(set1, set2 string, complement bool) lineTransform {
	mapRune := mapper(set1, set2, complement)
	return func(s string) string {
		return mapRunes(s, mapRune)
	}
}

// squeezeWith collapses consecutive runs of any rune in squeezeSet to one.
func squeezeWith(squeezeSet string) lineTransform {
	keep := keepAfterRun(squeezeSet)
	return func(s string) string {
		var out []rune
		var prev rune
		for i, r := range s {
			if keep(r, prev, i == 0) {
				out = append(out, r)
			}
			prev = r
		}
		return string(out)
	}
}

// keepAfterRun reports whether a rune survives squeezing given the previous
// rune and whether it is the first rune of the line.
func keepAfterRun(squeezeSet string) func(r, prev rune, first bool) bool {
	return func(r, prev rune, first bool) bool {
		return first || r != prev || !strings.ContainsRune(squeezeSet, r)
	}
}

// mapper builds the per-rune translation function. Without complement it maps
// set1[i]->set2[i] (clamping to the last set2 rune); with complement every
// rune outside set1 maps to the last set2 rune.
func mapper(set1, set2 string, complement bool) func(rune) rune {
	runes2 := []rune(set2)
	last := runes2[len(runes2)-1]
	if complement {
		return complementMapper(member(set1, false), last)
	}
	return tableMapper(translationTable(set1, runes2, last))
}

// complementMapper maps runes outside set1 (per inSet1) to last, others to self.
func complementMapper(inSet1 func(rune) bool, last rune) func(rune) rune {
	return func(r rune) rune {
		if inSet1(r) {
			return r
		}
		return last
	}
}

// tableMapper maps a rune through table, leaving absent runes unchanged.
func tableMapper(table map[rune]rune) func(rune) rune {
	return func(r rune) rune {
		if to, ok := table[r]; ok {
			return to
		}
		return r
	}
}

// translationTable maps each set1 rune to its set2 counterpart, clamping to
// last once set2 is exhausted. The first mapping for a rune wins, as in GNU tr.
func translationTable(set1 string, runes2 []rune, last rune) map[rune]rune {
	table := make(map[rune]rune)
	for i, r := range []rune(set1) {
		if _, seen := table[r]; !seen {
			table[r] = clampAt(runes2, i, last)
		}
	}
	return table
}

// clampAt returns runes2[i], or last when i is past the end of runes2.
func clampAt(runes2 []rune, i int, last rune) rune {
	if i < len(runes2) {
		return runes2[i]
	}
	return last
}

// filterRunes keeps only the runes for which keep reports true.
func filterRunes(s string, keep func(rune) bool) string {
	var out []rune
	for _, r := range s {
		if keep(r) {
			out = append(out, r)
		}
	}
	return string(out)
}

// mapRunes rewrites every rune of s through f.
func mapRunes(s string, f func(rune) rune) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		out = append(out, f(r))
	}
	return string(out)
}

// expandSet expands character ranges like "a-z" into their full sequence of
// runes ("abc...xyz"), leaving other runes untouched.
func expandSet(set string) string {
	runes := []rune(set)
	var out []rune
	for i := 0; i < len(runes); i++ {
		if isRange(runes, i) {
			out = appendRange(out, runes[i], runes[i+2])
			i += 2
			continue
		}
		out = append(out, runes[i])
	}
	return string(out)
}

// isRange reports whether position i begins an "a-z" range triple.
func isRange(runes []rune, i int) bool {
	return i+2 < len(runes) && runes[i+1] == '-'
}

// appendRange appends every rune from start through end inclusive.
func appendRange(out []rune, start, end rune) []rune {
	for ch := start; ch <= end; ch++ {
		out = append(out, ch)
	}
	return out
}
