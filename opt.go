package command

// TrSet is a tr character set argument (SET1 or SET2). Character ranges like
// "a-z" are expanded by the command before use.
type TrSet string

// trDeleteFlag controls character deletion mode (-d).
type trDeleteFlag bool

const (
	TrDelete   trDeleteFlag = true
	TrNoDelete trDeleteFlag = false
)

// trSqueezeFlag controls character squeeze mode (-s).
type trSqueezeFlag bool

const (
	TrSqueeze   trSqueezeFlag = true
	TrNoSqueeze trSqueezeFlag = false
)

// trComplementFlag controls set complement mode (-c).
type trComplementFlag bool

const (
	TrComplement   trComplementFlag = true
	TrNoComplement trComplementFlag = false
)

// flags is the parsed flag set for a tr run.
type flags struct {
	isDelete     trDeleteFlag
	isSqueeze    trSqueezeFlag
	isComplement trComplementFlag
}

// options folds the recognized tr option values out of opts into a flags
// value; later options override earlier ones and anything else is ignored.
func options(opts ...any) flags {
	var f flags
	for _, o := range opts {
		f = f.with(o)
	}
	return f
}

// with returns f updated by a single option value.
func (f flags) with(o any) flags {
	switch v := o.(type) {
	case trDeleteFlag:
		f.isDelete = v
	case trSqueezeFlag:
		f.isSqueeze = v
	case trComplementFlag:
		f.isComplement = v
	}
	return f
}
