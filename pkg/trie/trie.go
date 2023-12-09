package trie

import (
	"fmt"
	"strings"
)

type TrieNode struct {
	children map[rune]*TrieNode
	isWord   bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		isWord:   false,
	}
}

type Filter struct {
	root *TrieNode
}

func NewFilter() *Filter {
	return &Filter{
		root: NewTrieNode(),
	}
}

func (f *Filter) AddWord(word string) {
	node := f.root
	for _, ch := range word {
		if _, ok := node.children[ch]; !ok {
			node.children[ch] = NewTrieNode()
		}
		node = node.children[ch]
	}
	node.isWord = true
}

func (f *Filter) ContainsWord(word string) bool {
	node := f.root
	for _, ch := range word {
		if _, ok := node.children[ch]; !ok {
			return false
		}
		node = node.children[ch]
	}
	return node.isWord
}

func main() {
	words := []string{"word1", "word2"} // список слов-фильтров
	filter := NewFilter()
	for _, word := range words {
		filter.AddWord(word)
	}

	messages := []string{"msg1", "msg2"} // список сообщений
	for _, message := range messages {
		// Проверка наличия слов-фильтров в сообщении
		for _, word := range strings.Fields(message) {
			if filter.ContainsWord(word) {
				// Слово найдено, выполняем необходимые действия
				fmt.Printf("Слово '%s' найдено в сообщении '%s'\n", word, message)
			}
		}
	}
}
