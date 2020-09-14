# yamlwalk

Walk a yaml or json tree, performing an action at each node.

```
import (
	"github.com/breise/rstack"
	"github.com/breise/yamlwalk"
)

func main() {
  ...
  if err := yamlwalk.WalkDepthFirst(yamlDoc, myActionFunc); err != nil {
    log.Fatalf("Cannot yamlwalk.WalkDepthFirst(): %s", err)
  }
  fmt.Printf("The answer is: %v", myResult)
}

var myResult ...

func myActionFunc(node interface{}, ancestors *rstack.RStack) (bool, error) {
  ...
  // update myResult
}
```

As each node is visited, `myActionFunc` is invoked.  Its signature is
```
  func myActionFunc(node interface{}, ancestors *rstack.RStack) (bool, error)
```
`myActionFunc()` is passed the current node and its ancestors, and returns 2 values:
1. a `bool`
  * A `true` return value indicates to `WalkDepthFirst()` to _prune_ the tree at this point.  That is, do not proceed further down the current branch.
  * You can use this in `myActionFunc()`, returning `true` when encountering a node whose children you _don't_ want to traverse, and returning `false` when encountering a node whose children you _do_ want to traverse
  * If you want to traverse the entire tree, always return `false` from `myActionFunc()`
1. an `error`. A non-nil error causes `WalkDepthFirst()` to stop processing and report the error.

## `node`

Inside of `myActionFunc()`, the first thing you will probably want to do is determine what type of node you have.
```
  if m, isMap := node.(map[interface{}]interface{}); isMap {
    // work with the map, using type assertions as necessary
  } else if a, isArray := node.([]interface{}); isArray {
    // work with the array, using type assertions as necessary
  } else {
    // We have a scalar. Work with it, using type assertions as necessary
  }
```
If you are only interested in scalars (leaves in the yaml or json tree), you
can use the `NodeIsScalar()` convenience function, which basically performs the
check above, returning `false` if `isMap` or `isArray`, true otherwise:
```
  if yamlwalk.NodeIsScalar(node) {
    // Work with the scalar, using type assertions as necessary
  }
```

## `ancestors`

While working with your node, you may wish to examine or list its parent-key,
its parent's parent-key, etc... in other words, its ancestors.

`ancestors` is an [`*RStack`](https://github.com/breise/rstack).

`RStack` is a recursive stack, meaning that every node of the stack is itself an RStack.

The methods on `RStack` you may find useful in `myActionFunc()` are:
- `func (s *RStack) Pop() (*RStack, interface{}, error)`
- `func (s *RStack) ToSlice() []interface{}`
- `func (s *RStack) ToStringSlice() []string`
- `func (s *RStack) Join(sep string) string`

## The Closure

Since `myActionFunc()` does not return any value (other than the error type),
you will have to depend on side effects.  You can `Println()` or write to a
file from `myActionFunc()`, but you may prefer to store results in a closed
variable, that is, a package variable around which `myActionFunc()` forms a
closure.

See yamlwalk/example/example1.go for examples of two closed variables, 
```
var total float64
```
around which `func sumTotal()` forms a closure, and
```
var ccyPatchCommands []string
```
around which `func patchCcy()` forms a closure. 

In each case, `func main()`, after invoking `yamlwalk.WalkDepthFirst()`, prints
the value of the closed variable.

## Example

- `yamlwalk/example/example1.go`
