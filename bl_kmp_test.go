package main

import (
	"strconv"
	"strings"
	"testing"
)

// 测试用例结构
type testCase struct {
	name    string
	text    string
	pattern string
	next    []int // 用于KMP算法的next数组
}

// 生成指定长度的重复字符串
func generateRepeatedString(char string, length int) string {
	return strings.Repeat(char, length)
}

// 生成指定长度的递增字符串 (如: "abcdef...")
func generateSequentialString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = byte('a' + (i % 26))
	}
	return string(result)
}

func BenchmarkStringMatch(b *testing.B) {
	// 准备各种测试场景
	testCases := []testCase{
		// ASCII字符测试
		{
			name:    "短文本最佳情况",
			text:    "hello",
			pattern: "hello",
		},
		{
			name:    "短文本最坏情况",
			text:    "aaaab",
			pattern: "aaab",
		},
		// 中文基础测试
		{
			name:    "中文短文本完全匹配",
			text:    "你好世界",
			pattern: "你好",
		},
		{
			name:    "中文短文本不匹配",
			text:    "你好世界",
			pattern: "再见",
		},
		// 中文混合ASCII
		{
			name:    "中英混合文本",
			text:    "你好Hello世界World",
			pattern: "世界World",
		},
		// 中文重复模式
		{
			name:    "中文重复文本",
			text:    "你好你好你好你好世界",
			pattern: "你好你好",
		},
		{
			name:    "中文重复不匹配",
			text:    generateRepeatedString("你好", 1000) + "世界",
			pattern: generateRepeatedString("你好", 10) + "再见",
		},
		// 特殊中文场景
		{
			name:    "中文标点符号",
			text:    "你好，世界！这是一个测试。",
			pattern: "，世界！",
		},
		{
			name:    "中文空格混合",
			text:    "你好 世界  测试   结果",
			pattern: "世界  测试",
		},
		// 长文本测试
		{
			name:    "长中文文本短模式串",
			text:    generateRepeatedString("测试文本", 1000),
			pattern: "文本测",
		},
		{
			name:    "长中文文本长模式串",
			text:    generateRepeatedString("测试文本", 1000),
			pattern: generateRepeatedString("测试文本", 10),
		},
		// Unicode字符测试
		{
			name:    "Unicode表情符号",
			text:    "你好👋世界🌍真棒👍",
			pattern: "世界🌍",
		},
		// 实际应用场景
		{
			name:    "中文HTML标签",
			text:    "<div>这是一个测试文本</div>",
			pattern: "</div>",
		},
		{
			name:    "中文URL匹配",
			text:    "访问https://例子.com了解更多",
			pattern: "https://例子.com",
		},
		{
			name:    "中文JSON内容",
			text:    `{"name":"测试","value":"这是一个测试值"}`,
			pattern: `"value":"这是`,
		},
	}

	// 预处理KMP的next数组
	for i := range testCases {
		testCases[i].next = getNext(testCases[i].pattern)
	}

	// 对每个测试用例分别运行BF和KMP算法的基准测试
	for _, tc := range testCases {
		// 运行暴力匹配算法测试
		b.Run("BF_"+tc.name+"_Length_"+strconv.Itoa(len([]rune(tc.text))), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				BruteForceMatch(tc.text, tc.pattern)
			}
		})

		// 运行KMP算法测试
		b.Run("KMP_"+tc.name+"_Length_"+strconv.Itoa(len([]rune(tc.text))), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				KMPMatch(tc.text, tc.pattern, tc.next)
			}
		})
	}
}

// 验证两种算法结果一致性的测试
func TestMatchConsistency1(t *testing.T) {
	testCases := []testCase{
		{
			name:    "中文基本匹配",
			text:    "你好，世界！",
			pattern: "世界",
		},
		{
			name:    "中文不匹配",
			text:    "你好，世界！",
			pattern: "再见",
		},
		{
			name:    "中英混合",
			text:    "Hello你好World世界",
			pattern: "World世界",
		},
		{
			name:    "中文重复模式",
			text:    "你好你好你好世界",
			pattern: "你好你好",
		},
		{
			name:    "带标点符号",
			text:    "你好，世界！这是测试。",
			pattern: "，世界！",
		},
		{
			name:    "Unicode表情",
			text:    "你好👋世界🌍",
			pattern: "界🌍",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			next := getNext(tc.pattern)
			bfResult := BruteForceMatch(tc.text, tc.pattern)
			kmpResult := KMPMatch(tc.text, tc.pattern, next)

			if bfResult != kmpResult {
				t.Errorf("Results don't match for case %s: BF=%d, KMP=%d",
					tc.name, bfResult, kmpResult)
			}
		})
	}
}
