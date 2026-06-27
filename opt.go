package command

// trDeleteFlag controls character deletion mode (-d).
type trDeleteFlag bool

const (
	TrDelete   trDeleteFlag = true
	TrNoDelete trDeleteFlag = false
)

func (f trDeleteFlag) Configure(flags *flags) { flags.del = f }

// trSqueezeFlag controls character squeeze mode (-s).
type trSqueezeFlag bool

const (
	TrSqueeze   trSqueezeFlag = true
	TrNoSqueeze trSqueezeFlag = false
)

func (f trSqueezeFlag) Configure(flags *flags) { flags.squeeze = f }

// trComplementFlag controls set complement mode (-c).
type trComplementFlag bool

const (
	TrComplement   trComplementFlag = true
	TrNoComplement trComplementFlag = false
)

func (f trComplementFlag) Configure(flags *flags) { flags.complement = f }

// flags is the parsed flag set for a tr run.
type flags struct {
	del        trDeleteFlag
	squeeze    trSqueezeFlag
	complement trComplementFlag
}
