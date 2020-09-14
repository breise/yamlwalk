package yamlwalk

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"testing"

	rstack "github.com/breise/rstack"
)

var cases = []struct {
	name       string
	inputFile  string
	expectFile string
}{
	{"Pets as yaml", "testdata/pets_inp.yaml", "testdata/pets_exp.txt"},
	{"Pets as json", "testdata/pets_inp.json", "testdata/pets_exp.txt"},
	{"Array as yaml", "testdata/array_inp.yaml", "testdata/array_exp.txt"},
	{"Map as yaml", "testdata/map_inp.yaml", "testdata/map_exp.txt"},
	{"json with nulls", "testdata/nulls_inp.json", "testdata/nulls_exp.txt"},
}

func TestYamlWalk(t *testing.T) {
	for i, tc := range cases {
		desc := fmt.Sprintf("Test Case %d: %s", i, tc.name)
		t.Run(desc, func(t *testing.T) {
			b, err := ioutil.ReadFile(tc.inputFile)
			if err != nil {
				t.Fatalf("Cannot read file '%s'", tc.inputFile)
			}
			expB, err := ioutil.ReadFile(tc.expectFile)
			if err != nil {
				t.Fatalf("Cannot read file '%s'", tc.expectFile)
			}
			pathValues = []string{}
			if err := WalkDepthFirst(b, listPaths); err != nil {
				t.Fatalf("Cannot WalkDepthFirst(). Error: %s", err)
			}
			sort.Strings(pathValues)
			got := strings.TrimSpace(strings.Join(pathValues, "\n"))
			exp := strings.TrimSpace(string(expB))
			if exp != got {
				t.Errorf("%s:\nExp:\n{{{%s}}}\nGot:\n{{{%s}}}\n", desc, exp, got)
			}
		})
	}
}

var pathValues []string

func listPaths(node interface{}, ancestors *rstack.RStack) error {
	if NodeIsScalar(node) {
		els := []string{``} // force leading `/` in Join()
		els = append(els, ancestors.ToStringSlice()...)
		path := strings.Join(els, "/")
		n := fmt.Sprintf("%v", node) // node is a scalar, but we don't know its type.  Make it a string
		pathValue := strings.Join([]string{path, n}, ": ")
		pathValues = append(pathValues, pathValue)
	}
	return nil
}
