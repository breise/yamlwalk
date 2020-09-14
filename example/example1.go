package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/breise/rstack"
	"github.com/breise/yamlwalk"
)

func main() {

	// Note the variation in the depths of the `unitPrice` nodes
	shoppingList := []byte(`
  - dept: dairy
    items:
      - qty: 2
        desc: dozen eggs
        unitPrice: 1.20
      - qty: 1
        desc: mozzaralla 0.5 lb
        unitPrice: 1.90
  - dept: snacks
    type:
      healthy:
        - qty: 4
          desc: organic dried apples 0.5 lb
          unitPrice: 2.35
        - qty: 2
          desc: broccoli chips 1 lb
          unitPrice: 3.21
      not so healthy:
        - desc: deep fried potato chips 1 lb
          unitPrice: 3.20
  `)

	// Print the total of the items' qty * unitPrice
	if err := yamlwalk.WalkDepthFirst(shoppingList, sumTotal); err != nil {
		log.Fatalf("Cannot yamlwalk.WalkDepthFirst(): %s", err)
	}
	fmt.Printf("Total is %0.2f\n\n", total)

	// Produce a yaml-patch document with which we can patch the original to add
	// a sibling node `ccy: USD` to each `unitPrice` node.  (Patching yaml
	// documents is outside the scope of this repo.)
	if err := yamlwalk.WalkDepthFirst(shoppingList, patchCcy); err != nil {
		log.Fatalf("Cannot yamlwalk.WalkDepthFirst(): %s", err)
	}
	fmt.Printf("Suitable as input to yaml-patch:%s\n\n", strings.Join(ccyPatchCommands, ``))

	// show the type of each node
	if err := yamlwalk.WalkDepthFirst(shoppingList, showTypes); err != nil {
		log.Fatalf("Cannot yamlwalk.WalkDepthFirst(): %s", err)
	}
}

var total float64

func sumTotal(node interface{}, ancestors *rstack.RStack) (bool, error) {
	if m, isMap := node.(map[interface{}]interface{}); isMap {
		if uP, ok := m["unitPrice"]; ok {
			qty, haveQty := m["qty"]
			if !haveQty {
				qty = 1
			}
			price := uP.(float64) * float64(qty.(int))
			total += price
		}
	}
	return false, nil
}

var ccyPatchCommands []string

func patchCcy(node interface{}, ancestors *rstack.RStack) (bool, error) {
	if yamlwalk.NodeIsScalar(node) {
		if _, isFloat := node.(float64); isFloat {
			progenitors, parent, err := ancestors.Pop()
			if err != nil {
				return false, err
			}
			if parent == "unitPrice" {
				ccyPatchCommand := fmt.Sprintf(`
- op: add
  path: /%s/ccy
  value: USD`, progenitors.Join(`/`))
				ccyPatchCommands = append(ccyPatchCommands, ccyPatchCommand)
			}
		}
	}
	return false, nil
}

func showTypes(node interface{}, ancestors *rstack.RStack) (bool, error) {
	fmt.Printf("Type: %T; Value: %v\n", node, node)
	return false, nil
}
