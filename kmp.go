package main

// getNext 计算模式串的next数组
func getNext(pattern string) []int {
	patternRunes := []rune(pattern)
	m := len(patternRunes)
	next := make([]int, m)
	next[0] = 0
	j := 0

	for i := 1; i < m; {
		if patternRunes[i] == patternRunes[j] {
			j++
			next[i] = j
			i++
		} else if j > 0 {
			j = next[j-1]
		} else {
			next[i] = 0
			i++
		}
	}
	return next
}

// KMPMatch 使用KMP算法查找模式串在主串中的位置
func KMPMatch(text string, pattern string, next []int) int {
	// 将字符串转换为rune切片，以支持中文
	textRunes := []rune(text)
	patternRunes := []rune(pattern)

	n := len(textRunes)
	m := len(patternRunes)

	if m == 0 {
		return 0
	}
	if m > n {
		return -1
	}

	i, j := 0, 0
	for i < n && j < m {
		if textRunes[i] == patternRunes[j] {
			i++
			j++
		} else if j > 0 {
			j = next[j-1]
		} else {
			i++
		}
	}

	if j == m {
		return i - j
	}
	return -1
}

func kmpMain() {
	// 测试示例
	text := "Hello, World!"
	pattern := "World"
	next := getNext(pattern)
	pos := KMPMatch(text, pattern, next)
	println("Pattern found at position:", pos) // 应该输出：7

	// 展示next数组的计算
	pattern = "ABABC"
	next = getNext(pattern)
	println("Next array for pattern", pattern, ":")
	for i, v := range next {
		println("next[", i, "] =", v)
	}
}
