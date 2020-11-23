package boolean

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

// TruthValue defines possible values a boolean expression can have.
type TruthValue int

// Possible truth values
const (
	Undefined TruthValue = iota
	False
	True
)

func (v TruthValue) String() string {
	switch v {
	case Undefined:
		return "x"
	case False:
		return "0"
	case True:
		return "1"
	default:
		panic("unexpected value: " + strconv.Itoa(int(v)))
	}
}

// ParseValue returns the truth value represented by the string.
// It accepts x, undefined, UNDEFINED, Undefined, 0, f, F, FALSE, false, False, 1, t, T, TRUE, true, True.
// Any other value returns an error.
func ParseValue(s string) (TruthValue, error) {
	switch s {
	case "x", "undefined", "UNDEFINED", "Undefined":
		return Undefined, nil
	case "0", "f", "F", "false", "FALSE", "False":
		return False, nil
	case "1", "t", "T", "true", "TRUE", "True":
		return True, nil
	}
	return Undefined, fmt.Errorf("invalid truth value: %s", s)
}

// TruthTable defines a table of truth values.
type TruthTable [1 << MaxVars]TruthValue

// Evaluation computes the corresponding evaluation bit array and the bit mask for undefined values in the table.
func (table *TruthTable) Evaluation() (mask, eval uint64) {
	for i := 0; i < 64; i++ {
		// compute the corresponding table index
		var idx uint64
		for j := 0; j < MaxVars; j++ {
			// the i-th bit in varBits corresponds to the j-th bit in the table
			idx |= ((varBits[j] >> i) & 1) << j
		}
		// depending on the value at that particular index
		switch table[idx] {
		case Undefined:
		case False:
			mask |= 1 << i
		case True:
			eval |= 1 << i
			mask |= 1 << i
		default:
			panic("invalid truth value")
		}
	}
	return mask, eval
}

// NumVars computes and returns the number of relevant variables in the table.
func (table *TruthTable) NumVars() int {
	for i := uint(len(table)) - 1; i >= 0; i-- {
		if table[i] != Undefined {
			return bits.Len(i)
		}
	}
	return 0
}

func (table *TruthTable) String() string {
	var buf strings.Builder
	for _, v := range table {
		buf.WriteString(v.String())
	}
	return buf.String()
}

// ParseTable parses s as a truth table and returns the result.
func ParseTable(s string) (*TruthTable, error) {
	if l := len(s); l > len(TruthTable{}) {
		return nil, fmt.Errorf("invalid length: %d", l)
	}
	table := &TruthTable{}
	for i, c := range s {
		v, err := ParseValue(string(c))
		if err != nil {
			return nil, fmt.Errorf("index %d: %e", i, err)
		}
		table[i] = v
	}
	return table, nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (table *TruthTable) MarshalText() ([]byte, error) {
	return []byte(table.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (table *TruthTable) UnmarshalText(text []byte) error {
	x, err := ParseTable(string(text))
	if err != nil {
		return err
	}
	*table = *x
	return nil
}

// Print formats the table writes it to standard output.
func (table *TruthTable) Print() {
	numVars := table.NumVars()

	var columns []string
	for j := 0; j < numVars; j++ {
		columns = append(columns, strings.ToUpper(Variables[j].String()))
	}
	columns = append(columns, "OUTPUT")
	header := strings.Join(columns, " ")

	fmt.Println(header)
	fmt.Println(strings.Repeat("-", len(header)))

	for i := 0; i < 1<<numVars; i++ {
		var row []string
		for j := 0; j < numVars; j++ {
			row = append(row, fmt.Sprint((i>>j)&1))
		}
		row = append(row, table[i].String())
		fmt.Println(strings.Join(row, " "))
	}
}
