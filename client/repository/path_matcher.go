package repository

import (
	"regexp"
	"strings"
	"tr369-wss-client/client/model"
)

// MatchType 匹配类型
type MatchType int

const (
	MatchTypeNone     MatchType = iota // 不匹配
	MatchTypeExact                     // 精确匹配
	MatchTypePrefix                    // 前缀匹配（订阅路径是变化路径的前缀）
	MatchTypeWildcard                  // 通配符匹配（{i} 匹配任意实例）
)

// String 返回匹配类型的字符串表示
func (m MatchType) String() string {
	switch m {
	case MatchTypeNone:
		return "None"
	case MatchTypeExact:
		return "Exact"
	case MatchTypePrefix:
		return "Prefix"
	case MatchTypeWildcard:
		return "Wildcard"
	default:
		return "Unknown"
	}
}

// MatchResult 匹配结果
type MatchResult struct {
	Matched     bool      // 是否匹配
	MatchType   MatchType // 匹配类型
	MatchedPath string    // 匹配的订阅路径
}

// PathMatcher 路径匹配器
type PathMatcher struct{}

// NewPathMatcher 创建路径匹配器实例
func NewPathMatcher() *PathMatcher {
	return &PathMatcher{}
}

// Match 检查变化路径是否匹配订阅路径
// subscriptionPath: 订阅的路径（可能包含 {i}）
// changedPath: 发生变化的参数路径
func (m *PathMatcher) Match(subscriptionPath, changedPath string) MatchResult {
	// 1. 精确匹配
	if subscriptionPath == changedPath {
		return MatchResult{
			Matched:     true,
			MatchType:   MatchTypeExact,
			MatchedPath: subscriptionPath,
		}
	}

	// 2. 前缀匹配（订阅路径以 . 结尾，且是变化路径的前缀）
	if strings.HasSuffix(subscriptionPath, ".") && strings.HasPrefix(changedPath, subscriptionPath) {
		return MatchResult{
			Matched:     true,
			MatchType:   MatchTypePrefix,
			MatchedPath: subscriptionPath,
		}
	}

	// 3. 通配符匹配
	if strings.Contains(subscriptionPath, model.WildcardPlaceholder) {
		if m.matchWildcard(subscriptionPath, changedPath) {
			return MatchResult{
				Matched:     true,
				MatchType:   MatchTypeWildcard,
				MatchedPath: subscriptionPath,
			}
		}
	}

	return MatchResult{
		Matched:   false,
		MatchType: MatchTypeNone,
	}
}

// matchWildcard 检查带通配符的路径匹配
// pattern: 包含 {i} 的订阅路径
// path: 发生变化的参数路径
func (m *PathMatcher) matchWildcard(pattern, path string) bool {
	// 将 {i} 替换为正则表达式 \d+（匹配一个或多个数字）
	regexPattern := regexp.QuoteMeta(pattern)
	regexPattern = strings.ReplaceAll(regexPattern, `\{i\}`, `\d+`)

	if strings.HasSuffix(pattern, ".") {
		// 前缀匹配模式：订阅路径以 . 结尾，匹配所有子路径
		regexPattern = "^" + regexPattern
	} else {
		// 精确匹配模式：完全匹配
		regexPattern = "^" + regexPattern + "$"
	}

	matched, err := regexp.MatchString(regexPattern, path)
	if err != nil {
		return false
	}
	return matched
}

// IsPrefix 检查 subscriptionPath 是否是 changedPath 的前缀
func (m *PathMatcher) IsPrefix(subscriptionPath, changedPath string) bool {
	if !strings.HasSuffix(subscriptionPath, ".") {
		return false
	}
	return strings.HasPrefix(changedPath, subscriptionPath)
}
