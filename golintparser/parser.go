package golintparser

import (
	"bufio"
	"io"
	"regexp"
	"strings"
	"path/filepath"

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

	var report = &parser.Report{make([]parser.Package, 0)}
	var prevPath string
	var test *parser.Test

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

		path := strings.Split(line, ":")[0]
		dir := filepath.Dir(path)
		dir = strings.Replace(strings.TrimPrefix(dir, "/"), "/", ".", -1)
		name := filepath.Base(path)

		// line corresponds to a new file
		if prevPath != path {
			prevPath = path

			test = &parser.Test{
				Name:   name,
				Result: parser.FAIL,
				Output: make([]string, 0),
			}
			pkg := parser.Package{
				Name: dir,
				Tests: []*parser.Test{test},
			}
			report.Packages = append(report.Packages, pkg)
		}

		test.Output = append(test.Output, line)
	}

	return report, nil
}
