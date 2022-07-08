package gee

type node struct {
	pattern  string
	part     string
	children []*node
	wild     bool
	f        HandleFunc
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.wild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	ret := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.wild {
			ret = append(ret, child)
		}
	}
	return ret
}

func (n *node) insert(pattern string, parts []string, height int, f HandleFunc) {
	if height == len(parts) {
		n.pattern = pattern
		n.f = f
		return
	}
	part := parts[height]
	nn := n.matchChild(part)
	if nn == nil {
		nn = &node{
			part:     part,
			children: make([]*node, 0),
			wild:     n.wild || (len(part) > 0 && (part[0] == ':' || part[0] == '*')),
		}
		n.children = append(n.children, nn)
	}
	nn.insert(pattern, parts, height+1, f)
}

func (n *node) search(pattern string, parts []string, height int) *node {
	if height == len(parts) {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	nns := n.matchChildren(part)
	for _, nn := range nns {
		res := nn.search(pattern, parts, height+1)
		if res != nil {
			return res
		}
	}
	return nil
}
