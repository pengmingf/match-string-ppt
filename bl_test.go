package main

import (
	"strconv"
	"testing"
)

func BenchmarkBruteForceMatch(b *testing.B) {
	// 测试用例
	testCases := []struct {
		text    string
		pattern string
	}{
		// 最佳情况：模式串在开头就匹配
		{
			text:    "HelloWorld",
			pattern: "Hello",
		},
		// 最坏情况：模式串在末尾匹配
		{
			text:    "WorldHello",
			pattern: "Hello",
		},
		// 未匹配情况
		{
			text:    "HelloWorld",
			pattern: "Python",
		},
		// 较长文本的情况
		{
			text:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Hello dolor sit amet.",
			pattern: "dolor",
		},
	}

	// 运行基准测试
	for _, tc := range testCases {
		b.Run("Length_"+strconv.Itoa(len(tc.text)), func(b *testing.B) {
			// 重置计时器
			b.ResetTimer()
			// 运行 N 次测试
			for i := 0; i < b.N; i++ {
				BruteForceMatch(tc.text, tc.pattern)
			}
		})
	}
}

// 功能测试
func TestBruteForceMatch(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		pattern  string
		expected int
	}{
		{
			name:     "Normal match",
			text:     "Hello, World!",
			pattern:  "World",
			expected: 7,
		},
		{
			name:     "Pattern at beginning",
			text:     "Hello, World!",
			pattern:  "Hello",
			expected: 0,
		},
		{
			name:     "No match",
			text:     "Hello, World!",
			pattern:  "Python",
			expected: -1,
		},
		{
			name:     "Empty pattern",
			text:     "Hello, World!",
			pattern:  "",
			expected: 0,
		},
		{
			name:     "Pattern longer than text",
			text:     "Hi",
			pattern:  "Hello",
			expected: -1,
		},
		{
			name:     "中文匹配",
			text:     "你好，世界！",
			pattern:  "世界",
			expected: 3,
		},
		{
			name:     "中英混合",
			text:     "Hello你好World世界",
			pattern:  "你好World",
			expected: 5,
		},
		{
			name:     "中文标点",
			text:     "你好，世界！",
			pattern:  "，",
			expected: 2,
		},
		{
			name:     "表情符号",
			text:     "你好👋世界🌍",
			pattern:  "👋世界",
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BruteForceMatch(tt.text, tt.pattern)
			if got != tt.expected {
				t.Errorf("BruteForceMatch() = %v, want %v", got, tt.expected)
			}
		})
	}
}
