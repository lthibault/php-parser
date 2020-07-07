package parser_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	parser "github.com/lthibault/php-parser/pkg"
	"github.com/lthibault/php-parser/pkg/ast"
)

const testdata = "testdata"

func TestParser(t *testing.T) {
	var b bytes.Buffer
	for _, tc := range testfiles() {
		t.Run(tc.Name(), func(t *testing.T) {
			defer b.Reset()

			require.NoError(t, load(tc, &b))

			root, errs := parser.New(b.String()).Parse()
			assertNoErrors(t, errs, "parse error (%s)", tc.Name())
			tc.assertExpectedAST(t, root)
		})
	}
}

func assertNoErrors(t *testing.T, errs []error, msgAndArgs ...interface{}) bool {
	if assert.Empty(t, errs, "errors were encountered") {
		return true
	}

	for _, err := range errs {
		t.Log(err)
	}

	return false
}

func load(info os.FileInfo, w io.Writer) error {
	f, err := os.Open(relpath(info))
	if err != nil {
		return errors.Wrap(err, "fopen")
	}
	defer f.Close()

	io.Copy(w, f)
	return nil
}

func relpath(info os.FileInfo) string {
	return filepath.Join(testdata, info.Name())
}

type testCase struct {
	os.FileInfo
}

func (tc testCase) assertExpectedAST(t *testing.T, root []ast.Node) bool {
	if tc.noAST() {
		return assert.Empty(t, root, "unexpected AST for %s", tc.Name())
	}

	return assert.NotEmpty(t, root, "missing AST for %s", tc.Name())
}

func (tc testCase) noAST() bool {
	for _, part := range strings.Split(tc.Name(), ".") {
		if part == "noast" {
			return true
		}
	}

	return false
}

func testfiles() []testCase {
	fs, err := ioutil.ReadDir(testdata)
	if err != nil {
		panic(err)
	}

	tcs := make([]testCase, len(fs))
	for i, info := range fs {
		tcs[i].FileInfo = info
	}

	return tcs
}
