package constants

type Precedence int

const (
	HighHighPrec Precedence = 4
	LowHighPrec  Precedence = 3
	HighLowPrec  Precedence = 2
	LowLowPrec   Precedence = 1
	ParanPrec    Precedence = 0
	None         Precedence = -1
)

var (
	Precedences = map[string]Precedence{
		"*": HighHighPrec,
		"/": LowHighPrec,
		"+": HighLowPrec,
		"-": LowLowPrec,
		"(": ParanPrec,
		")": ParanPrec,
		"":  None,
	}
)
