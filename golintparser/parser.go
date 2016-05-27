package golintparser

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/ifn/go-junit-report/parser"
)

var (
	regexGolintLine = regexp.MustCompile(`\.go:\d+:\d+: `)
)

// A GolintParser represents golint report parser.
type GolintParser struct{}

// New returns a new GolintParser.
func New() *GolintParser {
	return &GolintParser{}
}

// Parse parses golint report and forms it as the test report.
func (glp *GolintParser) Parse(r io.Reader, _pkgName string) (*parser.Report, error) {
	reader := bufio.NewReader(r)

	report := &parser.Report{make([]parser.Package, 0)}
	var tests []*parser.Test
	var prevPkgName string

	// parse lines
	for {
		l, _, err := reader.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		line := string(l)

		// skip invalid lines
		if !regexGolintLine.MatchString(line) {
			continue
		}

		// package name is a file path
		pkgName := strings.Split(line, ":")[0]
		pkgName = strings.Replace(strings.TrimPrefix(pkgName, "/"), "/", ".", -1)
		// test name is the error position in line
		testName := strings.Join(strings.Split(line, ":")[1:3], ":")

		// line corresponds to a new file
		if prevPkgName != pkgName {
			report.Packages = append(report.Packages, parser.Package{
				Name: pkgName,
			})
			tests = make([]*parser.Test, 0)
			prevPkgName = pkgName
		}

		// append test case to a package
		tests = append(tests, &parser.Test{
			Name:   testName,
			Result: parser.FAIL,
			Output: []string{line},
		})
		report.Packages[len(report.Packages)-1].Tests = tests
	}

	return report, nil
}
