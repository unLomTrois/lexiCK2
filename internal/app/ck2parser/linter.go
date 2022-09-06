package ck2parser

import (
	"fmt"
)

type Linter struct {
	Filepath   string  `json:"filepath"`
	Data       []*Node `json:"data"`
	Level      int     `json:"level"`
	towrite    []byte
	singleline bool
}

func NewLinter(file_path string, data []*Node) *Linter {
	return &Linter{
		Filepath:   file_path,
		Data:       data,
		Level:      0,
		towrite:    []byte{},
		singleline: false,
	}
}

func (l *Linter) Lint() {

	for _, node := range l.Data {
		l.LintNode(node)
		if node.Type != Comment {
			l.towrite = append(l.towrite, byte('\n'))
		}
	}

	fmt.Println("bytes:", len(l.towrite))
}

func (l *Linter) LintNode(node *Node) {
	// fmt.Println("node", node)

	switch node.Type {
	case Comment:
		l.LintComment(node)
	case Property, Comparison:
		l.LintProperty(node)
	case Entity, Block:
		l.LintBlock(node)
	default:
		panic("[Linter] unknown node type: " + node.Type)
	}
}

func (l *Linter) LintComment(node *Node) {
	if len(l.towrite) > 0 && l.towrite[len(l.towrite)-1] != ' ' {
		l.Intend()
	}
	l.towrite = append(l.towrite, []byte(node.Data.(string))...)
	l.towrite = append(l.towrite, byte('\n'))
}

func (l *Linter) LintProperty(node *Node) {
	if len(l.towrite) > 0 && l.towrite[len(l.towrite)-1] != ' ' {
		l.Intend()
	}

	l.towrite = append(l.towrite, node.KeyLiteral()...)
	l.towrite = append(l.towrite, byte(' '))
	l.towrite = append(l.towrite, []byte(node.Operator)...)
	l.towrite = append(l.towrite, byte(' '))
	l.towrite = append(l.towrite, node.DataLiteral()...)

	if l.singleline {
		l.towrite = append(l.towrite, byte(' '))
	} else {
		l.towrite = append(l.towrite, byte('\n'))
	}
}

func (l *Linter) LintBlock(node *Node) {
	children := node.Data.([]*Node)
	if len(children) == 1 && children[0].Type != Block {
		l.singleline = true
	}

	if len(l.towrite) > 0 && l.towrite[len(l.towrite)-1] != ' ' {
		l.Intend()
	}
	if !l.singleline {
		l.Level++
	}

	l.towrite = append(l.towrite, node.KeyLiteral()...)
	l.towrite = append(l.towrite, byte(' '))
	l.towrite = append(l.towrite, []byte(node.Operator)...)
	l.towrite = append(l.towrite, byte(' '))
	l.towrite = append(l.towrite, byte('{'))

	if l.singleline || children[0].Type == Comment {
		l.towrite = append(l.towrite, byte(' '))
	} else {
		l.towrite = append(l.towrite, byte('\n'))
	}

	for _, c := range children {

		l.LintNode(c)
	}

	if !l.singleline {
		l.Level--
		l.Intend()
	}
	l.towrite = append(l.towrite, byte('}'))
	l.towrite = append(l.towrite, byte('\n'))

	l.singleline = false
}

func (l *Linter) Intend() {
	i := 0
	for i < l.Level {
		l.towrite = append(l.towrite, byte('\t'))
		i++
	}
}

func (l *Linter) Next() *Node {
	if len(l.Data) == 0 {
		return nil
	}

	next := l.Data[0]
	// l.Data = l.Data[1:]
	return next
}

func (l *Linter) LintedData() []byte {
	return l.towrite
}
