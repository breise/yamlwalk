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
	prune      bool
}{
	{name: "Pets as yaml", inputFile: "testdata/pets_inp.yaml", expectFile: "testdata/pets_exp.txt"},
	{name: "Pets as json", inputFile: "testdata/pets_inp.json", expectFile: "testdata/pets_exp.txt"},
	{name: "Pets as yaml, prune_arrays", inputFile: "testdata/pets_inp.yaml", expectFile: "testdata/pets_prune_exp.txt", prune: true},
	{name: "Pets as json, prune_arrays", inputFile: "testdata/pets_inp.json", expectFile: "testdata/pets_prune_exp.txt", prune: true},
	{name: "Array as yaml", inputFile: "testdata/array_inp.yaml", expectFile: "testdata/array_exp.txt"},
	{name: "Map as yaml", inputFile: "testdata/map_inp.yaml", expectFile: "testdata/map_exp.txt"},
	{name: "json with nulls", inputFile: "testdata/nulls_inp.json", expectFile: "testdata/nulls_exp.txt"},
	{name: "yaml with numeric keys", inputFile: "testdata/num_inp.yaml", expectFile: "testdata/num_exp.txt"},
}

var prune bool

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
			prune = tc.prune
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

func listPaths(node interface{}, ancestors *rstack.RStack) (bool, error) {
	if prune && NodeIsArray(node) {
		return true, nil
	}
	if NodeIsScalar(node) {
		els := []string{``} // force leading `/` in Join()
		els = append(els, ancestors.ToStringSlice()...)
		path := strings.Join(els, "/")
		n := fmt.Sprintf("%v", node) // node is a scalar, but we don't know its type.  Make it a string
		pathValue := strings.Join([]string{path, n}, ": ")
		pathValues = append(pathValues, pathValue)
	}
	return false, nil
}
