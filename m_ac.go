package main

// 树的节点
type node struct {
	fail   *node          // 失败指针
	isEnd  bool           // 是否词组结尾
	child  map[rune]*node // 子节点
	output []string       // 输出集合缓存
	value  string         // 完整模式串值
}

// 初始化一个节点
func newNode() *node {
	return &node{
		fail:   nil,
		isEnd:  false,
		child:  make(map[rune]*node),
		output: make([]string, 0, 4), // 预分配空间
	}
}

// ac自动机树
type acTree struct {
	root    *node         // root节点
	charSet map[rune]bool // 字符集缓存
}

// NewAc AC自动机，词匹配
func NewAc() *acTree {
	return &acTree{
		root:    newNode(),
		charSet: make(map[rune]bool),
	}
}

// Build 构建树
func (a *acTree) Build(words []string) error {
	// 构建字符集
	for _, word := range words {
		for _, r := range []rune(word) {
			a.charSet[r] = true
		}
	}

	// 构建Trie树
	for _, word := range words {
		nodePtr := a.root
		runes := []rune(word)
		for _, r := range runes {
			if _, ok := nodePtr.child[r]; !ok {
				nodePtr.child[r] = newNode()
			}
			nodePtr = nodePtr.child[r]
		}
		nodePtr.isEnd = true
		nodePtr.value = word
	}

	// 构建fail指针
	a.BuildFail()
	return nil
}

// BuildFail 构建树的fail指针
func (a *acTree) BuildFail() {
	// 使用切片替代通用队列，提高性能
	queue := make([]*node, 0, 256)
	queue = append(queue, a.root)
	a.root.fail = a.root

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// 处理当前节点的所有子节点
		for r, child := range current.child {
			queue = append(queue, child)

			if current == a.root {
				child.fail = a.root
				continue
			}

			// 修正的失败指针构建逻辑
			failNode := current.fail
			for failNode != a.root && failNode.child[r] == nil {
				failNode = failNode.fail
			}

			if failNode.child[r] != nil {
				child.fail = failNode.child[r]
				// 合并输出集合
				child.output = append(child.output, child.fail.output...)
			} else {
				child.fail = a.root
			}

			// 如果是模式串结尾，添加到输出集合
			if child.isEnd {
				child.output = append(child.output, child.value)
			}
		}
	}
}

// findNextState 查找下一个状态（优化的状态转移）
func (a *acTree) findNextState(current *node, r rune) *node {
	for current != a.root && current.child[r] == nil {
		current = current.fail
	}
	if current.child[r] != nil {
		return current.child[r]
	}
	return a.root
}

// Scan 扫描树
func (a *acTree) Scan(text string) []string {
	result := make([]string, 0, 64) // 预分配结果空间
	current := a.root
	runeText := []rune(text)

	// 遍历文本
	for _, r := range runeText {
		current = a.findNextState(current, r)

		// 使用预计算的输出集合
		if len(current.output) > 0 {
			result = append(result, current.output...)
		}
	}

	return result
}
