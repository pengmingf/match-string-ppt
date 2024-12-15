package main

import (
	"strconv"
	"testing"
)

func BenchmarkKMPMatch(b *testing.B) {
	// 测试用例
	testCases := []struct {
		name    string
		text    string
		pattern string
		next    []int
	}{
		{
			name:    "Best Case",
			text:    "HelloWorld",
			pattern: "Hello",
		},
		{
			name:    "Worst Case",
			text:    "AAAAAAAAAAAAAAAAAAAAAB",
			pattern: "AAAAB",
		},
		{
			name:    "Repeated Pattern",
			text:    "ABABABABABABABABC",
			pattern: "ABABC",
		},
		{
			name:    "Long Text",
			text:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			pattern: "dolor",
		},
		{
			name:    "No Match",
			text:    "HelloWorld",
			pattern: "Python",
		},
	}

	// 预处理所有测试用例的next数组
	for i := range testCases {
		testCases[i].next = getNext(testCases[i].pattern)
	}

	// 运行基准测试
	for _, tc := range testCases {
		b.Run(tc.name+"_Length_"+strconv.Itoa(len(tc.text)), func(b *testing.B) {
			// 重置计时器
			b.ResetTimer()
			// 运行 N 次测试
			for i := 0; i < b.N; i++ {
				KMPMatch(tc.text, tc.pattern, tc.next)
			}
		})
	}
}

// 测试 getNext 函数
func TestGetNext(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expected []int
	}{
		{
			name:     "Simple Pattern",
			pattern:  "ABABC",
			expected: []int{0, 0, 1, 2, 0},
		},
		{
			name:     "Repeated Pattern",
			pattern:  "AAAA",
			expected: []int{0, 1, 2, 3},
		},
		{
			name:     "No Repeats",
			pattern:  "ABCD",
			expected: []int{0, 0, 0, 0},
		},
		{
			name:     "Single Character",
			pattern:  "A",
			expected: []int{0},
		},
		{
			name:     "Empty Pattern",
			pattern:  "",
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getNext(tt.pattern)
			if len(got) != len(tt.expected) {
				t.Errorf("getNext() length = %v, want %v", len(got), len(tt.expected))
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("getNext()[%d] = %v, want %v", i, got[i], tt.expected[i])
				}
			}
		})
	}
}

// 功能测试
func TestKMPMatch(t *testing.T) {
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
		{
			name:     "中文重复模式",
			text:     "你好你好你好世界",
			pattern:  "你好你好",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := getNext(tt.pattern)
			got := KMPMatch(tt.text, tt.pattern, next)
			if got != tt.expected {
				t.Errorf("KMPMatch() = %v, want %v", got, tt.expected)
			}
		})
	}
}
