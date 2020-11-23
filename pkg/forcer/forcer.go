package forcer

import (
	"github.com/wollac/boolforcer/pkg/boolean"
)

// Forcer is a brute force boolean expression minimizer.
type Forcer struct {
	mask, eval uint64
}

const initialCapacity = 256

// New creates a new Forcer based on the truth table.
func New(truthTable *boolean.TruthTable) *Forcer {
	mask, eval := truthTable.Evaluation()
	return &Forcer{
		mask: mask,
		eval: eval,
	}
}

// Run starts the brute forcing to find minimal terms up to a complexity of maxComplexity.
func (f *Forcer) Run(maxComplexity int) []boolean.Term {
	sol := &solution{}
	f.buildTerms(maxComplexity, sol)
	return sol.terms
}

func (f *Forcer) add(sol *solution, result *evalTerms, terms ...boolean.Term) {
	for _, t := range terms {
		// apply the mask to discard values for undefined inputs
		eval := t.Eval() & f.mask
		// when this matches the expected evaluation, we found a solution
		if eval == f.eval {
			sol.terms = append(sol.terms, t)
		}
		// do not add equivalent terms
		if _, contains := result.evals[eval]; contains {
			continue
		}
		// add as a cachedTerm
		result.evals[eval] = struct{}{}
		result.terms = append(result.terms, &cachedTerm{t, eval, t.Complexity()})
	}
}

// buildTerms constructs all possible Boolean boolean.Terms of the given complexity.
func (f *Forcer) buildTerms(level int, sol *solution) *evalTerms {
	// the constants have the lowest complexity and represent the base case
	if level == boolean.ConstComplexity {
		result := newResult()
		f.add(sol, result, boolean.Truth, boolean.Falsity)
		return result
	}

	// construct all distinct boolean.Terms with lower complexity
	result := f.buildTerms(level-1, sol)
	if sol.found() {
		return result
	}
	// non constant terms of the lower complexity
	lower := result.terms[2:]

	// add variables
	if level == boolean.VarComplexity {
		for i := range boolean.Variables {
			f.add(sol, result, boolean.Variables[i])
		}
		if sol.found() {
			return result
		}
	}

	// add negations
	if level >= boolean.UnaryComplexity+boolean.VarComplexity {
		for _, p := range lower {
			// a negation must match the desired complexity level
			if p.Complexity()+boolean.UnaryComplexity != level {
				continue
			}
			// do not negate twice
			if _, ok := p.term.(*boolean.NOT); ok {
				continue
			}
			f.add(sol, result, boolean.NewNOT(p.term))
		}
		if sol.found() {
			return result
		}
	}

	// add binary operations
	if level >= boolean.BinaryComplexity+2*boolean.VarComplexity {
		for i, p := range lower {
			// since constants are already excluded the minimal complexity increase is a combination with a variable
			// is this is already larger than the desired complexity level we can skip p
			if p.Complexity()+boolean.VarComplexity+boolean.BinaryComplexity > level {
				continue
			}
			for j, q := range lower {
				// ignore the order as all binary operators are commutative
				if i >= j {
					continue
				}
				// a combination must match the desired complexity level
				if p.Complexity()+q.Complexity()+boolean.BinaryComplexity != level {
					continue
				}
				f.add(sol, result,
					boolean.NewAND(p.term, q.term), boolean.NewOR(p.term, q.term), boolean.NewXOR(p.term, q.term),
				)
			}
		}
	}
	return result
}

// cachedTerm is a boolean.Term with cached method values.
type cachedTerm struct {
	term       boolean.Term
	eval       uint64
	complexity int
}

func (t *cachedTerm) Eval() uint64    { return t.eval }
func (t *cachedTerm) Complexity() int { return t.complexity }
func (t *cachedTerm) String() string  { return t.term.String() }

type evalTerms struct {
	evals map[uint64]struct{}
	terms []*cachedTerm
}

type solution struct {
	terms []boolean.Term
}

func (s *solution) found() bool {
	return len(s.terms) > 0
}

func newResult() *evalTerms {
	return &evalTerms{
		evals: make(map[uint64]struct{}, initialCapacity),
		terms: make([]*cachedTerm, 0, initialCapacity),
	}
}
