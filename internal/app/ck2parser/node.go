package ck2parser

type NodeType string

const (
	Root     NodeType = "Root"
	NextLine NodeType = "NextLine"
	Comment  NodeType = "Comment"
	Entity   NodeType = "Entity"
	Block    NodeType = "Block"
	Property NodeType = "Property"
)

type RootNode struct {
	Type  NodeType     `json:"type"`
	Data  []*Node      `json:"data"`
	Scope *interface{} `json:"-"`
}

func (root *RootNode) InsertNode(node *Node) {
	root.Data = append(root.Data, node)
}

type Node struct {
	// element is either a Block or a Property
	Type   NodeType `json:"type"`
	Parent *any     `json:"-"`
	// Level  int      `json:"level"`
	// not nil if type is Property
	Key   string `json:"key,omitempty"`
	Value string `json:"value"`
	// not nil if type is Block
	Data []*Node `json:"data,omitempty"`
}

func NewNode(node_type NodeType, key string, value string) *Node {
	return &Node{
		Type:   node_type,
		Parent: nil,
		Key:    key,
		Value:  value,
		Data:   []*Node{},
	}
}
