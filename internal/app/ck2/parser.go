package ck2

import (
	"bytes"
	"fmt"
	"strconv"
)

type CK2Parser struct {
	Entities  []*Entity `json:"entities"`
	Scope     *Entity   `json:"scope"`
	PrevScope *Entity
	Depth     int `json:"depth"`
}

type Entity struct {
	Name     string           `json:"name"`
	Elements []*EntityElement `json:"elements"`
}

type EntityType string

const (
	Property EntityType = "Property"
	Block    EntityType = "Block"
)

type EntityElement struct {
	// element is either a Block or a Property
	Type EntityType `json:"type"`
	// not nil if type is Block
	Data []*EntityElement `json:"data"`
	// not nil if type is Property
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (parser *CK2Parser) ParseLine(line []byte) []byte {
	if bytes.Contains(line, []byte("=")) {
		kv := bytes.Split(line, []byte(" = "))
		key := kv[0]
		value := kv[1]
		fmt.Println("key:", strconv.Quote(string(key)), "value:", strconv.Quote(string(value)))

		if value[0] == byte('{') {
			// enter into entity scope
			if parser.Scope == nil {
				parser.Entities = append(parser.Entities, &Entity{
					Name:     string(key),
					Elements: []*EntityElement{},
				})
				parser.Scope = parser.Entities[len(parser.Entities)-1]
				parser.PrevScope = parser.Scope
			} else {
				// enter another scope
				fmt.Println("enter into scope of:", strconv.Quote(string(key)))
				parser.Scope.Elements = append(parser.Scope.Elements, &EntityElement{
					Type:  Block,
					Data:  []*EntityElement{},
					Key:   string(key),
					Value: "",
				})
				parser.Scope = parser.Entities[len(parser.Entities)-1]
				fmt.Println("scope:", parser.Scope)
				// last.Elements = append(last.Elements, &EntityElement{
				// 	_type: Property,
				// 	Data:  nil,
				// 	Key:   string(key),
				// 	Value: string(value),
				// })
			}
		} else {
			parser.Scope.Elements = append(parser.Scope.Elements, &EntityElement{
				Type:  Property,
				Data:  nil,
				Key:   string(key),
				Value: string(value),
			})
			fmt.Println("property:", strconv.Quote(string(key)), "value:", strconv.Quote(string(value)))
		}
	}
	if len(line) > 0 && line[0] == byte('}') {
		fmt.Println("END of scope")
		parser.Scope = parser.PrevScope
	}

	fmt.Print("\n")

	return line
}
