package boolean

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTruthTable_Evaluation(t *testing.T) {
	var tests = []struct {
		table TruthTable
		term  Term
	}{
		{
			TruthTable{
				False, False, False, False, False, False, False, False, False, False, False, False, False, False, False, False,
				False, False, False, False, False, False, False, False, False, False, False, False, False, False, False, False,
				False, False, False, False, False, False, False, False, False, False, False, False, False, False, False, False,
				False, False, False, False, False, False, False, False, False, False, False, False, False, False, False, False},
			Falsity,
		},
		{
			TruthTable{
				True, True, True, True, True, True, True, True, True, True, True, True, True, True, True, True,
				True, True, True, True, True, True, True, True, True, True, True, True, True, True, True, True,
				True, True, True, True, True, True, True, True, True, True, True, True, True, True, True, True,
				True, True, True, True, True, True, True, True, True, True, True, True, True, True, True, True},
			Truth,
		},
		{
			TruthTable{False, True},
			Variables[0],
		},
		{
			TruthTable{False, False, True, True},
			Variables[1],
		},
		{
			TruthTable{False, False, False, True},
			NewAND(Variables[0], Variables[1]),
		},
		{
			TruthTable{False, False, False, False, True, True, True, True},
			Variables[2],
		},
	}

	for _, tt := range tests {
		t.Run(tt.table.String(), func(t *testing.T) {
			mask, eval := tt.table.Evaluation()
			assert.Equal(t, tt.term.Eval()&mask, eval)
		})
	}
}

func TestParseTable(t *testing.T) {
	var tests = []struct {
		s     string
		table TruthTable
	}{
		{
			"",
			TruthTable{},
		},
		{
			"01",
			TruthTable{False, True},
		},
		{
			"01xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			TruthTable{False, True},
		},
		{
			"0x01",
			TruthTable{False, Undefined, False, True},
		},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			table, err := ParseTable(tt.s)
			require.NoError(t, err)
			require.Equal(t, *table, tt.table)
		})
	}
}
