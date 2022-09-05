package ck2parser

import (
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

func (l *Linter) Lint() {

	for _, node := range l.Data {
		l.LintNode(node)
	}

	// l.LintNode()

	fmt.Println("to write:", strconv.Quote(string(l.towrite)))
	fmt.Println("bytes:", len(l.towrite))
}

func (l *Linter) LintNode(node *Node) {
	fmt.Println("node", node)
	if node.Type == Comment {
		if len(l.towrite) > 0 && l.towrite[len(l.towrite)-1] != ' ' {
			l.Intend()
		}

		l.towrite = append(l.towrite, []byte(node.Data.(string))...)
		l.towrite = append(l.towrite, byte('\n'))
	}
	if node.Type == Property {
		l.Intend()
		l.towrite = append(l.towrite, []byte(node.KeyLiteral())...)
		l.towrite = append(l.towrite, byte(' '))
		l.towrite = append(l.towrite, []byte(node.Operator)...)
		l.towrite = append(l.towrite, byte(' '))
		l.towrite = append(l.towrite, []byte(node.DataLiteral())...)
		l.towrite = append(l.towrite, byte('\n'))
	}
	if node.Type == Block {
		l.Intend()
		l.Level++

		l.towrite = append(l.towrite, []byte(node.KeyLiteral())...)
		l.towrite = append(l.towrite, byte(' '))
		l.towrite = append(l.towrite, []byte(node.Operator)...)
		l.towrite = append(l.towrite, byte(' '))
		l.towrite = append(l.towrite, byte('{'))

		if node.Data.([]*Node)[0].Type == Comment {
			l.towrite = append(l.towrite, byte(' '))
		} else {
			l.towrite = append(l.towrite, byte('\n'))
		}

		for _, c := range node.Data.([]*Node) {
			// fmt.Println(c)

			l.LintNode(c)
		}

		l.Level--
		l.Intend()
		l.towrite = append(l.towrite, byte('}'))
		l.towrite = append(l.towrite, byte('\n'))

		// json, _ := json.MarshalIndent(node.Data, "", " ")
		// fmt.Println(string(json))

	}
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
