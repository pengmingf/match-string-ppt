package main

import "unicode/utf8"

// BruteForceMatch 使用暴力匹配算法查找模式串在主串中的位置
// 参数:
//   text: 主串
//   pattern: 模式串
// 返回值:
//   找到则返回第一次匹配的起始位置（按rune计算），未找到返回-1
func BruteForceMatch(text string, pattern string) int {
	// 获取主串和模式串的长度（按rune计算）
	n := utf8.RuneCountInString(text)
	m := utf8.RuneCountInString(pattern)

	// 如果模式串为空，返回0
	if m == 0 {
		return 0
	}

	// 如果模式串长度大于主串，不可能匹配成功
	if m > n {
		return -1
	}

	// 将字符串转换为rune切片，以支持中文
	textRunes := []rune(text)
	patternRunes := []rune(pattern)

	// 外层循环遍历主串的每个可能的起始位置
	for i := 0; i <= n-m; i++ {
		j := 0
		// 内层循环比较从i开始的子串是否与模式串匹配
		for j < m && textRunes[i+j] == patternRunes[j] {
			j++
		}
		// 如果j等于模式串长度，说明完全匹配
		if j == m {
			return i
		}
	}

	// 未找到匹配，返回-1
	return -1
}

func blMain() {
	// 测试示例
	text := "Hello, World!"
	pattern := "World"
	pos := BruteForceMatch(text, pattern)
	println("Pattern found at position:", pos) // 应该输出：7
}
