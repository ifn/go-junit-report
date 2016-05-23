package parser

import (
	"io"
)

// Parser is the interface that wraps the Parse method.
//
// Parse reads from r and forms the test report.
// pkgName may be used to name a package,
// e.g. if it was not determined for any test collection while parsing.
type Parser interface {
	Parse(r io.Reader, pkgName string) (*Report, error)
}

// Result represents a test result.
type Result int

// Test result constants
const (
	PASS Result = iota
	FAIL
	SKIP
)

// Report is a collection of package tests.
type Report struct {
	Packages []Package
}

// Failures counts the number of failed tests in this report
func (r *Report) Failures() int {
	count := 0

	for _, p := range r.Packages {
		for _, t := range p.Tests {
			if t.Result == FAIL {
				count++
			}
		}
	}

	return count
}

// Package contains the test results of a single package.
type Package struct {
	Name        string
	Time        int
	Tests       []*Test
	CoveragePct string
}

// Test contains the results of a single test.
type Test struct {
	Name   string
	Time   int
	Result Result
	Output []string
}
