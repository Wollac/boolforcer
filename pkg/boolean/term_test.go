package boolean

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var tests = []struct {
	term       Term
	complexity int
	eval       func(...bool) bool
}{
	{
		Truth,
		0,
		func(...bool) bool { return true },
	},
	{
		Falsity,
		0,
		func(...bool) bool { return false },
	},
	{
		Variables[0],
		1,
		func(v ...bool) bool { return v[0] },
	},
	{
		Variables[1],
		1,
		func(v ...bool) bool { return v[1] },
	},
	{
		Variables[2],
		1,
		func(v ...bool) bool { return v[2] },
	},
	{
		Variables[3],
		1,
		func(v ...bool) bool { return v[3] },
	},
	{
		Variables[4],
		1,
		func(v ...bool) bool { return v[4] },
	},
	{
		Variables[5],
		1,
		func(v ...bool) bool { return v[5] },
	},
	{
		&NOT{Variables[0]},
		2,
		func(v ...bool) bool { return !v[0] },
	},
	{
		&AND{Variables[0], Variables[1]},
		3,
		func(v ...bool) bool { return v[0] && v[1] },
	},
	{
		&OR{Variables[0], Variables[1]},
		3,
		func(v ...bool) bool { return v[0] || v[1] },
	},
	{
		&XOR{Variables[0], Variables[1]},
		3,
		func(v ...bool) bool { return v[0] != v[1] },
	},
	{
		&XOR{Variables[0], Truth},
		2,
		func(v ...bool) bool { return v[0] != true },
	},
	{
		&NOT{&OR{Variables[0], Variables[1]}},
		4,
		func(v ...bool) bool { return !(v[0] || v[1]) },
	},
	{
		&AND{
			&NOT{Variables[0]},
			&XOR{
				&NOT{Variables[1]},
				&OR{
					&AND{
						&NOT{Variables[2]},
						&NOT{Variables[3]},
					},
					&XOR{
						&NOT{Variables[4]},
						&NOT{Variables[5]},
					},
				},
			},
		},
		17,
		func(v ...bool) bool { return !v[0] && (!v[1] != ((!v[2] && !v[3]) || (!v[4] != !v[5]))) },
	},
}

func TestEval(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.term.String(), func(t *testing.T) {
			for i := 0; i < 1<<MaxVars; i++ {
				v := []bool{i&0x1 != 0, i&0x2 != 0, i&0x4 != 0, i&0x8 != 0, i&0x10 != 0, i&0x20 != 0}
				require.Equalf(t,
					tt.eval(v...), Eval(tt.term, v...),
					"%v with a=%b, b=%b, c=%b, d=%b, e=%b, f=%b", tt.term, v[0], v[1], v[2], v[3], v[4], v[5],
				)
			}
		})
	}
}

func TestComplexity(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.term.String(), func(t *testing.T) {
			require.Equal(t, tt.complexity, tt.term.Complexity())
		})
	}
}

func BenchmarkEval(b *testing.B) {
	term := tests[14].term
	b.Log(term)

	for i := 0; i < b.N; i++ {
		_ = term.Eval()
	}
}
