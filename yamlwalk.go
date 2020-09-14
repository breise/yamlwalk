package yamlwalk

import (
	"fmt"
	"strconv"

	rst "github.com/breise/rstack"
	yaml "gopkg.in/yaml.v2"
)

func WalkDepthFirst(b []byte, fn func(node interface{}, ancestors *rst.RStack) error) error {
	var v interface{}
	if err := yaml.Unmarshal(b, &v); err != nil {
		return fmt.Errorf("WalkDepthFirst(): Cannot Unmarshal(): %s", err)
	}
	if err := wdf(v, rst.New(), fn); err != nil {
		return err
	}
	return nil
}

func wdf(node interface{}, ancestors *rst.RStack, fn func(node interface{}, ancestors *rst.RStack) error) error {
	if err := fn(node, ancestors); err != nil {
		return err
	}
	if thisMap, isMap := node.(map[interface{}]interface{}); isMap {
		for k, v := range thisMap {
			nextAncestors := ancestors.Push(fmt.Sprintf("%v", k))
			if err := wdf(v, nextAncestors, fn); err != nil {
				return err
			}
		}
	} else if thisArray, isArray := node.([]interface{}); isArray {
		for i, v := range thisArray {
			nextAncestors := ancestors.Push(strconv.Itoa(i))
			if err := wdf(v, nextAncestors, fn); err != nil {
				return err
			}
		}
	}
	// node is a scalar.  We've already called fn().  Nothing further to do.
	return nil
}

func NodeIsMap(node interface{}) bool {
	_, rv := node.(map[interface{}]interface{})
	return rv
}

func NodeIsArray(node interface{}) bool {
	_, rv := node.([]interface{})
	return rv
}

func NodeIsScalar(node interface{}) bool {
	return !(NodeIsMap(node) || NodeIsArray(node))
}
