package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/wollac/boolforcer/pkg/boolean"
	"github.com/wollac/boolforcer/pkg/forcer"
)

var (
	table = flag.String(
		"table",
		"1110100010000000",
		"truth table",
	)
	complexity = flag.Int(
		"complexity",
		16,
		"maximum complexity considered",
	)
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	t, err := boolean.ParseTable(*table)
	if err != nil {
		return fmt.Errorf("invalid truth table: %e", err)
	}
	t.Print()
	fmt.Println()

	f := forcer.New(t)
	solutions := f.Run(*complexity)
	if len(solutions) == 0 {
		return errors.New("no solution found")
	}
	fmt.Println("The following minimal solutions have been found:")
	for i, s := range solutions {
		fmt.Printf(" %2d: %v with complexity %d\n", i+1, s, s.Complexity())
	}
	return nil
}
