package main

import (
	"strconv"
	"strings"
	"testing"
)

type multiPatternTestCase struct {
	name     string
	text     string
	patterns []string
	ac       *acTree
	trie     *Trie
}

// 生成重复文本
func generateRepeatedText(base string, count int) string {
	return strings.Repeat(base, count)
}

// 生成模式串数组
func generatePatterns(base []string, mutations int) []string {
	patterns := make([]string, 0, len(base)*mutations)
	patterns = append(patterns, base...)

	// 生成变种模式串
	for _, p := range base {
		for i := 0; i < mutations-1; i++ {
			patterns = append(patterns, p+strconv.Itoa(i))
		}
	}
	return patterns
}

func BenchmarkMultiPatternMatch(b *testing.B) {
	// 基础模式串集合
	basePatterns := []string{
		"你好", "世界", "测试", "模式串",
		"hello", "world", "test", "pattern",
		"你好世界", "helloworld",
	}

	testCases := []multiPatternTestCase{
		{
			name:     "短文本少模式串",
			text:     "你好，世界！Hello, World!",
			patterns: basePatterns[:4],
		},
		{
			name:     "短文本多模式串",
			text:     "你好，世界！Hello, World!",
			patterns: generatePatterns(basePatterns, 5),
		},
		{
			name:     "中文文本重复模式",
			text:     generateRepeatedText("你好世界测试模式串", 100),
			patterns: generatePatterns([]string{"你好", "世界", "测试"}, 3),
		},
		{
			name:     "英文文本重复模式",
			text:     generateRepeatedText("hello world test pattern", 100),
			patterns: generatePatterns([]string{"hello", "world", "test"}, 3),
		},
		{
			name:     "混合文本大量模式",
			text:     generateRepeatedText("你好hello世界world测试test模式串pattern", 50),
			patterns: generatePatterns(basePatterns, 10),
		},
		{
			name:     "长文本少量模式",
			text:     generateRepeatedText("这是一个很长的测试文本，包含中英文mixed content测试内容", 200),
			patterns: basePatterns[:3],
		},
		{
			name:     "长文本大量模式",
			text:     generateRepeatedText("这是一个很长的测试文本，包含中英文mixed content测试内容", 200),
			patterns: generatePatterns(basePatterns, 20),
		},
		{
			name:     "模式串前缀重叠",
			text:     generateRepeatedText("测试测试测和测试测试测", 100),
			patterns: []string{"测试", "测试测", "测试测试", "测试测试测"},
		},
		{
			name:     "模式串后缀重叠",
			text:     generateRepeatedText("测试测试测和测试测试测试", 100),
			patterns: []string{"试测", "试测试", "试测试测", "试测试测试"},
		},
		{
			name:     "特殊字符混合",
			text:     generateRepeatedText("你好👋世界🌍测试✨模式串💻", 50),
			patterns: []string{"👋世界🌍", "测试✨", "模式串💻", "你好👋"},
		},
		{
			name:     "HTML文本",
			text:     generateRepeatedText("<div class='test'>这是测试文本</div><p>这是段落</p>", 50),
			patterns: []string{"<div", "</div>", "class", "测试文本", "<p>", "</p>"},
		},
		{
			name:     "URL文本",
			text:     generateRepeatedText("https://example.com/测试?param=value&中文=测试", 50),
			patterns: []string{"https://", "example", ".com", "测试", "param=", "中文="},
		},
		{
			name:     "JSON文本",
			text:     generateRepeatedText(`{"name":"测试","value":"test","中文":"数据"}`, 50),
			patterns: []string{`"name"`, `"value"`, `"test"`, `"中文"`, `"数据"`},
		},
	}

	// 预处理所有测试用例
	for i := range testCases {
		ac := NewAc()
		ac.Build(testCases[i].patterns)
		testCases[i].ac = ac

		testCases[i].trie = BuildTrie(testCases[i].patterns)
	}

	// 运行基准测试
	for _, tc := range testCases {
		// Trie树测试
		b.Run("Trie_"+tc.name+"_Patterns_"+strconv.Itoa(len(tc.patterns)), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tc.trie.SearchList(tc.text)
			}
		})

		// 修改AC自动机测试，使用Scan方法
		b.Run("AC_"+tc.name+"_Patterns_"+strconv.Itoa(len(tc.patterns)), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tc.ac.Scan(tc.text)
			}
		})
	}
}

// 测试结果一致性
func TestMatchConsistency(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		patterns []string
	}{
		{
			name:     "基本匹配",
			text:     "你好，世界！Hello, World!",
			patterns: []string{"你好", "世界", "Hello", "World"},
		},
		{
			name:     "重叠模式",
			text:     "测试测试测试",
			patterns: []string{"测试", "测试测", "测试测试"},
		},
		{
			name:     "特殊字符",
			text:     "你好👋世界🌍",
			patterns: []string{"👋世界", "世界🌍", "你好👋"},
		},
		{
			name:     "HTML内容",
			text:     "<div>测试</div>",
			patterns: []string{"<div>", "</div>", "测试"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ac := NewAc()
			ac.Build(tc.patterns)
			trie := BuildTrie(tc.patterns)

			acResult := make(map[string][]int)
			// 将Scan的���果转换为map格式以便比较
			acMatches := ac.Scan(tc.text)
			for _, match := range acMatches {
				if acResult[match] == nil {
					acResult[match] = make([]int, 0)
				}
				// 查找匹配位置
				pos := strings.Index(tc.text, match)
				if pos != -1 {
					acResult[match] = append(acResult[match], pos)
				}
			}

			trieResult := trie.Search(tc.text)

			// 比较结果
			if len(acResult) != len(trieResult) {
				t.Errorf("Result count mismatch: AC=%d, Trie=%d", len(acResult), len(trieResult))
			}

			for pattern, acPos := range acResult {
				triePos, exists := trieResult[pattern]
				if !exists {
					t.Errorf("Pattern %s found in AC but not in Trie", pattern)
					continue
				}

				if len(acPos) != len(triePos) {
					t.Errorf("Position count mismatch for pattern %s: AC=%d, Trie=%d",
						pattern, len(acPos), len(triePos))
				}

				// 比较位置
				for i := range acPos {
					if acPos[i] != triePos[i] {
						t.Errorf("Position mismatch for pattern %s: AC=%d, Trie=%d",
							pattern, acPos[i], triePos[i])
					}
				}
			}
		})
	}
}
