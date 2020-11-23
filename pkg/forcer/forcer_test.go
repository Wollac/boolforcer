package forcer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wollac/boolforcer/pkg/boolean"
)

func TestForcer_Run(t *testing.T) {
	var tests = []struct {
		table *boolean.TruthTable
		term  boolean.Term
	}{
		{
			&boolean.TruthTable{boolean.True},
			boolean.Truth,
		},
		{
			&boolean.TruthTable{boolean.False},
			boolean.Falsity,
		},
		{
			&boolean.TruthTable{boolean.False, boolean.True},
			boolean.Variables[0],
		},
		{
			&boolean.TruthTable{boolean.True, boolean.False},
			boolean.NewNOT(boolean.Variables[0]),
		},
		{
			&boolean.TruthTable{boolean.False, boolean.False, boolean.False, boolean.True},
			boolean.NewAND(boolean.Variables[0], boolean.Variables[1]),
		},
		{
			&boolean.TruthTable{boolean.False, boolean.True, boolean.True, boolean.True},
			boolean.NewOR(boolean.Variables[0], boolean.Variables[1]),
		},
		{
			&boolean.TruthTable{boolean.False, boolean.True, boolean.True, boolean.False},
			boolean.NewXOR(boolean.Variables[0], boolean.Variables[1]),
		},
		{
			&boolean.TruthTable{boolean.True, boolean.False, boolean.True, boolean.True},
			boolean.NewOR(boolean.NewNOT(boolean.Variables[0]), boolean.Variables[1]),
		},
		{
			&boolean.TruthTable{boolean.True, boolean.False, boolean.False, boolean.True},
			boolean.NewNOT(boolean.NewXOR(boolean.Variables[0], boolean.Variables[1])),
		},
		{
			&boolean.TruthTable{boolean.False, boolean.True, boolean.True, boolean.True, boolean.True, boolean.True, boolean.True, boolean.True},
			boolean.NewOR(boolean.Variables[0], boolean.NewOR(boolean.Variables[1], boolean.Variables[2])),
		},
		{
			&boolean.TruthTable{boolean.True, boolean.False, boolean.False, boolean.False, boolean.False, boolean.False, boolean.False, boolean.False},
			boolean.NewNOT(boolean.NewOR(boolean.Variables[0], boolean.NewOR(boolean.Variables[1], boolean.Variables[2]))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.table.String(), func(t *testing.T) {
			terms := New(tt.table).Run(32)
			for _, term := range terms {
				assert.Equalf(t, tt.term.Eval(), term.Eval(), "term %v is not equivalent to expected %v", term, tt.term)
				assert.Equalf(t, tt.term.Complexity(), term.Complexity(), "complexity of term %v:%d is not equal to expected %d", term, term.Complexity(), tt.term.Complexity())
			}
		})
	}
}
