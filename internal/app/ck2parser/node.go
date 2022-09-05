package ck2parser

type NodeType string

const (
	NextLine   NodeType = "NextLine"
	Comment    NodeType = "Comment"
	Entity     NodeType = "Entity"
	Block      NodeType = "Block"
	Property   NodeType = "Property"
	Comparison NodeType = "Comparison"
)

type Node struct {
	// Parent *any     `json:"-"`
	Type     NodeType    `json:"type"`
	Key      interface{} `json:"key,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Data     interface{} `json:"value,omitempty"`
}

func (n *Node) Node() *Node {
	return n.Data.(*Node)
}

func (n *Node) KeyLiteral() string {
	return n.Key.(string)
}

func (n *Node) DataLiteral() string {
	return n.Data.(string)
}
