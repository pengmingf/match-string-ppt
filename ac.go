package main

import "fmt"

// ACNode AC自动机的节点结构
type ACNode struct {
	children [256]*ACNode // 使用固定大小数组替代map，提高访问速度
	fail     *ACNode      // 失败指针
	isEnd    bool         // 是否是模式串的结尾
	value    string       // 存储该节点对应的完整字符串
	output   []string     // 输出集合缓存
	length   int          // 模式串长度缓存
}

// NewACNode 创建新的AC节点
func NewACNode() *ACNode {
	return &ACNode{
		children: [256]*ACNode{},
		fail:     nil,
		isEnd:    false,
		output:   make([]string, 0, 4), // 预分配空间
	}
}

// AC 定义AC自动机结构
type AC struct {
	root *ACNode
	// 添加字符映射缓存
	charMap map[rune]byte // 字符到数组索引的映射
	maxChar byte          // 最大字符索引
}

// NewAC 创建新的AC自动机
func NewAC() *AC {
	return &AC{
		root:    NewACNode(),
		charMap: make(map[rune]byte, 256),
	}
}

// buildCharMap 构建字符映射
func (ac *AC) buildCharMap(patterns []string) {
	charSet := make(map[rune]bool)
	for _, pattern := range patterns {
		for _, r := range pattern {
			charSet[r] = true
		}
	}

	var index byte
	for r := range charSet {
		ac.charMap[r] = index
		index++
	}
	ac.maxChar = index
}

// Insert 向AC自动机中插入一个模式串
func (ac *AC) Insert(pattern string) {
	node := ac.root
	runes := []rune(pattern)
	for _, r := range runes {
		index := ac.charMap[r]
		if node.children[index] == nil {
			node.children[index] = NewACNode()
		}
		node = node.children[index]
	}
	node.isEnd = true
	node.value = pattern
	node.length = len(runes)
}

// BuildFail 构建失败指针（预处理）
func (ac *AC) BuildFail() {
	// 使用切片替代通用队列，提高性能
	queue := make([]*ACNode, 0, 256)
	queue = append(queue, ac.root)
	ac.root.fail = ac.root

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// 处理当前节点的所有子节点
		for i := byte(0); i < ac.maxChar; i++ {
			child := current.children[i]
			if child == nil {
				continue
			}

			if current == ac.root {
				child.fail = ac.root
			} else {
				failNode := current.fail
				for failNode != ac.root && failNode.children[i] == nil {
					failNode = failNode.fail
				}
				if failNode.children[i] != nil {
					child.fail = failNode.children[i]
					// 合并输出集合
					child.output = append(child.output, child.fail.output...)
				} else {
					child.fail = ac.root
				}
			}
			if child.isEnd {
				child.output = append(child.output, child.value)
			}
			queue = append(queue, child)
		}
	}
}

// BuildAC 预处理构建AC自动机
func BuildAC(patterns []string) *AC {
	ac := NewAC()
	// 首先构建字符映射
	ac.buildCharMap(patterns)
	// 插入所有模式串
	for _, pattern := range patterns {
		ac.Insert(pattern)
	}
	// 构建失败指针
	ac.BuildFail()
	return ac
}

// findNextState 查找下一个状态（优化的状态转移）
func (ac *AC) findNextState(current *ACNode, r rune) *ACNode {
	index := ac.charMap[r]
	for current != ac.root && current.children[index] == nil {
		current = current.fail
	}
	if current.children[index] != nil {
		return current.children[index]
	}
	return ac.root
}

// Search 在文本中搜索所有模式串出现的位置
func (ac *AC) Search(text string) map[string][]int {
	result := make(map[string][]int)
	runes := []rune(text)
	current := ac.root

	// 预分配结果空间
	posCache := make([]int, 0, 64)
	patternCache := make([]string, 0, 64)

	// 遍历文本
	for i, r := range runes {
		current = ac.findNextState(current, r)

		// 使用预计算的输出集合
		if len(current.output) > 0 {
			for _, pattern := range current.output {
				patternCache = append(patternCache, pattern)
				posCache = append(posCache, i-len([]rune(pattern))+1)
			}
		}
	}

	// 批量处理结果
	for i, pattern := range patternCache {
		if result[pattern] == nil {
			result[pattern] = make([]int, 0, 8)
		}
		result[pattern] = append(result[pattern], posCache[i])
	}

	return result
}

func acMain() {
	// 测试示例
	patterns := []string{"he", "she", "his", "hers"}
	text := "ushers"

	// 构建AC自动机
	ac := BuildAC(patterns)

	// 搜索模式串
	matches := ac.Search(text)

	// 输出结果
	for pattern, positions := range matches {
		fmt.Printf("Pattern '%s' found at positions: %v\n", pattern, positions)
	}
}
