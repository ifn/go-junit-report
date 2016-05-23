package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jstemmer/go-junit-report/golintparser"
	"github.com/jstemmer/go-junit-report/gotestparser"
	"github.com/jstemmer/go-junit-report/parser"
)

var (
	noXMLHeader bool
	packageName string
	setExitCode bool
	golint      bool
)

func init() {
	flag.BoolVar(&noXMLHeader, "no-xml-header", false, "do not print xml header")
	flag.StringVar(&packageName, "package-name", "", "specify a package name (compiled test have no package name in output)")
	flag.BoolVar(&setExitCode, "set-exit-code", false, "set exit code to 1 if tests failed")
	flag.BoolVar(&golint, "golint", false, "parse input as golint output")
}

func main() {
	flag.Parse()

	// Determine the parser
	var pr parser.Parser
	if golint {
		pr = golintparser.New()
	} else {
		pr = gotestparser.New()
	}

	// Read input
	report, err := pr.Parse(os.Stdin, packageName)
	if err != nil {
		fmt.Printf("Error reading input: %s\n", err)
		os.Exit(1)
	}

	// Write xml
	err = JUnitReportXML(report, noXMLHeader, os.Stdout)
	if err != nil {
		fmt.Printf("Error writing XML: %s\n", err)
		os.Exit(1)
	}

	if setExitCode && report.Failures() > 0 {
		os.Exit(1)
	}
}
