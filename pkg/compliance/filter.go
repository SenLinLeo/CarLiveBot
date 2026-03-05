package compliance

import (
	"bufio"
	"os"
	"strings"
)

// Filter 合规词过滤
type Filter struct {
	words map[string]struct{}
}

// NewFilter 从文件加载禁用词（每行一个，# 为注释）
func NewFilter(path string) (*Filter, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	words := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		words[line] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &Filter{words: words}, nil
}

// ContainsForbidden 检查文本是否包含任意禁用词
func (f *Filter) ContainsForbidden(text string) bool {
	for w := range f.words {
		if strings.Contains(text, w) {
			return true
		}
	}
	return false
}

// ReplaceForbidden 将命中词替换为 replaceWith（空则跳过替换仅检测）
func (f *Filter) ReplaceForbidden(text, replaceWith string) string {
	out := text
	for w := range f.words {
		if replaceWith == "" {
			continue
		}
		out = strings.ReplaceAll(out, w, replaceWith)
	}
	return out
}

// Allowed 若不含禁用词返回 true
func (f *Filter) Allowed(text string) bool {
	return !f.ContainsForbidden(text)
}
