package ck2parser

import (
	"bufio"
	"ck2-parser/internal/app/lexer"
	"encoding/json"
	"fmt"
	"os"
)

type QueueParser struct {
	queue  []*lexer.Token
	cursor int
	data   []*Node
}

func NewQueueParser(queue []*lexer.Token) *QueueParser {
	return &QueueParser{
		queue:  queue,
		cursor: 0,
	}
}

func (p *QueueParser) Current() *lexer.Token {
	return p.queue[p.cursor]
}

func (q *QueueParser) hasMoreTokens() bool {
	return q.cursor < len(q.queue)-1
}

func (q *QueueParser) Next() *lexer.Token {
	if !q.hasMoreTokens() {
		return nil
	}

	q.cursor++

	return q.queue[q.cursor]
}

func (q *QueueParser) Parse() {
	root := &RootNode{
		Type:  Root,
		Data:  []*Node{},
		Scope: nil,
	}

	for {
		current := q.Current()
		fmt.Println(current)
		switch current.Type {
		case lexer.COMMENT:
			new_node := &Node{
				Type:   Comment,
				Parent: nil,
				Key:    "",
				Value:  string(current.Value),
				Data:   nil,
			}
			root.InsertNode(new_node)
			q.Next()
		case lexer.WORD:
			fmt.Println(string(current.Value))
			next := q.Next()
			if next.Type != lexer.EQUALS {
				continue
			}
			fmt.Println(string(next.Value))
			next = q.Next()
			// handle element
			if next.Type == lexer.WORD {
				fmt.Println(string(next.Value))

				property_node := &Node{
					Type:   Property,
					Parent: nil,
					Key:    string(current.Value),
					Value:  string(next.Value),
					Data:   nil,
				}
				root.InsertNode(property_node)

				// q.Next()
				// continue
			}
			if next.Type == lexer.START {
				fmt.Println("kek", string(next.Value))

				block_node := &Node{
					Type:   Block,
					Parent: nil,
					Key:    string(current.Value),
					Value:  string(next.Value),
					Data:   []*Node{},
				}
				root.InsertNode(block_node)

				q.Next()
				continue
			}

		}
		if !q.hasMoreTokens() {
			break
		}
		// fmt.Println(q.Next())

	}

	new_file, err := os.Create("./tmp/meta.json")
	if err != nil {
		panic(err)
	}
	defer new_file.Close()

	w := bufio.NewWriter(new_file)

	JSON, _ := json.MarshalIndent(root, "", "  ")
	w.Write(JSON)
	w.Flush()

}
