package main

// Author: johnmatthiggins@gmail.com
// Program: JSX Import Sorter
//
// A program that sorts JSX imports. It can be invoked by
// hitting "!" and typing the name of the binary after
// selecting some text in visual mode.

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"
)

type ImportStatement struct {
	componentName string
	importPath    string
	rawLine       string
}

func main() {
	var lines []string

	// Read in lines until you get to EOT character
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "." {
			break
		}
		lines = append(lines, line)
	}

	statements := make([]ImportStatement, len(lines))

	for i := 0; i < len(statements); i++ {
		importStatement, err := parseImportStatement(lines[i])
		if err != nil {
			// print out original text if we can't figure out
			// how to parse it...
			fmt.Print(strings.Join(lines, "\n"))
			fmt.Fprintf(os.Stderr, "Could not parse lines...\n")
			return
		}
		statements[i] = *importStatement
	}

	slices.SortFunc(statements, compareImportStatements)
	for _, statement := range statements {
		statement.PrintStatement()
	}
}

func compareImportStatements(a, b ImportStatement) int {
	result := cmp.Compare(a.importPath, b.importPath)
	if result == 0 {
		result = cmp.Compare(a.componentName, b.componentName)
	}
	return result
}

func (stmt *ImportStatement) PrintStatement() {
	var sb strings.Builder

	for _, c := range []rune(stmt.rawLine) {
		if c == ' ' {
			sb.WriteString(" ")
		} else {
			break
		}
	}

	fmt.Printf("%simport %s from %s\n", sb.String(), stmt.componentName, stmt.importPath)
}

func parseImportStatement(line string) (*ImportStatement, error) {
	segments := strings.Split(strings.TrimSpace(line), " ")
	if len(segments) != 4 {
		err := fmt.Errorf(
			"Expected 4 segments after splitting on space character, but found %d.",
			len(segments),
		)
		return nil, err
	}
	importStatement := ImportStatement{
		componentName: segments[1],
		importPath:    segments[3],
		rawLine:       line,
	}

	return &importStatement, nil
}
