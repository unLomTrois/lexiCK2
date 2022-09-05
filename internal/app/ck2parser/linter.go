package ck2parser

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Linter struct {
	Filepath string  `json:"filepath"`
	Data     []*Node `json:"data"`
	Level    int     `json:"level"`
	towrite  []byte
}

func NewLinter(file_path string, data []*Node) *Linter {
	return &Linter{
		Filepath: file_path,
		Data:     data,
		towrite:  []byte{},
	}
}

func (l *Linter) Lint(node *Node) {
	fmt.Println(node)

	if node.Type == Comment {
		l.towrite = append(l.towrite, []byte(node.Data.(string))...)
		l.towrite = append(l.towrite, byte('\n'))
	}
	if node.Type == Property {
		l.Intend()
		l.towrite = append(l.towrite, []byte(node.Key.(string))...)
		l.towrite = append(l.towrite, byte(' '))
		l.towrite = append(l.towrite, []byte(node.Operator)...)
		l.towrite = append(l.towrite, byte(' '))
		l.towrite = append(l.towrite, []byte(node.Data.(string))...)
		l.towrite = append(l.towrite, byte('\n'))
	}
	if node.Type == Block {
		l.Level++

		l.towrite = append(l.towrite, []byte(node.Key.(string))...)
		l.towrite = append(l.towrite, byte(' '))
		l.towrite = append(l.towrite, []byte(node.Operator)...)
		l.towrite = append(l.towrite, byte(' '))
		l.towrite = append(l.towrite, byte('{'))

		if node.Data.([]*Node)[0].Type == Comment {
			l.towrite = append(l.towrite, byte(' '))
		} else {
			l.towrite = append(l.towrite, byte('\n'))
			l.Intend()
		}

		// iterate over children
		child1 := node.Data.([]*Node)[0]
		child2 := node.Data.([]*Node)[1]
		// fmt.Println("child1", child1)
		// fmt.Println("child2", child2)
		l.Lint(child1)
		l.Lint(child2)
		l.Lint(node.Data.([]*Node)[2])

		l.Level--
		l.towrite = append(l.towrite, byte('}'))

		json, _ := json.MarshalIndent(node.Data, "", " ")
		fmt.Println(string(json))

	}
	fmt.Println("to write:", strconv.Quote(string(l.towrite)))
	fmt.Println("bytes:", len(l.towrite))
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
	l.Data = l.Data[1:]
	return next
}

func (l *Linter) LintedData() []byte {
	return l.towrite
}
