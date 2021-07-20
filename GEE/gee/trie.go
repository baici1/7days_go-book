package gee

import "strings"

//实现动态路由最常用的数据结构，被称为前缀树(Trie树)。
//前缀树路由，每个节点的信息
type node struct {
	pattern  string  //待匹配的路由 例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node //子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

// 单个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
//路由插入（注册）
//不断查询每个路由的部分（以/分割的部分）如果有跳过，没有就添加子节点，直到全部完成
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

//路由查询（匹配）
//递归法进行匹配  从根节点获取子节点 子节点查询直到都匹配成功
//退出要求 1.匹配到了* 2.pattern==“” （没有结束）3.匹配到最后一个节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") { //strings.HasPrefix()函数用来检测字符串是否以指定的前缀开头。
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
