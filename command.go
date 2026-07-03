package command

import (
	"strings"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// lineText is one line of input text, without its trailing newline.
type lineText string

// charSet is a fully expanded character set: every "a-z" range of a TrSet has
// been unrolled into its individual runes.
type charSet string

// clampRune is the last rune of set2, to which translation clamps: set1 runes
// past the end of set2 and, under -c, every rune outside set1 map to it.
type clampRune rune

// setIndex is the position of a rune within a character set; a set1 rune at
// index i translates to the set2 rune at the same index.
type setIndex int

// rangeBound is an endpoint rune of an "a-z" style character range.
type rangeBound rune

// lineTransform rewrites one input line into its translated form.
type lineTransform func(lineText) lineText

// Tr returns a Command that translates, deletes, or squeezes characters.
// from is the source character set (set1), to is the replacement character set
// (set2). Character ranges like "a-z" are expanded automatically.
//
// Flags:
//   - TrDelete (-d):     delete characters in set1
//   - TrSqueeze (-s):    squeeze repeated characters
//   - TrComplement (-c): operate on the complement of set1
func Tr(from, to TrSet, opts ...any) gloo.Command[[]byte, []byte] {
	transform := compose(plan(expandSet(from), expandSet(to), options(opts...)))
	return patterns.Map(func(b []byte) ([]byte, error) {
		return []byte(transform(lineText(b))), nil
	})
}

// plan selects the ordered transforms a given flag set applies. The returned
// stages are run left to right by compose.
func plan(set1, set2 charSet, f flags) []lineTransform {
	stages := []lineTransform{primary(set1, set2, f)}
	if bool(f.isSqueeze) {
		stages = append(stages, squeezeStage(set1, set2))
	}
	return stages
}

// primary chooses the delete-or-translate stage that runs before squeezing. A
// run with neither -d nor a set2 to translate into leaves the line unchanged.
func primary(set1, set2 charSet, f flags) lineTransform {
	if bool(f.isDelete) {
		return deleteWith(member(set1, f.isComplement))
	}
	if set2 != "" {
		return translateWith(set1, set2, f.isComplement)
	}
	return identity
}

// squeezeStage builds the -s stage. The squeezed set is set2 when translating
// into one, otherwise set1, matching GNU tr.
func squeezeStage(set1, set2 charSet) lineTransform {
	squeezeSet := set2
	if squeezeSet == "" {
		squeezeSet = set1
	}
	return squeezeWith(squeezeSet)
}

// compose folds a list of transforms into one applied left to right.
func compose(stages []lineTransform) lineTransform {
	return func(s lineText) lineText {
		for _, stage := range stages {
			s = stage(s)
		}
		return s
	}
}

// identity returns its input unchanged.
func identity(s lineText) lineText { return s }

// member reports membership in set, inverted under -c. The predicate names the
// characters a delete keeps and a complement translate leaves untouched.
func member(set charSet, isComplement trComplementFlag) func(rune) bool {
	return func(r rune) bool {
		return strings.ContainsRune(string(set), r) != bool(isComplement)
	}
}

// deleteWith drops every rune for which drop reports true, keeping the rest.
func deleteWith(drop func(rune) bool) lineTransform {
	return func(s lineText) lineText {
		return filterRunes(s, func(r rune) bool { return !drop(r) })
	}
}

// translateWith maps each rune through a translation table. Under -c, every
// rune outside set1 becomes the last rune of set2.
func translateWith(set1, set2 charSet, isComplement trComplementFlag) lineTransform {
	mapRune := mapper(set1, set2, isComplement)
	return func(s lineText) lineText {
		return mapRunes(s, mapRune)
	}
}

// squeezeWith collapses consecutive runs of any rune in squeezeSet to one.
func squeezeWith(squeezeSet charSet) lineTransform {
	keep := keepAfterRun(squeezeSet)
	return func(s lineText) lineText {
		var out []rune
		var prev rune
		for i, r := range s {
			if keep(r, prev, i == 0) {
				out = append(out, r)
			}
			prev = r
		}
		return lineText(out)
	}
}

// keepAfterRun reports whether a rune survives squeezing given the previous
// rune and whether it is the first rune of the line.
func keepAfterRun(squeezeSet charSet) func(r, prev rune, isFirst bool) bool {
	return func(r, prev rune, isFirst bool) bool {
		return isFirst || r != prev || !strings.ContainsRune(string(squeezeSet), r)
	}
}

// mapper builds the per-rune translation function. Without -c it maps
// set1[i]->set2[i] (clamping to the last set2 rune); under -c every rune
// outside set1 maps to the last set2 rune.
func mapper(set1, set2 charSet, isComplement trComplementFlag) func(rune) rune {
	runes2 := []rune(string(set2))
	last := clampRune(runes2[len(runes2)-1])
	if bool(isComplement) {
		return complementMapper(member(set1, TrNoComplement), last)
	}
	return tableMapper(translationTable(set1, runes2, last))
}

// complementMapper maps runes outside set1 (per inSet1) to last, others to self.
func complementMapper(inSet1 func(rune) bool, last clampRune) func(rune) rune {
	return func(r rune) rune {
		if inSet1(r) {
			return r
		}
		return rune(last)
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
func translationTable(set1 charSet, runes2 []rune, last clampRune) map[rune]rune {
	table := make(map[rune]rune)
	for i, r := range []rune(string(set1)) {
		if _, seen := table[r]; !seen {
			table[r] = clampAt(runes2, setIndex(i), last)
		}
	}
	return table
}

// clampAt returns runes2[i], or last when i is past the end of runes2.
func clampAt(runes2 []rune, i setIndex, last clampRune) rune {
	if int(i) < len(runes2) {
		return runes2[int(i)]
	}
	return rune(last)
}

// filterRunes keeps only the runes for which keep reports true.
func filterRunes(s lineText, keep func(rune) bool) lineText {
	var out []rune
	for _, r := range s {
		if keep(r) {
			out = append(out, r)
		}
	}
	return lineText(out)
}

// mapRunes rewrites every rune of s through f.
func mapRunes(s lineText, f func(rune) rune) lineText {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		out = append(out, f(r))
	}
	return lineText(out)
}

// expandSet expands character ranges like "a-z" into their full sequence of
// runes ("abc...xyz"), leaving other runes untouched.
func expandSet(set TrSet) charSet {
	runes := []rune(string(set))
	var out []rune
	for i := 0; i < len(runes); i++ {
		if isRange(runes, setIndex(i)) {
			out = appendRange(out, rangeBound(runes[i]), rangeBound(runes[i+2]))
			i += 2
			continue
		}
		out = append(out, runes[i])
	}
	return charSet(out)
}

// isRange reports whether position i begins an "a-z" range triple.
func isRange(runes []rune, i setIndex) bool {
	return int(i)+2 < len(runes) && runes[int(i)+1] == '-'
}

// appendRange appends every rune from start through end inclusive.
func appendRange(out []rune, start, end rangeBound) []rune {
	for ch := rune(start); ch <= rune(end); ch++ {
		out = append(out, ch)
	}
	return out
}
