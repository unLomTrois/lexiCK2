package ck2

import (
	"bytes"
	"fmt"
	"strconv"
)

type CK2Parser struct {
	Namespace string     `json:"namespace"`
	Depth     int        `json:"depth"`
	Elements  []*Element `json:"elements"`
	Scope     *Element   `json:"-"`
	PrevScope *Element   `json:"-"`
}

type ElementType string

const (
	Entity   ElementType = "Entity"
	Block    ElementType = "Block"
	Property ElementType = "Property"
)

type Element struct {
	// element is either a Block or a Property
	Type ElementType `json:"type"`
	// not nil if type is Property
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
	// not nil if type is Block
	Data []*Element `json:"data,omitempty"`
}

func NewParser() *CK2Parser {
	return &CK2Parser{
		Namespace: "",
		Depth:     0,
		Elements:  []*Element{},
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
					parser.Elements = append(parser.Elements, &Element{
						Type:  Entity,
						Data:  []*Element{},
						Key:   string(key),
						Value: "",
					})
					parser.Scope = parser.Elements[len(parser.Elements)-1]
					parser.PrevScope = parser.Scope
				} else {
					// enter another scope
					fmt.Println("enter into scope of:", strconv.Quote(string(key)))
					parser.Scope.Data = append(parser.Scope.Data, &Element{
						Type:  Block,
						Data:  []*Element{},
						Key:   string(key),
						Value: string(value),
					})
					parser.Scope = parser.Scope.Data[len(parser.Scope.Data)-1]
					fmt.Println("scope:", parser.Scope)
				}
			} else {
				parser.Scope.Data = append(parser.Scope.Data, &Element{
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
