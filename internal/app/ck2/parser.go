package ck2

import (
	"bytes"
	"fmt"
	"strconv"
)

type CK2Parser struct {
	Namespace string  `json:"namespace"`
	Depth     int     `json:"depth"`
	Nodes     []*Node `json:"elements"`
	Scope     *Node   `json:"-"`
	PrevScope *Node   `json:"-"`
}

type NodeType string

const (
	Entity   NodeType = "Entity"
	Block    NodeType = "Block"
	Property NodeType = "Property"
)

type Node struct {
	// element is either a Block or a Property
	Type NodeType `json:"type"`
	// not nil if type is Property
	Key   string `json:"key"`
	Value string `json:"value"`
	// not nil if type is Block
	Data []*Node `json:"data"`
}

func NewParser() *CK2Parser {
	return &CK2Parser{
		Namespace: "",
		Depth:     0,
		Nodes:     []*Node{},
		Scope:     nil,
		PrevScope: nil,
	}
}

func (parser *CK2Parser) ParseLine(line []byte) []byte {
	if bytes.Contains(line, []byte("=")) {
		kv := bytes.SplitN(line, []byte(" = "), 2)
		key := kv[0]
		value := kv[1]
		fmt.Println("key:", strconv.Quote(string(key)), "value:", strconv.Quote(string(value)))

		// namespace
		if parser.Scope == nil && parser.Namespace == "" {
			parser.Namespace = string(value)
		} else {
			fmt.Println("ELSE")
			if value[0] == byte('{') {
				// enter into entity scope
				if parser.Scope == nil {
					parser.Nodes = append(parser.Nodes, &Node{
						Type:  Entity,
						Data:  []*Node{},
						Key:   string(key),
						Value: "",
					})
					parser.Scope = parser.Nodes[len(parser.Nodes)-1]
					parser.PrevScope = parser.Scope
				} else {
					// enter another scope
					fmt.Println("enter into scope of:", strconv.Quote(string(key)))
					parser.Scope.Data = append(parser.Scope.Data, &Node{
						Type:  Block,
						Data:  []*Node{},
						Key:   string(key),
						Value: string(value),
					})
					parser.Scope = parser.Scope.Data[len(parser.Scope.Data)-1]
					fmt.Println("scope:", parser.Scope)
				}
			} else {
				parser.Scope.Data = append(parser.Scope.Data, &Node{
					Type:  Property,
					Data:  nil,
					Key:   string(key),
					Value: string(value),
				})
				fmt.Println("property:", strconv.Quote(string(key)), "value:", strconv.Quote(string(value)))
			}
		}
	}

	if len(line) > 0 && line[0] == byte('}') {
		fmt.Println("END of scope")
		if parser.Scope == parser.PrevScope {
			parser.Scope = nil
		} else {
			parser.Scope = parser.PrevScope
		}
	}

	fmt.Print("\n")

	return line
}
