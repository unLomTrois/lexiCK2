package ck2

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
