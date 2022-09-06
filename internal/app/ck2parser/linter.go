package ck2parser

import (
	"fmt"
	"strconv"
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

	// l.LintNode()

	fmt.Println("to write:", strconv.Quote(string(l.towrite)))
	fmt.Println("bytes:", len(l.towrite))
}

func (l *Linter) LintNode(node *Node) {
	// fmt.Println("node", node)

	if node.Type == Comment {
		if len(l.towrite) > 0 && l.towrite[len(l.towrite)-1] != ' ' {
			l.Intend()
		}

		l.towrite = append(l.towrite, []byte(node.Data.(string))...)
		l.towrite = append(l.towrite, byte('\n'))
	}
	if node.Type == Property {
		if len(l.towrite) > 0 && l.towrite[len(l.towrite)-1] != ' ' {
			l.Intend()
		}
		l.towrite = append(l.towrite, []byte(node.KeyLiteral())...)
		l.towrite = append(l.towrite, byte(' '))
		l.towrite = append(l.towrite, []byte(node.Operator)...)
		l.towrite = append(l.towrite, byte(' '))
		l.towrite = append(l.towrite, []byte(node.DataLiteral())...)

		if l.singleline {
			l.towrite = append(l.towrite, byte(' '))
		} else {
			l.towrite = append(l.towrite, byte('\n'))
		}
	}
	if node.Type == Block {
		if len(node.Data.([]*Node)) == 1 && (node.Key == "NOT" || node.Key == "limit") {
			l.singleline = true
		}

		// l.Intend()
		fmt.Println("node", node)

		if len(l.towrite) > 0 && l.towrite[len(l.towrite)-1] != ' ' {
			l.Intend()
		}
		if !l.singleline {
			l.Level++
		}

		l.towrite = append(l.towrite, []byte(node.KeyLiteral())...)
		l.towrite = append(l.towrite, byte(' '))
		l.towrite = append(l.towrite, []byte(node.Operator)...)
		l.towrite = append(l.towrite, byte(' '))
		l.towrite = append(l.towrite, byte('{'))

		if node.Data.([]*Node)[0].Type == Comment || l.singleline {
			l.towrite = append(l.towrite, byte(' '))
		} else {
			l.towrite = append(l.towrite, byte('\n'))
		}

		for _, c := range node.Data.([]*Node) {
			// fmt.Println(c)

			l.LintNode(c)
		}

		fmt.Println(l.singleline)
		if !l.singleline {
			l.Level--
			l.Intend()
		}
		l.towrite = append(l.towrite, byte('}'))
		l.towrite = append(l.towrite, byte('\n'))

		// json, _ := json.MarshalIndent(node.Data, "", " ")
		// fmt.Println(string(json))

		l.singleline = false
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
	// l.Data = l.Data[1:]
	return next
}

func (l *Linter) LintedData() []byte {
	return l.towrite
}
