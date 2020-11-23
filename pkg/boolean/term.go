package boolean

// MaxVars is the maximum number of Boolean variables in a Term.
const MaxVars = 6

// Complexity values for the different Boolean operations.
const (
	ConstComplexity  = 0
	VarComplexity    = 1
	UnaryComplexity  = 1
	BinaryComplexity = 1
)

const (
	trueNotation  = "1"
	falseNotation = "0"
	notNotation   = "¬"
	andNotation   = " ∧ "
	orNotation    = " ∨ "
	xorNotation   = " ⊕ "
)

var (
	// Truth defines a Term that is always true.
	Truth = &T{}
	// Falsity defines a Term that is always false.
	Falsity = &F{}

	// Variables defines all Boolean variables.
	Variables = [MaxVars]*Variable{
		newVariable(0), newVariable(1), newVariable(2),
		newVariable(3), newVariable(4), newVariable(5),
	}
)

// Term defines a Boolean term.
type Term interface {
	// Eval evaluates the term for all possible truth values.
	Eval() uint64

	// Complexity returns the complexity of the term.
	Complexity() int

	// String returns a human-readable notation of the term.
	String() string
}

// Eval evaluates the Boolean term t for the provided truth values.
func Eval(t Term, vs ...bool) bool {
	if len(vs) > MaxVars {
		panic("too many truth values")
	}
	// compute the bit index that corresponds to vs
	var idx uint
	for i := range vs {
		if vs[i] {
			idx |= 1 << i
		}
	}
	// evaluate the term and return the corresponding bit
	return (t.Eval()>>idx)&1 != 0
}

// NewNOT creates a new Term from NOT p.
func NewNOT(p Term) *NOT {
	if p == nil {
		panic("nil term")
	}
	return &NOT{p}
}

// NewAND creates a new Term from p AND q.
func NewAND(p, q Term) *AND {
	if p == nil || q == nil {
		panic("nil term")
	}
	return &AND{p, q}
}

// NewOR creates a new Term from p OR q.
func NewOR(p, q Term) *OR {
	if p == nil || q == nil {
		panic("nil term")
	}
	return &OR{p, q}
}

// NewXOR creates a new Term from p XOR q.
func NewXOR(p, q Term) *XOR {
	if p == nil || q == nil {
		panic("nil term")
	}
	return &XOR{p, q}
}

type T struct{}

func (T) Eval() uint64    { return ^uint64(0) }
func (T) String() string  { return trueNotation }
func (T) Complexity() int { return ConstComplexity }

type F struct{}

func (F) Eval() uint64    { return 0 }
func (F) String() string  { return falseNotation }
func (F) Complexity() int { return ConstComplexity }

type Variable struct {
	eval uint64
	name string
}

func (t *Variable) Eval() uint64   { return t.eval }
func (t *Variable) String() string { return t.name }
func (Variable) Complexity() int   { return VarComplexity }

type NOT struct {
	p Term
}

func (t *NOT) Eval() uint64    { return ^t.p.Eval() }
func (t *NOT) String() string  { return notNotation + t.p.String() }
func (t *NOT) Complexity() int { return UnaryComplexity + t.p.Complexity() }

type AND struct {
	p, q Term
}

func (t *AND) Eval() uint64    { return t.p.Eval() & t.q.Eval() }
func (t *AND) String() string  { return "(" + t.p.String() + andNotation + t.q.String() + ")" }
func (t *AND) Complexity() int { return BinaryComplexity + t.p.Complexity() + t.q.Complexity() }

type OR struct {
	p, q Term
}

func (t *OR) Eval() uint64    { return t.p.Eval() | t.q.Eval() }
func (t *OR) String() string  { return "(" + t.p.String() + orNotation + t.q.String() + ")" }
func (t *OR) Complexity() int { return BinaryComplexity + t.p.Complexity() + t.q.Complexity() }

type XOR struct {
	p, q Term
}

func (t *XOR) Eval() uint64    { return t.p.Eval() ^ t.q.Eval() }
func (t *XOR) String() string  { return "(" + t.p.String() + xorNotation + t.q.String() + ")" }
func (t *XOR) Complexity() int { return BinaryComplexity + t.p.Complexity() + t.q.Complexity() }

// each bit index corresponds to a possible combination of truth values
var varBits = [MaxVars]uint64{
	0xaaaaaaaaaaaaaaaa,
	0xcccccccccccccccc,
	0xf0f0f0f0f0f0f0f0,
	0xff00ff00ff00ff00,
	0xffff0000ffff0000,
	0xffffffff00000000,
}

func newVariable(index uint) *Variable {
	if index >= MaxVars {
		panic("invalid index")
	}
	return &Variable{varBits[index], string(rune('a' + index))}
}
