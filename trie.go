package main

import "fmt"

// TrieNode 定义Trie树的节点结构
type TrieNode struct {
	children map[rune]*TrieNode // 子节点映射表
	isEnd    bool               // 标记是否是单词结尾
	value    string             // 存储该节点对应的完整字符串（仅在isEnd为true时有效）
}

// NewTrieNode 创建新的Trie节点
func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		isEnd:    false,
	}
}

// Trie 定义Trie树结构
type Trie struct {
	root *TrieNode
}

// NewTrie 创建新的Trie树
func NewTrie() *Trie {
	return &Trie{
		root: NewTrieNode(),
	}
}

// Insert 向Trie树中插入一个单词
func (t *Trie) Insert(word string) {
	node := t.root
	runes := []rune(word)
	for _, r := range runes {
		if node.children[r] == nil {
			node.children[r] = NewTrieNode()
		}
		node = node.children[r]
	}
	node.isEnd = true
	node.value = word
}

// BuildTrie 预处理构建Trie树
func BuildTrie(patterns []string) *Trie {
	trie := NewTrie()
	for _, pattern := range patterns {
		trie.Insert(pattern)
	}
	return trie
}

// Search 在文本中搜索所有模式串出现的位置
// 返回一个map，key是模式串，value是该模式串在文本中出现的所有位置的切片
func (t *Trie) Search(text string) map[string][]int {
	result := make(map[string][]int)
	runes := []rune(text)
	n := len(runes)

	// 对文本中的每个位置进行匹配
	for i := 0; i < n; i++ {
		node := t.root
		for j := i; j < n; j++ {
			if node.children[runes[j]] == nil {
				break
			}
			node = node.children[runes[j]]
			if node.isEnd {
				if result[node.value] == nil {
					result[node.value] = make([]int, 0)
				}
				result[node.value] = append(result[node.value], i)
			}
		}
	}
	return result
}

func trieMain() {
	// 测试示例
	patterns := []string{"he", "she", "his", "hers"}
	text := "ushers"

	// 构建Trie树
	trie := BuildTrie(patterns)

	// 搜索模式串
	matches := trie.Search(text)

	// 输出结果
	for pattern, positions := range matches {
		fmt.Printf("Pattern '%s' found at positions: %v\n", pattern, positions)
	}
}
